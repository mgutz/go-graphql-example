# GraphQL subscription example

Thi is a basic Go graphql example with subscriptions. It's based on this [original
example](https://github.com/matiasanaya/go-graphql-subscription-example). Assets
are externalized and embedded in the binary with [pkger](). A custom `Any`
scalar is implemented to represent any valid JSON value which can be string, int,
object ...

This application uses:

* [github.com/graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)
* [github.com/graph-gophers/graphql-transport-ws](https://github.com/graph-gophers/graphql-transport-ws)


## How to use

Run the application:

```
go run main.go
```

Navigate to [localhost:8080](http://localhost:8080) and use GraphiQL to subscribe using the following example:

```
subscription onHelloSaid {
  helloSaid {
    id
    msg
  }
}
```

On a separate tab run:

```
mutation SayHello{
  sayHello(msg: "Hello world!") {
    id
    msg
  }
}
```

