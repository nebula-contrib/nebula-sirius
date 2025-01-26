**Nebula Go SDK**
=====================

**What is Nebula Go SDK?**
----------------------------

Nebula Go SDK is a Go library that provides a simple and efficient way to interact with the Nebula Graph database. It allows you to connect to Nebula Graph, execute queries, and retrieve results.

**Why do we need Nebula Go SDK?**
--------------------------------------

In many applications, interacting with the Nebula Graph database can be complex. Nebula Go SDK simplifies this process by providing a standardized way to connect to Nebula Graph, execute queries, and retrieve results. This library is particularly useful in building applications that need to interact with the Nebula Graph database.

**Installation**
---------------

To install Nebula Go SDK into your project, use the following command:

```bash
go get github.com/egasimov/nebula-go-sdk
```

**Usage**
---------

To use Nebula Go SDK, simply import the library and create a new instance of the `GraphClient` struct.

```go
package main

import (
	"context"
	"fmt"
	"github.com/egasimov/nebula-go-sdk"
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
	log.Println(fmt.Sprintf("Got a Thrift client: %v", client))

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
	if err != nil {
		log.Fatalf("Error executing query via graph client: %v", err)
	}

	log.Println(nebula_go_sdk.GenResultSet(a1))

}


```

**Examples**
--------------
You may refer the working samples located under [examples](./examples) folder.

**Contribution**
--------------

We welcome contributions to Nebula Go SDK. If you're interested in contributing, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix.
3. Write tests for your changes.
4. Submit a pull request.

**LICENSE**
-------
Note

This project includes code copied and pasted from the Nebula Go repository (https://github.com/vesoft-inc/nebula-go). The original code is licensed under the Apache License 2.0, and we acknowledge the original authors of this code.


Nebula Go SDK is released under the Apache 2.0 License.

**Code of Conduct**
------------------

We follow the Go community's Code of Conduct. Please read it before contributing.

**Support**
----------

If you have any questions or need help with Nebula Go SDK, please open an issue on GitHub.

**Thank you for using Nebula Go SDK!**