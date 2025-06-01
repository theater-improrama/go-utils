package optional

type Optional[T any] struct {
	value T
	ok    bool
}

func (o Optional[T]) Value() T {
	return o.value
}

func (o Optional[T]) OK() bool {
	return o.ok
}

func FromNil[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return From(*value)
}

func From[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		ok:    true,
	}
}

func Empty[T any]() Optional[T] {
	return Optional[T]{
		ok: false,
	}
}
