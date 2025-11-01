package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/haranton/go-graphql-blog/graph"
	"github.com/haranton/go-graphql-blog/graph/directives"
	"github.com/haranton/go-graphql-blog/internals/config"
	"github.com/haranton/go-graphql-blog/internals/middleware"
	"github.com/haranton/go-graphql-blog/internals/service"
	"github.com/haranton/go-graphql-blog/internals/storage"
	"github.com/haranton/go-graphql-blog/internals/storage/memory"
	"github.com/haranton/go-graphql-blog/internals/storage/postgres"
	"github.com/haranton/go-graphql-blog/internals/storage/postgres/migrator"
	"github.com/vektah/gqlparser/v2/ast"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	store  storage.Storage
	srv    *handler.Server
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	var store storage.Storage

	if cfg.Storage.Type == "postgres" {
		logger.Info("Using Postgres storage")
		dbConn := postgres.GetDBConnect(cfg, logger)
		store = postgres.NewPostgresStorage(dbConn)
	} else {
		logger.Info("Using in-memory storage")
		store = memory.NewMemoryStorage()
	}

	serv := service.NewService(store)
	resolver := graph.NewResolver(serv, logger)

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

	return &App{
		cfg:    cfg,
		logger: logger,
		store:  store,
		srv:    srv,
	}
}

func (a *App) MustStart() {

	migrator.MustRunMigrations(a.cfg, a.logger)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.DataLoaderMiddleware(a.store)(a.srv))

	addr := fmt.Sprintf(":%d", a.cfg.App.Port)
	a.logger.Info("GraphQL server running", slog.String("addr", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func (a *App) Close() error {
	a.store.Close()
	return nil
}
