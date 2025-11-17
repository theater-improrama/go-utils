package example

//go:generate go run ./../ -filter=FilterBuilder -order=OrderByBuilder

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

type FilterBuilder interface {
	query.FilterBuilderBase[FilterBuilder]

	NameEq(name string) FilterBuilder
	CreatedAfter(t time.Time) FilterBuilder
}

type OrderByBuilder interface {
	CreatedAt(o query.Order) OrderByBuilder
}

type Repository interface {
	List(
		ctx context.Context,
		opts ...query.Option[FilterBuilder, OrderByBuilder],
	)
}
