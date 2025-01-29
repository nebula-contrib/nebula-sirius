/*
 *
 * Copyright (c) 2023 Elchin Gasimov. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nebula_go_sdk

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/egasimov/nebula-go-sdk/nebula/graph"
	"github.com/egasimov/nebula-go-sdk/nebula/meta"
	"github.com/egasimov/nebula-go-sdk/nebula/storage"
	pool "github.com/jolestar/go-commons-pool"
	"golang.org/x/net/http2"
	"math"
	"net"
	"net/http"
	"strconv"
)

// NebulaClientFactory represents a factory for creating new instances of the
// Nebula client.
//
// The factory is configured using the provided NebulaClientConfig and a
// logger. The logger is used to log any errors that occur while creating or
// using the client instances. The client name generator function is used to
// generate a client name for each instance of the Nebula client.
type NebulaClientFactory struct {
	conf              *NebulaClientConfig
	log               Logger
	genClientNameFunc func(ctx context.Context) (string, error)
}

// ClientNameGeneratorFunc is a function that generates a client name based on the context.
type ClientNameGeneratorFunc func(ctx context.Context) (string, error)

// InitNebulaClientFactoryWithDefault creates a new NebulaClientFactory
// with the given configuration and a default logger.
//
// The returned factory can be used to create new instances of the Nebula
// client, which can be used to interact with the Nebula graph database.
//
// The given configuration will be used to initialize the new client
// instances. The default logger will be used to log any errors that occur
// while creating or using the client instances.
func InitNebulaClientFactoryWithDefault(conf *NebulaClientConfig) *NebulaClientFactory {
	return NewNebulaClientFactory(conf, DefaultLogger{}, DefaultClientNameGenerator)
}

// NewNebulaClientFactory creates a new NebulaClientFactory with the given
// configuration and logger.
//
// The returned factory can be used to create new instances of the Nebula
// client, which can be used to interact with the Nebula graph database.
//
// The given configuration will be used to initialize the new client
// instances. The logger will be used to log any errors that occur while
// creating or using the client instances.
func NewNebulaClientFactory(conf *NebulaClientConfig, log Logger, genClientNameFunc ClientNameGeneratorFunc) *NebulaClientFactory {
	return &NebulaClientFactory{
		conf:              conf,
		log:               log,
		genClientNameFunc: genClientNameFunc,
	}
}

// MakeObject is the implementation of the ObjectFactory interface method.
//
// This method will create a new instance of the Nebula client using the
// configuration provided when creating the factory.
//
// The returned PooledObject will contain the newly created client and can
// be used to interact with the Nebula graph database.
//
// The ctx context will be used to generate the client instance. If the
// context is canceled before the client instance is generated, the
// method will return an error.
func (f *NebulaClientFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := f.createWrappedNebulaClient(ctx)
	if err != nil {
		return nil, err
	}

	return pool.NewPooledObject(c), nil
}

// DestroyObject is the implementation of the ObjectFactory interface method.
//
// This method is responsible for destroying a pooled object. It closes the
// transport of the underlying WrappedNebulaClient, effectively terminating
// the connection to the Nebula graph database.
//
// The ctx context is not used directly in this method but is included to
// satisfy the interface requirements.
//
// Returns an error if there is a failure in closing the transport.
func (f *NebulaClientFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client := object.Object.(*WrappedNebulaClient)
	return client.GetTransport().Close()
}

// ValidateObject checks whether the given object is valid or not.
//
// The "validity" of an object is determined by whether its underlying
// transport is open or not.
//
// This is used by the pool to remove dead connections from the pool.
func (f *NebulaClientFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// Check if the context is cancelled before proceeding
	if err := ctx.Err(); err != nil {
		return false
	}

	// do validate
	client := object.Object.(*WrappedNebulaClient)

	// check graph version endpoint ?
	return client.GetTransport().IsOpen()
}

// ActivateObject is called when an object is borrowed from the pool.
// It may be used to reset or initialize the connection. In this case,
// it will open the transport if it is not already open, and then verify
// the client version.
func (f *NebulaClientFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// Optionally reset or initialize the connection
	client := object.Object.(*WrappedNebulaClient)

	if !client.GetTransport().IsOpen() {
		f.log.Debug(fmt.Sprintf("[%s] - client was not open, going to open transport before activated...", client.GetClientName()))
		err := client.GetTransport().Open()
		if err != nil {
			f.log.Error(fmt.Sprintf("[%s] - %v", client.GetClientName(), err))
			return err
		}
		f.log.Debug(fmt.Sprintf("[%s] - client is opened transport, activated succesfully", client.GetClientName()))
	}

	return client.verifyClientVersion(ctx)
}

// PassivateObject is called when an object is returned to the pool.
//
// It may be used to reset or close the connection. In this case,
// it will close the transport if it is already open.
func (f *NebulaClientFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate
	client := object.Object.(*WrappedNebulaClient)
	if client.GetTransport().IsOpen() {
		err := client.GetTransport().Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetClientNameGenerator returns the client name generator function
func (f *NebulaClientFactory) GetClientNameGenerator() ClientNameGeneratorFunc {
	return f.genClientNameFunc
}

// prepareTransportAndProtocolFactory creates a new instance of thrift.TTransport
func (f *NebulaClientFactory) prepareTransportAndProtocolFactory(ctx context.Context) (thrift.TTransport, thrift.TProtocolFactory, error) {
	if ctx.Err() == context.Canceled {
		return nil, nil, ctx.Err()
	}

	hostAddress := f.conf.HostAddress
	timeout := f.conf.Timeout
	sslConfig := f.conf.SslConfig

	newAdd := net.JoinHostPort(hostAddress.Host, strconv.Itoa(hostAddress.Port))

	var transport thrift.TTransport
	var pf thrift.TProtocolFactory
	var sock thrift.TTransport
	if sslConfig != nil {
		sock = thrift.NewTSSLSocketConf(newAdd, &thrift.TConfiguration{
			ConnectTimeout: timeout, // Use 0 for no timeout
			SocketTimeout:  timeout, // Use 0 for no timeout

			TLSConfig: sslConfig,
		})

		//sock, err = thrift.NewTSSLSocketTimeout(newAdd, sslConfig, timeout, timeout)
	} else {
		sock = thrift.NewTSocketConf(newAdd, &thrift.TConfiguration{
			ConnectTimeout: timeout, // Use 0 for no timeout
			SocketTimeout:  timeout, // Use 0 for no timeout
		})
		//sock, err = thrift.NewTSocketTimeout(newAdd, timeout, timeout)
	}

	// Set transport
	bufferSize := 128 << 10
	bufferedTransFactory := thrift.NewTBufferedTransportFactory(bufferSize)
	buffTransport, err := bufferedTransFactory.GetTransport(sock)
	if err != nil {
		return nil, nil, err
	}

	//transport = thrift.NewTHeaderTransport(buffTransport)
	transport = thrift.NewTHeaderTransportConf(buffTransport, &thrift.TConfiguration{})

	//pf = thrift.NewTHeaderProtocolFactory()
	pf = thrift.NewTHeaderProtocolFactoryConf(
		&thrift.TConfiguration{})

	return transport, pf, nil
}

func (f *NebulaClientFactory) getTransportAndProtocolFactoryForHttp2(ctx context.Context) (thrift.TTransport, thrift.TProtocolFactory, error) {
	if ctx.Err() == context.Canceled {
		return nil, nil, ctx.Err()
	}

	hostAddress := f.conf.HostAddress
	sslConfig := f.conf.SslConfig
	httpHeader := f.conf.HttpHeader

	newAdd := net.JoinHostPort(hostAddress.Host, strconv.Itoa(hostAddress.Port))
	var (
		err       error
		transport thrift.TTransport
		pf        thrift.TProtocolFactory
	)

	if sslConfig != nil {
		transport, err = thrift.NewTHttpClientWithOptions("https://"+newAdd,
			thrift.THttpClientOptions{
				Client: &http.Client{
					Transport: &http2.Transport{
						TLSClientConfig: sslConfig,
					},
				},
			})
	} else {
		transport, err = thrift.NewTHttpClientWithOptions("https://"+newAdd, thrift.THttpClientOptions{
			Client: &http.Client{
				Transport: &http2.Transport{
					// So http2.Transport doesn't complain the URL scheme isn't 'https'
					AllowHTTP: true,
					// Pretend we are dialing a TLS endpoint. (Note, we ignore the passed tls.Config)
					DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
						_ = cfg
						var d net.Dialer
						return d.DialContext(ctx, network, addr)
					},
				},
			},
		})

	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a net.Conn-backed Transport,: %s", err.Error())
	}

	//pf = thrift.NewTBinaryProtocolFactoryDefault()
	pf = thrift.NewTBinaryProtocolFactoryConf(&thrift.TConfiguration{})

	if httpHeader != nil {
		client, ok := transport.(*thrift.THttpClient)
		if !ok {
			return nil, nil, fmt.Errorf("failed to get thrift http client")
		}
		for k, vv := range httpHeader {
			if k == "Content-Type" {
				// fbthrift will add "Content-Type" header, so we need to skip it
				continue
			}
			for _, v := range vv {
				// fbthrift set header with http.Header.Add, so we need to set header one by one
				client.SetHeader(k, v)
			}
		}
	}

	return transport, pf, nil
}

// createWrappedNebulaClient creates a new instance of WrappedNebulaClient
func (f *NebulaClientFactory) createWrappedNebulaClient(ctx context.Context) (interface{}, error) {
	var (
		err       error
		transport thrift.TTransport
		pf        thrift.TProtocolFactory
	)

	if f.conf.UseHTTP2 {
		transport, pf, err =
			f.getTransportAndProtocolFactoryForHttp2(ctx)
	} else {
		transport, pf, err = f.prepareTransportAndProtocolFactory(ctx)
	}

	if err != nil {
		f.log.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	graphClient := graph.NewGraphServiceClientFactory(transport, pf)
	metaClient := meta.NewMetaServiceClientFactory(transport, pf)
	storageClient := storage.NewGraphStorageServiceClientFactory(transport, pf)

	clientName, err := f.genClientNameFunc(ctx)
	if err != nil {
		return nil, err
	}
	return newWrappedNebulaClient(graphClient, storageClient, metaClient, transport, clientName, f.log), nil
}

// DefaultClientNameGenerator is a default implementation of ClientNameGenerator.
// It generates a random hex string of length 10 and prefix it with "NebulaClient_".
// The generated name is used to identify the client in the logs.
func DefaultClientNameGenerator(ctx context.Context) (string, error) {
	l := 10
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(buff)
	return fmt.Sprintf("NebulaClient_%s", str[:l]), nil // strip 1 extra character we get from odd length results
}
