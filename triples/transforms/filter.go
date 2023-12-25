package transforms

import (
	"fmt"
	"regexp"

	t "github.com/mholzen/information/triples"
)

func AlwaysTripleMatch(triple t.Triple) bool {
	return true
}

func NewSubjectTripleMatch(subject t.Node) t.TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Subject == subject
	}
}

func NewSubjectsTripleMatch(subjects t.NodeSet) t.TripleMatch {
	return func(triple t.Triple) bool {
		return subjects.Contains(triple.Subject)
	}
}

func NewPredicateTripleMatch(predicate t.Node) t.TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Predicate == predicate
	}
}

func NewObjectTripleMatch(object t.Node) t.TripleMatch {
	return func(triple t.Triple) bool {
		return triple.Object == object
	}
}

func NewObjectRegexpTripleMatch(object *regexp.Regexp) t.TripleMatch {
	return func(triple t.Triple) bool {
		b := []byte(triple.Object.String())
		return object.Match(b)
	}
}

func NewObjectsTripleMatch(objects t.NodeSet) t.TripleMatch {
	return func(triple t.Triple) bool {
		return objects.Contains(triple.Object)
	}
}

func OrMatches(matches ...t.TripleMatch) t.TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if match(triple) {
				return true
			}
		}
		return false
	}
}

func AndMatches(matches ...t.TripleMatch) t.TripleMatch {
	return func(triple t.Triple) bool {
		for _, match := range matches {
			if !match(triple) {
				return false
			}
		}
		return true
	}
}

func NotMatch(match t.TripleMatch) t.TripleMatch {
	return func(triple t.Triple) bool {
		return !match(triple)
	}
}

func NewPredicateFilter(destination *t.Triples, predicate t.Node) t.Transformer {
	return NewFilterTransformer(destination, NewPredicateTripleMatch(predicate))
}

func NewFilterTransformer(destination *t.Triples, filter t.TripleMatch) t.Transformer {
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

func Filter(filter t.TripleMatch) t.Mapper {
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

func NewTripleMatchFromTriples(filter *t.Triples) (t.TripleMatch, error) {
	matches := make([]t.TripleMatch, 0)

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

func ReferenceTripleMatch(triple t.Triple) bool {
	return triple.Predicate == t.Subject ||
		triple.Predicate == t.Predicate ||
		triple.Predicate == t.Object
}

func NewContainedTripleMatch(triples *t.Triples) t.TripleMatch {
	return func(triple t.Triple) bool {
		return triples.Contains(triple)
	}
}
