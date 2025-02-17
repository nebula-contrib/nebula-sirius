package main

import (
	"context"
	"fmt"
	"github.com/jolestar/go-commons-pool"
	nebula_sirius "github.com/nebula-contrib/nebula-sirius"
	"github.com/nebula-contrib/nebula-sirius/nebula"
	nebulagraph_light_deployment "github.com/nebula-contrib/nebula-sirius/nebulagraph-light-deployment"
	"log"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()

	// Configure ClientFactory that serves creation of nebula clients based on the provided configuration
	clientFactory := nebula_sirius.NewNebulaClientFactory(
		&nebula_sirius.NebulaClientConfig{
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

	spaceName := "person_pool_test_4"
	tagName := "Persona"
	{
		ngl := `DROP SPACE IF EXISTS %s;
	`
		nglQuery := fmt.Sprintf(ngl, spaceName)
		log.Println(nglQuery)
		a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
		if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
			log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
		}

	}

	time.Sleep(5 * time.Second)
	{
		nglQuery := fmt.Sprintf(`CREATE SPACE %s(partition_num=1, replica_factor=1, vid_type = FIXED_STRING(30));`, spaceName)
		log.Println(nglQuery)
		a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
		if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
			log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
		}

		log.Printf("Space: %s successfully created", spaceName)
	}

	time.Sleep(10 * time.Second)
	{
		nglQuery := fmt.Sprintf(`USE %s; CREATE TAG %s();`, spaceName, tagName)
		log.Println(nglQuery)
		a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
		if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
			log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
		}

		log.Println("Person tag successfully created")
	}

	time.Sleep(10 * time.Second)

	{
		var persons []string
		for i := 0; i < 100; i++ {
			persons = append(persons, fmt.Sprintf(`"P%d":()`, i))
		}

		nglTemplate := `USE %s; INSERT VERTEX %s() VALUES %s;`
		nglQuery := fmt.Sprintf(nglTemplate, spaceName, tagName, strings.Join(persons, ", "))
		log.Println(nglQuery)
		a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
		if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
			log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
		}

		log.Println("Persons are successfully inserted")
	}

	time.Sleep(10 * time.Second)

	{
		nglTemplate := `USE %s; MATCH (v) RETURN count(v) as col1;`
		nglQuery := fmt.Sprintf(nglTemplate, spaceName)
		log.Println(nglQuery)
		a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
		if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED {
			log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
		}

		for _, row := range a1.Data.GetRows() {
			for _, col := range row.Values {
				log.Printf("total person count: %v", col.GetIVal())
			}
		}
	}

}
