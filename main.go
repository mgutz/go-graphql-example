package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mgutz/go-graphql-example/scalars"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"github.com/markbates/pkger"
)

var httpPort = 8080

func init() {
	port := os.Getenv("PORT")
	if port != "" {
		var err error
		httpPort, err = strconv.Atoi(port)
		if err != nil {
			panic(err)
		}
	}
}

func mustReadSchema(filename string) string {
	f, err := pkger.Open("/schemas/main.graphql")
	if err != nil {
		panic(err)
	}
	schema, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return string(schema)
}

func main() {
	schema := mustReadSchema("/schemas/main.graphql")

	// graphiql handler
	assets := http.FileServer(pkger.Dir("/public"))
	http.Handle("/", assets)

	// init graphQL schema
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	s, err := graphql.ParseSchema(string(schema), newResolver(), opts...)
	if err != nil {
		panic(err)
	}

	// graphQL handler
	graphQLHandler := graphqlws.NewHandlerFunc(s, &relay.Handler{Schema: s})
	http.HandleFunc("/graphql", graphQLHandler)

	fmt.Println(fmt.Sprintf("Listening on http://localhost:%d", httpPort))

	// start HTTP server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		panic(err)
	}
}

type resolver struct {
	helloSaidEvents     chan *helloSaidEvent
	helloSaidSubscriber chan *helloSaidSubscriber
}

func newResolver() *resolver {
	r := &resolver{
		helloSaidEvents:     make(chan *helloSaidEvent),
		helloSaidSubscriber: make(chan *helloSaidSubscriber),
	}

	go r.broadcastHelloSaid()

	return r
}

func (r *resolver) Hello() string {
	return "Hello world!"
}

func (r *resolver) SayHello(args struct{ Msg string }) *helloSaidEvent {
	e := &helloSaidEvent{Msg: args.Msg, ID: randomID(), Meta: scalars.NewAny(123)}
	go func() {
		select {
		case r.helloSaidEvents <- e:
		case <-time.After(1 * time.Second):
		}
	}()
	return e
}

type helloSaidSubscriber struct {
	stop   <-chan struct{}
	events chan<- *helloSaidEvent
}

func (r *resolver) broadcastHelloSaid() {
	subscribers := map[string]*helloSaidSubscriber{}
	unsubscribe := make(chan string)

	// NOTE: subscribing and sending events are at odds.
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-r.helloSaidSubscriber:
			subscribers[randomID()] = s
		case e := <-r.helloSaidEvents:
			for id, s := range subscribers {
				go func(id string, s *helloSaidSubscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}

					select {
					case <-s.stop:
						unsubscribe <- id
					case s.events <- e:
					case <-time.After(time.Second):
					}
				}(id, s)
			}
		}
	}
}

func (r *resolver) HelloSaid(ctx context.Context) <-chan *helloSaidEvent {
	c := make(chan *helloSaidEvent)
	// NOTE: this could take a while
	r.helloSaidSubscriber <- &helloSaidSubscriber{events: c, stop: ctx.Done()}

	return c
}

type helloSaidEvent struct {
	ID   string
	Msg  string
	Meta *scalars.Any
}

func randomID() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 16)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
