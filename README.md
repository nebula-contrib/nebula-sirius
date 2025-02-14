**nebula-sirius**
=====================

**What is nebula-sirius?**
----------------------------

_**nebula-sirius**_ is a Go library that provides a simple and efficient way to interact with the Nebula Graph database. It allows you to connect to Nebula Graph, execute queries, and retrieve results.

**Why do we need nebula-sirius?**
--------------------------------------

In many applications, interacting with the Nebula Graph database can be complex. nebula-sirius simplifies this process by providing a standardized way to connect to Nebula Graph, execute queries, and retrieve results. This library is particularly useful in building applications that need to interact with the Nebula Graph database.

**What makes current project different from [nebula-go](https://github.com/vesoft-inc/nebula-go)**
----------------------------

The motivation behind creating a new client SDK for the Nebula Graph database is as follows:

- **Robust Connection Pooling**: Current resource management is done poorly in the current [nebula-go](https://github.com/vesoft-inc/nebula-go), which is written in a scratch way. nebula-sirius introduces connection pooling for robust resource management via library, [go-commons-pool](https://github.com/jolestar/go-commons-pool).
- **Context Cancellation**: Better control over requests and handling graceful shutdowns, which is a missing feature in the current nebula-go. nebula-sirius supports context cancellation for improved request handling.
- **Utilization of different code-generator**: Utilized [apache/thrift](https://github.com/apache/thrift) for code generation, which is more feature rich than the current nebula-go's code-generator which is [vesoft-inc/fbthrift](https://github.com/vesoft-inc/fbthrift).

_To bring these features into nebula-go would cause breaking changes, which is critical for projects that currently rely on it. Therefore, we decided to re-write it with a different name, **nebula-sirius**, and release it._

**Installation**
---------------

To install nebula-sirius into your project, use the following command:

```bash
go get github.com/egasimov/nebula-sirius
```

**Usage**
---------

To use nebula-sirius, simply import the library and create a new instance of the `GraphClient` struct.



Here is the restructured Usage section with multiple code snippets for easier understanding:

### Usage

#### Step 1: Configure the Nebula Client Factory

First, we need to configure the Nebula Client Factory with the provided configuration.

```go
clientFactory := nebula_sirius.NewNebulaClientFactory(
	&nebula_sirius.NebulaClientConfig{
		HostAddress: nebula_sirius.HostAddress{
			Host: nebulagraph_light_deployment.HostGraphD,
			Port: nebulagraph_light_deployment.PortGraphD,
		},
	},
	nebula_sirius.DefaultLogger{},
)
```

#### Step 2: Create a Pool of Nebula Clients

Next, we create a pool of Nebula clients based on the client factory and pool configuration.
P:S for full reference of ObjectPoolConfig, please refer to [go-commons-pool documentation](https://pkg.go.dev/github.com/jolestar/go-commons-pool#ObjectPoolConfig)
```go
nebulaClientPool := pool.NewObjectPool(
	ctx,
	clientFactory,
	&pool.ObjectPoolConfig{
		MaxIdle:  5,
		MaxTotal: 10,
		//MaxWaitMillis: 1000,
	},
)
```

#### Step 3: Borrow a Thrift Client from the Pool

We then borrow a Thrift client from the pool.

```go
clientObj, err := nebulaClientPool.BorrowObject(ctx)
if err != nil {
	log.Fatalf("Error borrowing object from pool: %s", err)
}
```

#### Step 4: Get the Graph Client

We take the GraphClient to execute Nebula queries on the Nebula graph service.

```go
g, err := client.GraphClient()
if err != nil {
	log.Fatalf("Error getting graph client: %v", err)
}
```

#### Step 5: Authenticate with the Nebula Database

We make an authentication request to the Nebula database.

```go
a, err := g.Authenticate(
	ctx,
	[]byte(nebulagraph_light_deployment.USERNAME),
	[]byte(nebulagraph_light_deployment.PASSWORD),
)
if err != nil {
	log.Fatalf("Error executing query via graph client: %v", err)
}
```

#### Step 6: Execute a Query

Finally, we execute a query using the GraphClient.

```go
nglQuery := `SHOW HOSTS;`
a1, err := g.Execute(ctx, *a.SessionID, []byte(nglQuery))
if err != nil || a1.GetErrorCode() != nebula.ErrorCode_SUCCEEDED{
	log.Fatalf("Error executing query via graph client: %v %s", err, a1.ErrorMsg)
}
```

#### Step 7: Print the Result Set

We print the result set of the query.

```go
log.Println(nebula_sirius.GenResultSet(a1))
```

**Examples**
--------------
You may refer the working samples located under [examples](./examples) folder.

**Contribution**
--------------

We welcome contributions to nebula-sirius. If you're interested in contributing, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix.
3. Write tests for your changes.
4. Submit a pull request.

**LICENSE**
-------
Note

This project includes code copied and pasted from the Nebula Go repository (https://github.com/vesoft-inc/nebula-go). The original code is licensed under the Apache License 2.0, and we acknowledge the original authors of this code.


nebula-sirius is released under the Apache 2.0 License.

**Code of Conduct**
------------------

We follow the Go community's Code of Conduct. Please read it before contributing.

**Support**
----------

If you have any questions or need help with nebula-sirius, please open an issue on GitHub.

**Thank you for using nebula-sirius!**