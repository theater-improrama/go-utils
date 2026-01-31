package example

//go:generate go run ./../ -filterable=Filterable -orderable=Orderable

import (
	"context"
	"time"

	"github.com/theater-improrama/go-utils/query"
)

type User interface {
	ID() int
	Name() string
	CreatedAt() time.Time
}

type Filterable interface {
	NameEq(name string)
	CreatedAfter(t time.Time)
}

type Orderable interface {
	CreatedAt()
}

type Repository interface {
	List(
		ctx context.Context,
		opts ...query.Option[Filterable, Orderable],
	)
}
