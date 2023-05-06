package rel

type Aggregate[T any] interface {
	Relationship
	Contains(T) bool
	Add(T)
}

type AggregateList[T comparable] []T

func (l *AggregateList[T]) Add(i T) {
	*l = append(*l, i)
}

func (l *AggregateList[T]) Contains(item T) bool {
	for _, t := range *l {
		if t == item {
			return true
		}
	}
	return false
}

func (l *AggregateList[T]) Remove(item T) {
	for i, it := range *l {
		if it == item {
			res := (*l)[0:i]
			res = append(res, (*l)[i+1:]...)
			*l = res
		}
	}
}

type AggregateSet[T comparable] map[T]struct{}

func (a *AggregateSet[T]) Add(item T) {
	(*a)[item] = struct{}{}
}

func (a *AggregateSet[T]) Contains(item T) bool {
	_, ok := (*a)[item]
	return ok
}

func (a *AggregateSet[T]) Remove(item T) {
	delete(*a, item)
}
