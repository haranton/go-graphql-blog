package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/haranton/go-graphql-blog/graph"
	"github.com/haranton/go-graphql-blog/graph/directives"
	"github.com/haranton/go-graphql-blog/internals/service"
	"github.com/haranton/go-graphql-blog/internals/storage/memory"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

type contextKey string

const userCtxKey = contextKey("user")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	store := memory.NewMemoryStorage()

	serv := service.NewService(store)

	resolver := graph.NewResolver(serv)

	schema := graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			Auth: directives.AuthDirective,
		},
	})

	srv := handler.NewDefaultServer(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", basicAuthMiddleware(serv, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func basicAuthMiddleware(serv *service.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := serv.SrvAuth.Authenticate(r.Context(), login, pass)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
