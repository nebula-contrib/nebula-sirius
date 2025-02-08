package nebula_go_sdk

import (
	"context"
	"testing"

	"github.com/egasimov/nebula-go-sdk/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewWrappedNebulaClient(t *testing.T) {
	graphClient := mocks.NewGraphService(t)
	storageClient := mocks.NewGraphStorageService(t)
	metaClient := mocks.NewMetaService(t)
	transport := mocks.NewTTransport(t)
	clientName := "testClient"
	logger := &mocks.Logger{}

	client := newWrappedNebulaClient(graphClient, storageClient, metaClient, transport, clientName, logger)

	assert.NotNil(t, client)
	assert.Equal(t, clientName, client.GetClientName())
	assert.Equal(t, graphClient, client.graphClient)
	assert.Equal(t, storageClient, client.storageClient)
	assert.Equal(t, metaClient, client.metaClient)
	assert.Equal(t, transport, client.GetTransport())
	assert.Equal(t, logger, client.log)
}

func TestWrappedNebulaClient_Close(t *testing.T) {
	transport := mocks.NewTTransport(t)
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName: "testClient",
		transport:  transport,
		log:        logger,
	}

	logger.On("Debug", mock.Anything).Return(nil)
	transport.On("IsOpen").Return(true)
	transport.On("Close").Return(nil)

	err := client.Close()
	assert.NoError(t, err)

	transport.AssertExpectations(t)
}

func TestWrappedNebulaClient_GraphClient(t *testing.T) {
	graphClient := mocks.NewGraphService(t)
	transport := mocks.NewTTransport(t)
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName:  "testClient",
		graphClient: graphClient,
		transport:   transport,
		log:         logger,
	}

	logger.On("Debug", mock.Anything).Return(nil)
	transport.On("IsOpen").Return(true)

	_, err := client.GraphClient()
	assert.NoError(t, err)
}

func TestWrappedNebulaClient_MetaClient(t *testing.T) {
	metaClient := mocks.NewMetaService(t)
	transport := mocks.NewTTransport(t)
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName: "testClient",
		metaClient: metaClient,
		transport:  transport,
		log:        logger,
	}

	logger.On("Debug", mock.Anything).Return(nil)
	transport.On("IsOpen").Return(true)

	_, err := client.MetaClient()
	assert.NoError(t, err)
}

func TestWrappedNebulaClient_StorageClient(t *testing.T) {
	storageClient := mocks.NewGraphStorageService(t)
	transport := &mocks.TTransport{}
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName:    "testClient",
		storageClient: storageClient,
		transport:     transport,
		log:           logger,
	}

	logger.On("Debug", mock.Anything).Return(nil)
	transport.On("IsOpen").Return(true)

	_, err := client.StorageClient()
	assert.NoError(t, err)
}

func TestWrappedNebulaClient_VerifyClientVersion(t *testing.T) {
	ctx := context.Background()

	graphClient := mocks.NewGraphService(t)
	transport := &mocks.TTransport{}
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName:  "testClient",
		graphClient: graphClient,
		transport:   transport,
		log:         logger,
		clientCfg: NebulaClientConfig{
			HandshakeKey: "testKey",
		},
	}

	logger.On("Error", mock.Anything).Return(nil)
	graphClient.On("VerifyClientVersion", ctx, mock.Anything).Return(nil, assert.AnError)
	transport.On("Close").Return(nil)

	err := client.verifyClientVersion(ctx)
	assert.Error(t, err)

	graphClient.AssertExpectations(t)
}

func TestWrappedNebulaClient_OpenTransportIfNeeded(t *testing.T) {
	transport := &mocks.TTransport{}
	logger := &mocks.Logger{}
	client := &WrappedNebulaClient{
		clientName: "testClient",
		transport:  transport,
		log:        logger,
	}

	logger.On("Debug", mock.Anything).Return(nil)
	transport.On("IsOpen").Return(false)
	transport.On("Open").Return(nil)

	err := client.openTransportIfNeeded()
	assert.NoError(t, err)

	transport.AssertExpectations(t)
}
