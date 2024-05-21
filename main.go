package main

import (
	"flag"
	"log"
	"net/http"

	"PostCommentService/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func main() {
	useMemory := flag.Bool("useMemory", false, "Use in-memory storage")
	flag.Parse()
	resolver := graph.NewResolver(*useMemory)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
