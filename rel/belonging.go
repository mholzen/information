package rel

type AggregateDefinedNegatively[T any] interface {
	Aggregate[T]
	Exclude(T)
	Unexclude(T)
	// Contains
}

type MapBelonging[T comparable] struct {
	Members    map[T]struct{}
	NonMembers map[T]struct{}
}

func (b *MapBelonging[T]) contains(item T) *bool {
	if _, ok := b.Members[item]; ok {
		return &ok
	}
	if _, ok := b.NonMembers[item]; ok {
		ok = false
		return &ok
	}
	return nil
}

// func (b *MapBelonging[T]) Boolean() *bool {
// 	return b.contains(item)
// }
