package main

import (
	"context"
	"fmt"
	"github.com/egasimov/nebula-go-sdk"
	"github.com/egasimov/nebula-go-sdk/nebula"
	nebulagraph_light_deployment "github.com/egasimov/nebula-go-sdk/nebulagraph-light-deployment"
	"github.com/jolestar/go-commons-pool"
	"log"
)

func main() {
	ctx := context.Background()

	// Configure ClientFactory that serves creation of nebula clients based on the provided configuration
	clientFactory := nebula_go_sdk.NewNebulaClientFactory(
		&nebula_go_sdk.NebulaClientConfig{
			HostAddress: nebula_go_sdk.HostAddress{
				Host: nebulagraph_light_deployment.HostGraphD,
				Port: nebulagraph_light_deployment.PortGraphD,
			},
		},
		nebula_go_sdk.DefaultLogger{},
		nebula_go_sdk.DefaultClientNameGenerator,
	)

	// Build a pool of nebula clients based on clientFactory and poolConfig
	nebulaClientPool := pool.NewObjectPool(
		ctx,
		clientFactory,
		&pool.ObjectPoolConfig{
			MaxIdle:  5,
			MaxTotal: 10,
			//MaxWaitMillis: 1000,
		},
	)

	// Borrow a Thrift client from the pool
	clientObj, err := nebulaClientPool.BorrowObject(ctx)
	if err != nil {
		log.Fatalf("Error borrowing object from pool: %s", err)
	}

	// Return the object to the pool when done
	defer func(thriftPool *pool.ObjectPool, ctx context.Context, object interface{}) {
		err := thriftPool.ReturnObject(ctx, object)
		if err != nil {
			log.Fatalf("Thrift client error: %v", err)
		}
	}(nebulaClientPool, ctx, clientObj)

	client := clientObj.(*nebula_go_sdk.WrappedNebulaClient)

	// Use the client...
	log.Println(fmt.Sprintf("Got a Thrift client with name: %s %v", client.GetClientName(), client))

	// Take GraphClient to execute nebula queries on nebula graph service
	g, err := client.GraphClient()
	if err != nil {
		log.Fatalf("Error getting graph client: %v", err)
	}

	// First, Make authentication request(username, passwd) to nebula database
	a, err := g.Authenticate(
		ctx,
		[]byte(nebulagraph_light_deployment.USERNAME),
		[]byte(nebulagraph_light_deployment.PASSWORD),
	)
	if err != nil {
		log.Fatalf("Error executing query via graph client: %v", err)
	}

	log.Println(fmt.Sprintf("SessionId: %d, ErrorCode: %s, ErrorMessage: %s", a.GetSessionID(), a.GetErrorCode(), a.GetErrorMsg()))

	log.Println(" - - - - - - - - - - - - - - - - - - - - - - - - ")
	nglQuery := `SHOW HOSTS;`
	a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
	if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
		log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
	}

	log.Println(nebula_go_sdk.GenResultSet(a1))

}
