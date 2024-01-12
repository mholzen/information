package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type Computation struct {
	Variable VariableNode
	Function t.UnaryFunctionNode
	Expected t.Node
}

func (c Computation) Test(node t.Node) (bool, error) {
	res, err := c.Function(node)
	if err != nil {
		return false, err
	}
	// TODO: useful to log to another set of triples
	// log.Printf("res: %s, expected: %s", res, c.Expected)
	return res == c.Expected, nil
}

func (c Computation) GetTriple() t.Triple {
	return t.NewTripleFromNodes(c.Variable, c.Function, c.Expected)
}

func NewComputation(variable VariableNode, function t.UnaryFunctionNode, expected t.Node) Computation {
	return Computation{
		Variable: variable,
		Function: function,
		Expected: expected,
	}
}

func NewComputationFromTriple(query t.Triple) (Computation, error) {
	res := Computation{}
	variable, ok := query.Subject.(VariableNode)
	if !ok {
		return res, fmt.Errorf("subject %s is not a variable", query.Subject)
	}
	function, ok := query.Predicate.(t.UnaryFunctionNode)
	if !ok {
		return res, fmt.Errorf("predicate %s is not a function", query.Predicate)
	}
	return NewComputation(variable, function, query.Object), nil
}

type Computations []Computation

func NewComputations(computations ...Computation) Computations {
	return Computations(computations)
}

func NewComputationsFromTriples(query *t.Triples) (Computations, error) {
	res := make(Computations, 0)
	for _, triple := range query.TripleSet {
		c, err := NewComputationFromTriple(triple)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (c Computations) Test(vars VariableMap) bool {
	for _, computation := range c {
		if vars[computation.Variable] == nil {
			return false
		}
		ok, err := computation.Test(vars[computation.Variable])
		if err != nil || !ok {
			return false
		}
	}
	return true
}

func (c Computations) GetTriples() *t.Triples {
	res := t.NewTriples()
	for _, computation := range c {
		res.Add(computation.GetTriple())
	}
	return res
}
