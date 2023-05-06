package rel

type Order[T any] interface {
	IsBefore(T, T) Truth[T]
	SetBefore(T, T)
	UnsetBefore(T, T)
}

type List[T comparable] []T

// func (l *List[T]) IsBefore(item1 T, item2 T) Truth[T] {
// 	return false
// }

func (l *List[T]) SetBefore(item1 T, item2 T) {
}

func (l *List[T]) UnsetBefore(item1 T, item2 T) {
}
