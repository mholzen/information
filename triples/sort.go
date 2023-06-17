package triples

import "sort"

func (s TripleList) Sort() {
	sort.Sort(TripleSort{s, func(i, j int) bool {

		if s[i].Subject.LessThan(s[j].Subject) {
			return true
		} else if s[j].Subject.LessThan(s[i].Subject) {
			return false
		}
		// subjects are equal

		if s[i].Predicate.LessThan(s[j].Predicate) {
			return true
		} else if s[j].Predicate.LessThan(s[i].Predicate) {
			return false
		}
		// predicates are equal

		return s[i].Object.LessThan(s[j].Object)
	}})
}

type TripleSort struct {
	data     TripleList
	lessFunc func(i, j int) bool
}

func (s TripleSort) Len() int           { return len(s.data) }
func (s TripleSort) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s TripleSort) Less(i, j int) bool { return s.lessFunc(i, j) }
