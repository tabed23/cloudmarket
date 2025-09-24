package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/tabed23/cloudmarket-auth/graph"
	"github.com/tabed23/cloudmarket-auth/graph/config"
	"github.com/tabed23/cloudmarket-auth/graph/middleware"
	"github.com/tabed23/cloudmarket-auth/graph/repos/store"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err :=godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	db := config.InitDB()
	store := store.NewStore(db)
	r := mux.NewRouter()
	r.Use(middleware.AuthMiddleware)
	c :=  graph.Config{Resolvers: &graph.Resolver{Repository: store}}
	c.Directives.Auth = middleware.Auth 

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))


	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
