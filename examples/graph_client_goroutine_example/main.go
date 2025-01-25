package main

import (
	"context"
	"fmt"
	"github.com/egasimov/nebula-go-sdk"
	nebulagraph_light_deployment "github.com/egasimov/nebula-go-sdk/nebulagraph-light-deployment"
	"github.com/jolestar/go-commons-pool"
	"log"
	"sync"
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
	)

	// Build a pool of nebula clients based on clientFactory and poolConfig
	nebulaClientPool := pool.NewObjectPool(
		ctx,
		clientFactory,
		&pool.ObjectPoolConfig{
			MaxIdle:  1,
			MaxTotal: 2,
			//MaxWaitMillis: 1000,
			TestOnCreate:       true,
			TestOnBorrow:       true,
			TestWhileIdle:      true,
			TestOnReturn:       true,
			BlockWhenExhausted: true,
		},
	)

	var wg sync.WaitGroup
	goroutineCnt := 10
	wg.Add(goroutineCnt)

	for i := 0; i < goroutineCnt; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

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
			log.Println(fmt.Sprintf("Got a Thrift client: %v", client))

			if err := ExecSomeQuery(ctx, client); err != nil {
				if err != nil {
					log.Fatalf("Error getting graph client: %v", err)
				}
			}

		}(&wg)
	}

	wg.Wait()

	log.Println("Application finished.")
}

func ExecSomeQuery(ctx context.Context, client *nebula_go_sdk.WrappedNebulaClient) error {
	// Take GraphClient to execute nebula queries on nebula graph service
	g, err := client.GraphClient()

	if err != nil {
		return err
	}

	// First, Make authentication request(username, passwd) to nebula database
	a, err := g.Authenticate(
		ctx,
		[]byte(nebulagraph_light_deployment.USERNAME),
		[]byte(nebulagraph_light_deployment.PASSWORD),
	)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("SessionId: %d, ErrorCode: %s, ErrorMessage: %s", a.GetSessionID(), a.GetErrorCode(), a.GetErrorMsg()))

	log.Println(" - - - - - - - - - - - - - - - - - - - - - - - - ")

	nglQuery := `SHOW HOSTS;`
	a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
	if err != nil {
		return err
	}

	log.Println(nebula_go_sdk.GenResultSet(a1))

	return nil
}
