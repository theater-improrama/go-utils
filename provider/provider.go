package provider

import "context"

type ProviderFunc[T any] func(ctx context.Context) (T, error)
