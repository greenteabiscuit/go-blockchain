# go-blockchain

To run multiple nodes:

```shell
go run main.go :8080
go run main.go :8081
```

List of APIs:

```shell
$ curl localhost:8080/chain
$ curl localhost:8080/mine
$ curl localhost:8080/transactions/new -d \
'{"sender":"exampleSender", "recipient": "exampleRecipient", "amount":1}'
$ curl localhost:8080/nodes/register -d \
'{"Nodes": "localhost:8081"}'
$ curl localhost:8080/nodes/resolve
```

## Simulating the code.

Run multiple nodes

```shell
$ go run main.go :8080
$ go run main.go :8081
```

Register the node

```shell
$ curl localhost:8080/nodes/register -d \
'{"Nodes": "localhost:8081"}'
```

Mine on the second node

```shell
$ curl localhost:8081/mine
```

The second chain should be longer than the first chain now.
Resolve the conflict.

```shell
$ curl localhost:8080/nodes/resolve
```