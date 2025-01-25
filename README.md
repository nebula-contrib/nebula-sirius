## Working sample nebula go sdk

```go
package main

import (
	"context"
	"fmt"
	"github.com/egasimov/nebula-go-sdk"
	"github.com/jolestar/go-commons-pool"
	"log"
)

func main() {
	ctx := context.Background()

	nebulaConnPool := pool.NewObjectPool(
		ctx,
		&nebula_go_sdk.NebulaClientFactory{
			Conf: &nebula_go_sdk.NebulaClientConfig{
				HostAddress: nebula_go_sdk.HostAddress{
					Host: "localhost",
					Port: 9669,
				},
			},
		},
		&pool.ObjectPoolConfig{
			MaxIdle:  5,
			MaxTotal: 10,
			//MaxWaitMillis: 1000,
		},
	)
	// Borrow a Thrift client from the pool
	clientObj, err := nebulaConnPool.BorrowObject(ctx)
	if err != nil {
		log.Fatalf("Error borrowing object from pool: %v", err)
	}

	client := clientObj.(*nebula_go_sdk.WrappedNebulaClient)

	// Return the object to the pool when done
	defer func(thriftPool *pool.ObjectPool, ctx context.Context, object interface{}) {
		err := thriftPool.ReturnObject(ctx, object)
		if err != nil {
			log.Fatalf("Thrift client error: %v", err)
		}
	}(nebulaConnPool, ctx, clientObj)

	g, err := client.GraphClient()
	if err != nil {
		log.Fatalf("Error getting graph client: %v", err)
	}

	a, err := g.Authenticate(
		ctx, []byte("root"), []byte("Zq4mKfV7"))

	fmt.Println(a, err)

	a1, err := g.Execute(ctx, *a.SessionID, []byte("SHOW HOSTS;"))

	fmt.Println(a1, err)

	// Use the client...
	fmt.Printf("Got a Thrift client: %v", client)
	
}

```