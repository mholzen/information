package triples

import (
	"github.com/sirupsen/logrus"
)

type Transformer func(target *Triples) error

func (source *Triples) Transform(transformer Transformer) error {
	return transformer(source)
}

func NewObjectFilter(target *Triples, objectFn UnaryFunctionNode) Transformer {
	return func(source *Triples) error {
		for _, triple := range source.TripleSet {
			value, err := objectFn(triple.Object)
			if err != nil {
				return err
			}

			if value.(NumberNode).Value != 0 {
				target.Add(triple)
			}
		}
		return nil
	}
}

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

func NewTraverse(start Node, filter TripleMatch, dest Node, output *Triples) Transformer {
	visitedNodes := make(NodeSet)
	nodeQueue := make([]Node, 0)
	nodeQueue = append(nodeQueue, start)
	resultIndex := 0

	return func(source *Triples) error {
		for len(nodeQueue) > 0 {
			node := nodeQueue[0]
			nodeQueue = nodeQueue[1:]

			for _, triple := range source.GetTriplesForSubject(node) {
				if !filter(triple) {
					logrus.Debugf("%s fail", triple)
					continue
				}
				logrus.Debugf("%s pass", triple)
				tripleReference := output.AddTripleReference(triple)
				output.NewTripleFromNodes(dest, NewIndexNode(resultIndex), tripleReference)
				resultIndex++

				if !visitedNodes.Contains(triple.Object) {
					visitedNodes.Add(triple.Object)
					nodeQueue = append(nodeQueue, triple.Object)
				}
			}
		}
		return nil
	}
}

type TransformerWithResult struct {
	Transformer
	Result *Node
}
