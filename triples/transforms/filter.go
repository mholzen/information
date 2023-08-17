package transforms

import (
	. "github.com/mholzen/information/triples"
)

type TripleMatch func(triple Triple) bool

var AlwaysTripleMatch TripleMatch = func(triple Triple) bool {
	return true
}

func NewSubjectTripleMatch(subject Node) TripleMatch {
	return func(triple Triple) bool {
		return triple.Subject == subject
	}
}

func NewPredicateTripleMatch(predicate Node) TripleMatch {
	return func(triple Triple) bool {
		return triple.Predicate == predicate
	}
}

func NewObjectTripleMatch(nodeFn UnaryFunctionNode) TripleMatch {
	return func(triple Triple) bool {
		value, err := nodeFn(triple.Object)
		if err != nil {
			return false
		}

		return int(value.(NumberNode).Value) != 0
	}
}

func NewOrMatch(matches ...TripleMatch) TripleMatch {
	return func(triple Triple) bool {
		for _, match := range matches {
			if match(triple) {
				return true
			}
		}
		return false
	}
}

func NewNotMatch(match TripleMatch) TripleMatch {
	return func(triple Triple) bool {
		return !match(triple)
	}
}

func NewPredicateOrMatch(predicates ...Node) TripleMatch {
	return func(triple Triple) bool {
		for _, predicate := range predicates {
			if triple.Predicate == predicate {
				return true
			}
		}
		return false
	}
}

func NewTripleFilter(destination *Triples, filter TripleMatch) Transformer {
	if destination == nil {
		destination = NewTriples()
	}
	return func(source *Triples) error {
		for _, triple := range source.TripleSet {
			if filter(triple) {
				destination.Add(triple)
			}
		}
		return nil
	}
}

func NewPredicateFilter(destination *Triples, predicate Node) Transformer {
	return NewTripleFilter(destination, NewPredicateTripleMatch(predicate))
}
