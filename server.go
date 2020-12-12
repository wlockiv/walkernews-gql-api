package main

import (
	"github.com/joho/godotenv"
	"github.com/wlockiv/walkernews/internal/auth"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/wlockiv/walkernews/graph"
	"github.com/wlockiv/walkernews/graph/generated"
)

const defaultPort = "8080"

// TODO: Add permissions for Create
// TODO: Integrate JWT flow w/ Fauna
// TODO: Add delete/update for links
// TODO: Update, Deleting Links
// TODO: Add permissions for Updating account
// TODO: Delete the Users controller
// TODO: Remove DynamoDB dependencies
func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
