package set

type Set[T Comparable[T]] interface {
	Contains(T) bool
}

type MapSet[T comparable] map[T]struct{}

func (m MapSet[T]) Contains(item T) bool {
	_, ok := m[item]
	return ok
}

func (m MapSet[T]) Add(item T) {
	m[item] = struct{}{}
}

func NewMapSet[T comparable]() MapSet[T] {
	return make(MapSet[T], 0)
}

type RuleList[T Comparable[T]] interface {
	ContainsAt(int, T) bool
	Contains(List[T]) bool
}

func NewListRuleList[T Comparable[T]](items ...Set[T]) ListRuleList[T] {
	res := make(ListRuleList[T], 0)
	res = append(res, items...)
	return res
}

type ListRuleList[T Comparable[T]] []Set[T]

func (l ListRuleList[T]) ContainsAt(i int, item T) bool {
	if i < 0 || i > len(l) {
		return false
	}
	rule := l[i]
	return rule.Contains(item)
}

func (l ListRuleList[T]) Contains(items List[T]) bool {
	for i, item := range items {
		if i >= len(l) {
			return false
		}
		if !l.ContainsAt(i, item) {
			return false
		}
	}
	return true
}

type Comparable[T any] interface {
	Equal(T) bool
}

type List[T Comparable[T]] []T

func (l List[T]) Equal(t List[T]) bool {
	if len(t) != len(l) {
		return false
	} else {
		for i, item := range l {
			if t[i].Equal(item) {
				return false
			}
		}
		return true
	}
}

func (l List[T]) Contains(t T) bool {
	for _, item := range l {
		if item.Equal(t) {
			return true
		}
	}
	return false
}

func (l List[T]) Intersect(s Set[T]) List[T] {
	return l.Filter(func(t T) bool {
		return s.Contains(t)
	})
}

func (l List[T]) Difference(s Set[T]) List[T] {
	return l.Filter(func(t T) bool {
		return !s.Contains(t)
	})
}

func (l List[T]) Filter(filter func(t T) bool) List[T] {
	res := make(List[T], 0)
	for _, item := range l {
		if filter(item) {
			res = append(res, item)
		}
	}
	return res
}

func NewList[T Comparable[T]](items ...T) List[T] {
	res := make([]T, 0)
	res = append(res, items...)
	return res
}
