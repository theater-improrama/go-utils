package query

type OrderByFunc[B any] func(b B) B

type QueryBuilder[FB any, OB any] interface {
	Paginate(offset, limit int) QueryBuilder[FB, OB]
	OrderBy(fns ...OrderByFunc[OB]) QueryBuilder[FB, OB]
	Filter(fn FilterPredicate[FB]) QueryBuilder[FB, OB]
}

type FilterPredicate[B any] func(B) B

type FilterBuilderBase[B any] interface {
	Not(fn FilterPredicate[B]) B
	And(fns ...FilterPredicate[B]) B
	Or(fns ...FilterPredicate[B]) B
}

type FilterBase[B FilterBuilderBase[B]] struct{}

func (FilterBase[B]) Empty() FilterPredicate[B] {
	return func(b B) B {
		return b
	}
}

func (FilterBase[B]) Not(fn FilterPredicate[B]) FilterPredicate[B] {
	return func(b B) B {
		return b.Not(fn)
	}
}

func (FilterBase[B]) And(fns ...FilterPredicate[B]) FilterPredicate[B] {
	return func(b B) B {
		return b.And(fns...)
	}
}

func (FilterBase[B]) Or(fns ...FilterPredicate[B]) FilterPredicate[B] {
	return func(b B) B {
		return b.Or(fns...)
	}
}
