package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type TripleMatch func(triple t.Triple) bool

var AlwaysTripleMatch TripleMatch = func(triple t.Triple) bool {
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

func NewOrMatch(matches ...TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if match(triple) {
				return true
			}
		}
		return false
	}
}

func And(matches ...TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if !match(triple) {
				return false
			}
		}
		return true
	}
}

func NewNotMatch(match TripleMatch) TripleMatch {
	return func(triple t.Triple) bool {
		return !match(triple)
	}
}

func NewPredicateOrMatch(predicates ...t.Node) TripleMatch {
	return func(triple t.Triple) bool {
		for _, predicate := range predicates {
			if triple.Predicate == predicate {
				return true
			}
		}
		return false
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

func NewFilterMapperFromTriples(filter *t.Triples) t.Mapper {
	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, triple := range source.TripleSet {
			for _, filterTriple := range filter.TripleSet {
				f, ok := filterTriple.Object.(t.UnaryFunctionNode)
				if !ok {
					return nil, fmt.Errorf("object '%s' must be a UnaryFunctionNode", filterTriple.Object)
				}

				n, err := triple.GetNode(filterTriple.Predicate)
				if err != nil {
					return nil, err
				}
				match, err := f(n)
				if err != nil {
					return nil, err
				}
				if match.(t.NumberNode).Value == 1 {
					res.Add(triple)
				}
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

	// log.Printf("filter query: %s", filter)

	for _, filterTriple := range filter.TripleSet {
		filterTriple := filterTriple
		nodeGetter, err := t.GetNodeFunction(filterTriple.Predicate)
		if err != nil {
			return nil, err
		}

		f, ok := filterTriple.Object.(t.NodeBoolFunctionNode)
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
	return And(matches...), nil
}
