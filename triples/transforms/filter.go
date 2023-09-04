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

func NewObjectTripleMatch(object Node) TripleMatch {
	return func(triple Triple) bool {
		return triple.Object == object
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

func And(matches ...TripleMatch) TripleMatch {
	return func(triple Triple) bool {
		for _, match := range matches {
			if !match(triple) {
				return false
			}
		}
		return true
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

func NewPredicateFilter(destination *Triples, predicate Node) Transformer {
	return NewFilterTransformer(destination, NewPredicateTripleMatch(predicate))
}

func NewFilterTransformer(destination *Triples, filter TripleMatch) Transformer {
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

func NewFilterMapper(filter TripleFMatch) Mapper {
	return func(source *Triples) (*Triples, error) {
		res := NewTriples()
		for _, triple := range source.TripleSet {
			f, err := filter(triple)
			if err != nil {
				return nil, err
			}
			if f {
				res.Add(triple)
			}
		}
		return res, nil
	}
}

func Filter(filter TripleMatch) Mapper {
	return func(source *Triples) (*Triples, error) {
		res := NewTriples()
		for _, triple := range source.TripleSet {
			if filter(triple) {
				res.Add(triple)
			}
		}
		return res, nil
	}
}
