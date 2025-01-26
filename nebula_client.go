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
	"github.com/egasimov/nebula-go-sdk/nebula"
	"github.com/egasimov/nebula-go-sdk/nebula/graph"
	"github.com/egasimov/nebula-go-sdk/nebula/meta"
	"github.com/egasimov/nebula-go-sdk/nebula/storage"
	"math"
	"net/http"
	"time"
)

// NebulaClientConfig represents the configuration for the Nebula client.
type NebulaClientConfig struct {
	// UseHTTP2 indicates whether to use HTTP2
	UseHTTP2 bool

	// HttpHeader is the http headers for the connection when using HTTP2
	HttpHeader http.Header

	// client handshakeKey, make sure the client handshakeKey is in the white list of NebulaGraph server 'client_white_list'
	HandshakeKey string

	SslConfig *tls.Config

	// HostAddress represents network address as host and port
	HostAddress HostAddress

	// Socket timeout and Socket connection timeout, unit: seconds
	Timeout time.Duration
}

// WrappedNebulaClient represents a client for interacting with the Nebula graph database.
//
// It encapsulates the graph, meta, and storage service clients, along with the transport
// layer and logging functionality. The client can be configured using the provided
// NebulaClientConfig and logs errors and other information using the specified Logger.
type WrappedNebulaClient struct {
	clientName    string
	graphClient   *graph.GraphServiceClient
	metaClient    *meta.MetaServiceClient
	storageClient *storage.GraphStorageServiceClient
	transport     thrift.TTransport
	clientCfg     NebulaClientConfig
	log           Logger
}

// newWrappedNebulaClient creates a new instance of WrappedNebulaClient.
//
// It returns a new instance of WrappedNebulaClient with the given
// configuration and logger. The returned client can be used to interact with
// the Nebula graph database using the provided configuration.
//
// The given configuration will be used to initialize the client instances.
// The logger will be used to log any errors that occur while creating or using
// the client instances.
func newWrappedNebulaClient(
	graphClient *graph.GraphServiceClient,
	storageClient *storage.GraphStorageServiceClient,
	metaClient *meta.MetaServiceClient,
	transport thrift.TTransport,
	log Logger,
) *WrappedNebulaClient {
	return &WrappedNebulaClient{
		clientName:    fmt.Sprintf("NebulaClient_%s", randomBase16String(10)),
		graphClient:   graphClient,
		metaClient:    metaClient,
		storageClient: storageClient,
		transport:     transport,
		log:           log,
	}
}

func (c *WrappedNebulaClient) verifyClientVersion(ctx context.Context) error {
	req := graph.NewVerifyClientVersionReq()
	if c.clientCfg.HandshakeKey != "" {
		req.Version = []byte(c.clientCfg.HandshakeKey)
	}
	resp, err := c.graphClient.VerifyClientVersion(ctx, req)
	if err != nil {
		c.log.Error(fmt.Sprintf("[%s] - error: %v", c.clientName, err))
		defer c.transport.Close()
		return err
	}

	if resp.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
		c.log.Error(fmt.Sprintf("[%s] - incompatible handshakeKey between client and server: %s", c.clientName, string(resp.GetErrorMsg())))
		return fmt.Errorf("incompatible handshakeKey between client and server: %s", string(resp.GetErrorMsg()))
	}
	return nil
}

// Close closes the underlying transport.
// It is safe to call this method multiple times.
func (wc *WrappedNebulaClient) Close() error {
	wc.log.Debug(fmt.Sprintf("[%s] - Closing Nebula client: %+v", wc.clientName, wc.clientName))

	if wc.transport.IsOpen() && wc.transport != nil {
		return wc.transport.Close()
	}
	return nil
}

func (wc *WrappedNebulaClient) GetClientName() string {
	return wc.clientName
}

func (wc *WrappedNebulaClient) GetTransport() thrift.TTransport {
	return wc.transport
}

func (wc *WrappedNebulaClient) GraphClient() (*graph.GraphServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("[%s] - %v", wc.clientName, err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("[%s] - client opened transport", wc.clientName))
	return wc.graphClient, nil
}

func (wc *WrappedNebulaClient) MetaClient() (*meta.MetaServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("[%s] - %v", wc.clientName, err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("[%s] - client opened transport", wc.clientName))
	return wc.metaClient, nil
}

func (wc *WrappedNebulaClient) StorageClient() (*storage.GraphStorageServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("[%s] - %v", wc.clientName, err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("[%s] - client opened transport", wc.clientName))
	return wc.storageClient, nil
}

func (wc *WrappedNebulaClient) openTransportIfNeeded() error {
	if !wc.transport.IsOpen() {
		wc.log.Debug(fmt.Sprintf("[%s] - client did not open transport, and is going to open transport", wc.clientName))
		err := wc.transport.Open()
		return err
	}

	return nil
}

func randomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}
