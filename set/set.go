package set

type Set[V comparable] interface {
	Set(v V)
	Delete(v V)
	Values() []V
	Exists(v V) bool
	Intersect(v Set[V]) Set[V]
	Difference(v Set[V]) Set[V]
}

type set[V comparable] struct {
	vs map[V]bool
}

func (s *set[V]) Set(v V) {
	s.vs[v] = true
}

func (s *set[V]) Delete(v V) {
	delete(s.vs, v)
}

func (s *set[V]) Values() []V {
	vs := make([]V, 0)

	for v, _ := range s.vs {
		vs = append(vs, v)
	}

	return vs
}

func (s *set[V]) Exists(v V) bool {
	b, ok := s.vs[v]
	if !ok {
		return false
	}

	return b
}

func (s *set[V]) Intersect(v Set[V]) Set[V] {
	intersectS := New[V]()

	for k, _ := range s.vs {
		if v.Exists(k) {
			intersectS.Set(k)
		}
	}

	return intersectS
}

// Difference returns a new set for the set operation A-B (all elements of A minus the elements of B), where
//   - A is the set s, and
//   - B is the set v
func (s *set[V]) Difference(v Set[V]) Set[V] {
	diffS := New[V]()

	for k, _ := range s.vs {
		if v.Exists(k) {
			continue
		}

		diffS.Set(k)
	}

	return diffS
}

func New[V comparable]() Set[V] {
	return &set[V]{
		vs: make(map[V]bool),
	}
}

func makeMap[V comparable](vs []V) map[V]bool {
	m := make(map[V]bool)

	for _, v := range vs {
		m[v] = true
	}

	return m
}

func From[V comparable](vs []V) Set[V] {
	return &set[V]{
		vs: makeMap(vs),
	}
}
