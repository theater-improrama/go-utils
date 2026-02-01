package optional

type Optional[T any] struct {
	Value T
	IsSet bool
}

func FromNil[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return From(*value)
}

func From[T any](value T) Optional[T] {
	return Optional[T]{
		Value: value,
		IsSet: true,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{
		IsSet: false,
	}
}
