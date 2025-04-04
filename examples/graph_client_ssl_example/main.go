package main

import (
	"context"
	"fmt"
	"github.com/jolestar/go-commons-pool"
	nebula_sirius "github.com/nebula-contrib/nebula-sirius"
	nebulagraph_light_deployment "github.com/nebula-contrib/nebula-sirius/nebulagraph-light-deployment"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	wdDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rootCAPath := fmt.Sprintf("%s/nebulagraph-light-deployment/secrets/ca.crt", wdDir)
	certPath := fmt.Sprintf("%s/nebulagraph-light-deployment/secrets/client.crt", wdDir)
	privateKeyPath := fmt.Sprintf("%s/nebulagraph-light-deployment/secrets/client.key", wdDir)

	// Configure SSL
	sslConfig, err := nebula_sirius.GetDefaultSSLConfig(
		rootCAPath,
		certPath,
		privateKeyPath,
	)
	sslConfig.InsecureSkipVerify = false

	if err != nil {
		log.Fatal(fmt.Sprintf("%s.", err.Error()))
	}

	// Configure ClientFactory that serves creation of nebula clients based on the provided configuration
	clientFactory := nebula_sirius.NewNebulaClientFactory(
		&nebula_sirius.NebulaClientConfig{
			SslConfig: sslConfig,
			HostAddress: nebula_sirius.HostAddress{
				Host: nebulagraph_light_deployment.HostGraphD,
				Port: nebulagraph_light_deployment.PortGraphD,
			},
		},
		nebula_sirius.DefaultLogger{},
		nebula_sirius.DefaultClientNameGenerator,
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

	client := clientObj.(*nebula_sirius.WrappedNebulaClient)

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
	if err != nil {
		log.Fatalf("Error executing query via graph client: %v", err)
	}

	log.Println(nebula_sirius.GenResultSet(a1))

}
