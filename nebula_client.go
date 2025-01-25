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
}

func newWrappedNebulaClient(
	graphClient *graph.GraphServiceClient,
	storageClient *storage.GraphStorageServiceClient,
	metaClient *meta.MetaServiceClient,
	transport thrift.TTransport,
) *WrappedNebulaClient {
	return &WrappedNebulaClient{
		graphClient:   graphClient,
		metaClient:    metaClient,
		storageClient: storageClient,
		transport:     transport,
	}
}

func (c *WrappedNebulaClient) verifyClientVersion(ctx context.Context) error {
	req := graph.NewVerifyClientVersionReq()
	if c.clientCfg.HandshakeKey != "" {
		req.Version = []byte(c.clientCfg.HandshakeKey)
	}
	resp, err := c.graphClient.VerifyClientVersion(ctx, req)
	if err != nil {
		defer c.transport.Close()
		return err
	}

	if resp.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
		return fmt.Errorf("incompatible handshakeKey between client and server: %s", string(resp.GetErrorMsg()))
	}
	return nil
}

func (r *WrappedNebulaClient) Close() error {
	return r.transport.Close()
}

func (r *WrappedNebulaClient) GetTransport() thrift.TTransport {
	return r.transport
}

func (r *WrappedNebulaClient) GraphClient() (*graph.GraphServiceClient, error) {
	if !r.transport.IsOpen() {
		err := r.transport.Open()
		if err != nil {
			return nil, err
		}
	}
	return r.graphClient, nil
}

func (r *WrappedNebulaClient) MetaClient() (*meta.MetaServiceClient, error) {
	if !r.transport.IsOpen() {
		err := r.transport.Open()
		if err != nil {
			return nil, err
		}
	}
	return r.metaClient, nil

}

func (r *WrappedNebulaClient) StorageClient() (*storage.GraphStorageServiceClient, error) {
	if !r.transport.IsOpen() {
		err := r.transport.Open()
		if err != nil {
			return nil, err
		}
	}
	return r.storageClient, nil
}
