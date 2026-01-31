package query

type Option[FB, OB any] func(b Builder[FB, OB])

type OrderByFunc[B any] func(b B) B

type Builder[FB any, OB any] interface {
	Paginate(offset, limit int) Builder[FB, OB]
	OrderBy(fns ...OrderByFunc[OB]) Builder[FB, OB]
	Filter(fn FilterPredicate[FB]) Builder[FB, OB]
}

type FilterPredicate[B any] func(B) B

// FilterBuilderLogic provides logical operations for filter builders.
// This is embedded by generated filter builder interfaces.
type FilterBuilderLogic[B any] interface {
	Not(fn FilterPredicate[B]) B
	And(fns ...FilterPredicate[B]) B
	Or(fns ...FilterPredicate[B]) B
}

// FilterBase provides helper methods for constructing filter predicates.
// It is embedded by the generated Filter variable to provide Not/And/Or wrappers.
type FilterBase[B FilterBuilderLogic[B]] struct{}

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
