package storage

import (
	"context"

	"github.com/scottshotgg/gqlsql/graph/model"
)

type (
	Storage interface {
		GetTODOs(ctx context.Context, fields []string) ([]*model.Todo, error)
	}
)
