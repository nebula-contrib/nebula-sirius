package nebula_go_sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/egasimov/nebula-go-sdk/nebula/graph"
	"github.com/egasimov/nebula-go-sdk/nebula/meta"
	"github.com/egasimov/nebula-go-sdk/nebula/storage"
	pool "github.com/jolestar/go-commons-pool"
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"strconv"
)

type NebulaClientFactory struct {
	conf *NebulaClientConfig
	log  Logger
}

func InitNebulaClientFactoryWithDefaultLogger(conf *NebulaClientConfig) *NebulaClientFactory {
	return NewNebulaClientFactory(conf, DefaultLogger{})
}

func NewNebulaClientFactory(conf *NebulaClientConfig, log Logger) *NebulaClientFactory {
	return &NebulaClientFactory{
		conf: conf,
		log:  log,
	}
}

func (f *NebulaClientFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	c, err := f.createWrappedNebulaClient(ctx)
	if err != nil {
		return nil, err
	}

	return pool.NewPooledObject(c), nil
}

func (f *NebulaClientFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client := object.Object.(*WrappedNebulaClient)
	return client.GetTransport().Close()
}

func (f *NebulaClientFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	client := object.Object.(*WrappedNebulaClient)

	// check graph version endpoint ?
	return client.GetTransport().IsOpen()
}

func (f *NebulaClientFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// Optionally reset or initialize the connection
	client := object.Object.(*WrappedNebulaClient)

	if !client.GetTransport().IsOpen() {
		f.log.Debug(fmt.Sprintf("[%s] - client was not open, going to open transport before activated...", client.GetClientName()))
		err := client.GetTransport().Open()
		if err != nil {
			f.log.Error(fmt.Sprintf("%v", err))
			return err
		}
		f.log.Debug(fmt.Sprintf("[%s] - client is opened transport, activated succesfully", client.GetClientName()))
	}

	return client.verifyClientVersion(ctx)
}

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

// create socket based transport
func (f *NebulaClientFactory) prepareTransportAndProtocolFactory(ctx context.Context) (thrift.TTransport, thrift.TProtocolFactory, error) {
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

	transport = thrift.NewTHeaderTransport(buffTransport)

	//pf = thrift.NewTHeaderProtocolFactory()
	pf = thrift.NewTHeaderProtocolFactoryConf(
		&thrift.TConfiguration{})

	return transport, pf, nil
}

func (f *NebulaClientFactory) getTransportAndProtocolFactoryForHttp2(ctx context.Context) (thrift.TTransport, thrift.TProtocolFactory, error) {
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

// Factory function to create new Thrift client
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

	return newWrappedNebulaClient(graphClient, storageClient, metaClient, transport, f.log), nil
}
