package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/scottshotgg/gqlsql/graph/generated"
	"github.com/scottshotgg/gqlsql/graph/model"
	"github.com/scottshotgg/gqlsql/storage"
	"github.com/scottshotgg/gqlsql/storage/postgres"
)

const defaultPort = "8080"

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var (
		ctx    = context.Background()
		s, err = postgres.New(ctx)
	)

	if err != nil {
		log.Fatalln("err:", err)
	}

	var (
		conf = generated.Config{
			Resolvers: &resolver{
				s: s,
			},
		}

		srv = handler.NewDefaultServer(generated.NewExecutableSchema(conf))
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type resolver struct {
	s storage.Storage
}
type qresolver struct{}
type mresolver struct{}

func (r *resolver) Query() generated.QueryResolver {
	return r
}

func (r *resolver) Mutation() generated.MutationResolver {
	return r
}

func (r *resolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.s.GetTODOs(ctx, GetPreloads(ctx))
}

func (r *resolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	return r.s.CreateTODO(ctx, input.ID, input.Text, input.Done)
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) []string {
	var preloads []string

	for _, column := range fields {
		preloads = append(preloads, GetPreloadString(prefix, column.Name))
		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), preloads[0])...)
	}

	return preloads
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}

	return name
}
