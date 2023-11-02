package triples

import "sort"

type TripleList []Triple

func (l TripleList) GetObjects() []Node {
	var objects []Node
	for _, triple := range l {
		objects = append(objects, triple.Object)
	}
	return objects
}

func (l TripleList) GetObjectStrings() []string {
	var objects []string
	for _, triple := range l {
		objects = append(objects, triple.Object.String())
	}
	return objects
}

func (l TripleList) Sort() TripleList {
	sort.Sort(TripleSort{l, func(i, j int) bool {

		if l[i].Subject.LessThan(l[j].Subject) {
			return true
		} else if l[j].Subject.LessThan(l[i].Subject) {
			return false
		}

		// subjects are equal
		if l[i].Predicate.LessThan(l[j].Predicate) {
			return true
		} else if l[j].Predicate.LessThan(l[i].Predicate) {
			return false
		}

		// predicates are equal
		return l[i].Object.LessThan(l[j].Object)
	}})
	return l
}

func (l TripleList) SortBy(f func(i, j int) bool) TripleList {
	sort.Sort(TripleSort{l, f})
	return l
}

type TripleSort struct {
	data     TripleList
	lessFunc func(i, j int) bool
}

func (s TripleSort) Len() int           { return len(s.data) }
func (s TripleSort) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s TripleSort) Less(i, j int) bool { return s.lessFunc(i, j) }
