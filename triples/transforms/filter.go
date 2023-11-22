package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type TripleMatch func(triple t.Triple) bool

func AlwaysTripleMatch(triple t.Triple) bool {
	return true
}

func NewSubjectTripleMatch(subject t.Node) TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Subject == subject
	}
}

func NewPredicateTripleMatch(predicate t.Node) TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Predicate == predicate
	}
}

func NewObjectTripleMatch(object t.Node) TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Object == object
	}
}

func OrMatches(matches ...TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if match(triple) {
				return true
			}
		}
		return false
	}
}

func AndMatches(matches ...TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if !match(triple) {
				return false
			}
		}
		return true
	}
}

func NotMatch(match TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		return !match(triple)
	}
}

func NewPredicateFilter(destination *t.Triples, predicate t.Node) t.Transformer {
	return NewFilterTransformer(destination, NewPredicateTripleMatch(predicate))
}

func NewFilterTransformer(destination *t.Triples, filter TripleMatch) t.Transformer {
	if destination == nil {
		destination = t.NewTriples()
	}
	return func(source *t.Triples) error {
		for _, triple := range source.TripleSet {
			if filter(triple) {
				destination.Add(triple)
			}
		}
		return nil
	}
}

func NewFilterMapper(filter TripleFMatch) t.Mapper {
	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
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

func Filter(filter TripleMatch) t.Mapper {
	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, triple := range source.TripleSet {
			if filter(triple) {
				res.Add(triple)
			}
		}
		return res, nil
	}
}

func NewTripleMatchFromTriples(filter *t.Triples) (TripleMatch, error) {
	matches := make([]TripleMatch, 0)

	for _, filterTriple := range filter.TripleSet {
		filterTriple := filterTriple
		nodeGetter, err := t.GetNodeFunction(filterTriple.Predicate)
		if err != nil {
			return nil, fmt.Errorf("error getting node function for triple '%s': %w", filterTriple, err)
		}

		f, ok := filterTriple.Object.(t.NodeBoolFunction)
		if !ok {
			f = func(n t.Node) bool {
				return n == filterTriple.Object
			}
		}

		match := func(triple t.Triple) bool {
			return f(nodeGetter(triple))
		}
		matches = append(matches, match)
	}
	return AndMatches(matches...), nil
}
