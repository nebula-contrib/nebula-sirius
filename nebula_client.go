package nebula_go_sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/egasimov/nebula-go-sdk/nebula"
	"github.com/egasimov/nebula-go-sdk/nebula/graph"
	"github.com/egasimov/nebula-go-sdk/nebula/meta"
	"github.com/egasimov/nebula-go-sdk/nebula/storage"
	"net/http"
	"time"
)

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

type WrappedNebulaClient struct {
	graphClient   *graph.GraphServiceClient
	metaClient    *meta.MetaServiceClient
	storageClient *storage.GraphStorageServiceClient
	transport     thrift.TTransport
	clientCfg     NebulaClientConfig
	log           Logger
}

func newWrappedNebulaClient(
	graphClient *graph.GraphServiceClient,
	storageClient *storage.GraphStorageServiceClient,
	metaClient *meta.MetaServiceClient,
	transport thrift.TTransport,
	log Logger,
) *WrappedNebulaClient {
	return &WrappedNebulaClient{
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
		c.log.Error(fmt.Sprintf("error: %v", err))
		defer c.transport.Close()
		return err
	}

	if resp.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
		c.log.Error(fmt.Sprintf("incompatible handshakeKey between client and server: %s", string(resp.GetErrorMsg())))
		return fmt.Errorf("incompatible handshakeKey between client and server: %s", string(resp.GetErrorMsg()))
	}
	return nil
}

func (wc *WrappedNebulaClient) Close() error {
	return wc.transport.Close()
}

func (wc *WrappedNebulaClient) GetTransport() thrift.TTransport {
	return wc.transport
}

func (wc *WrappedNebulaClient) GraphClient() (*graph.GraphServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("client opened transport"))
	return wc.graphClient, nil
}

func (wc *WrappedNebulaClient) MetaClient() (*meta.MetaServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("client opened transport"))
	return wc.metaClient, nil
}

func (wc *WrappedNebulaClient) StorageClient() (*storage.GraphStorageServiceClient, error) {
	if err := wc.openTransportIfNeeded(); err != nil {
		wc.log.Error(fmt.Sprintf("%v", err))
		return nil, err
	}

	wc.log.Debug(fmt.Sprintf("client opened transport"))
	return wc.storageClient, nil
}

func (wc *WrappedNebulaClient) openTransportIfNeeded() error {
	if !wc.transport.IsOpen() {
		wc.log.Debug(fmt.Sprintf("client did not open transport, and is going to open transport"))
		err := wc.transport.Open()
		return err
	}

	return nil
}
