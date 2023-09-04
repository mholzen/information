package transforms

import (
	. "github.com/mholzen/information/triples"
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

func NewCompute() Transformer {
	return func(source *Triples) error {
		for _, triple := range source.TripleSet {
			if f, ok := triple.Predicate.(UnaryFunctionNode); ok {
				value, err := f(triple.Subject)
				if err != nil {
					return err
				}
				name := NewStringNode(f.String())
				source.AddTriple(triple.Subject, name, value)
			}
		}
		return nil
	}
}

func NewComputeWithDefinitions(definitions *Triples) Transformer {
	return func(source *Triples) error {

		for _, definition := range definitions.GetTriplesForPredicate(ComputeNode).TripleSet {
			if f, ok := definition.Subject.(UnaryFunctionNode); !ok {
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

var ComputeNode = NewStringNode("compute")

func ComputeTripleTransformer(subject Node, predicate UnaryFunctionNode, label Node) Transformer {
	return func(source *Triples) error {
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

func GetDefinitions() *Triples {
	var Definitions = NewTriples()
	Definitions.AddTriple(TypeNode, ComputeNode, "type")
	Definitions.AddTriple(SquareNode, ComputeNode, "square")
	Definitions.AddTriple(LengthFunction, ComputeNode, "length")
	return Definitions
}
