package transforms

import (
	t "github.com/mholzen/information/triples"
)

type ComputationGenerator struct {
	Variable          VariableNode
	FunctionGenerator UnaryFunctionNodeGenerator
	Expected          t.Node
}

func (cg ComputationGenerator) GetComputation(source *t.Triples) Computation {
	return NewComputation(cg.Variable, cg.FunctionGenerator(source), cg.Expected)
}

type ComputationGenerators []ComputationGenerator

func NewObjectCountFunction(triples *t.Triples) t.UnaryFunctionNode {
	return func(node t.Node) (t.Node, error) {
		return t.NewIndexNode(len(triples.GetTripleListForObject(node))), nil
	}
}

type UnaryFunctionNodeGenerator func(*t.Triples) t.UnaryFunctionNode
