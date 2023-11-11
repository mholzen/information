package transforms

import (
	t "github.com/mholzen/information/triples"
)

// does (x, f, y) means "f(x) = y"
// examples
// (2, "square", 4)
// (f, "computes", "square")

// How do we express that we want to compute the square of 2?
// (2, "square", x)

// A query might contain:
// (x, "square", 4)

// the instantiated query might contain, the needle
// (3, "square", 4)

// The source might contain:
// (marc, age, 3)
// (f, "computes", "square")

// The transformer looks for the instantiated query in the source.
// if it can't find it, it looks for a way to compute it
// (3, "square", 9)

func NewCompute() t.Transformer {
	return func(source *t.Triples) error {
		for _, triple := range source.TripleSet {
			if f, ok := triple.Predicate.(t.UnaryFunctionNode); ok {
				value, err := f(triple.Subject)
				if err != nil {
					return err
				}
				name := t.NewStringNode(f.String())
				source.AddTriple(triple.Subject, name, value)
			}
		}
		return nil
	}
}

func NewComputeWithDefinitions(definitions *t.Triples) t.Transformer {
	return func(source *t.Triples) error {

		for _, definition := range definitions.GetTriplesForPredicate(ComputeNode).TripleSet {
			if f, ok := definition.Subject.(t.UnaryFunctionNode); !ok {
				continue
			} else {
				label := definition.Object

				for _, triple := range source.GetTriplesForPredicate(label).TripleSet {
					if _, ok := triple.Object.(VariableNode); !ok {
						continue
					}

					value, err := f(triple.Subject)
					if err != nil {
						return err
					}
					source.AddTriple(triple.Subject, label, value)
				}
			}
		}
		return nil
	}
}

var ComputeNode = t.NewStringNode("compute")

func ComputeTripleTransformer(subject t.Node, predicate t.UnaryFunctionNode, label t.Node) t.Transformer {
	return func(source *t.Triples) error {
		object, err := predicate(subject)
		if err != nil {
			return err
		}
		_, err = source.AddTriple(subject, label, object)
		if err != nil {
			return err
		}
		return nil
	}
}

func GetDefinitions() *t.Triples {
	var Definitions = t.NewTriples()
	Definitions.AddTriple(t.TypeNode, ComputeNode, "type")
	Definitions.AddTriple(t.SquareNode, ComputeNode, "square")
	Definitions.AddTriple(t.LengthFunction, ComputeNode, "length")
	return Definitions
}
