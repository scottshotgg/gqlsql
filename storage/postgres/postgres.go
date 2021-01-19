package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/scottshotgg/gqlsql/graph/model"
)

type (
	Postgres struct {
		db *sqlx.DB
	}
)

const (
	todosTable = `create table todos (
		id varchar unique,
		text text,
		done boolean
	);`

	getTodos = `select %s from todos;`

	insertTodo = `insert into todos values ($1, $2, $3);`
)

func New(ctx context.Context) (*Postgres, error) {
	var conn, err = sqlx.ConnectContext(ctx, "postgres", "user=postgres password=postgres host=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	_, err = conn.ExecContext(ctx, todosTable)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}

	for i := 0; i < 10; i++ {
		_, err = conn.ExecContext(ctx, insertTodo, uuid.New().String(), "something_here_"+strconv.Itoa(i), i%2 == 0)
		if err != nil {
			return nil, err
		}
	}

	return &Postgres{
		db: conn,
	}, nil
}

func (p *Postgres) GetTODOs(ctx context.Context, fields []string) ([]*model.Todo, error) {
	var todos []*model.Todo

	var query = fmt.Sprintf(getTodos, strings.Join(fields, ", "))

	var err = p.db.SelectContext(ctx, &todos, query)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
