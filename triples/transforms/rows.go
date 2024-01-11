package transforms

import (
	"github.com/mholzen/information/triples"
)

func RowQuery() *triples.Triples {
	// TODO: should return outer most rows, not all nested, which needs more complex queries

	rowQuery := triples.NewTriples()
	rowQuery.AddTriple(triples.NodeMatchAnyAnonymous, triples.NodeMatchAnyIndex, triples.NodeMatchAnyAnonymous)
	return rowQuery
}

func RowMapper() (triples.Mapper, error) {
	query := RowQuery2()
	return func(source *triples.Triples) (*triples.Triples, error) {
		res, err := query.Apply(source)
		if err != nil {
			return nil, err
		}
		return res.GetAllTriples(), nil
	}, nil
}

func RowQuery2() Query2 {
	root := Var()
	cRoot := NewComputation(root, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	// Consider
	// NewComputation(root, triples.ObjectCountNode, triples.NewIndexNode(0))

	indices := Var()
	cIndices := NewComputation(indices, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))

	rows := Var()
	cRows := NewComputation(rows, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	rowQuery := triples.NewTriples()
	rowQuery.AddTriple(root, indices, rows)

	return NewQuery2(rowQuery, NewComputations(cRoot, cIndices, cRows))
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	query := RowQuery2()
	res, err := query.Apply(source)
	if err != nil {
		return nil, err
	}
	return res.GetAllTriples(), nil
}

func MatrixQuery() *triples.Triples {
	query := triples.NewTriples()
	x := NewVariableNode()
	t := query.AddTripleFromNodes(
		triples.NodeBoolFunction(triples.NodeMatchAnyAnonymous),
		triples.NodeBoolFunction(triples.NodeMatchAnyIndex),
		x)
	query.AddTriple(x, triples.NodeMatchAnyIndex, triples.NodeMatchAny)
	query.AddTripleFromNodes(x, triples.TypeFunctionNode, triples.NewStringNode("triples.AnonymousNode"))

	a := query.AddTripleReference(t)
	query.AddTriple(triples.NewAnonymousNode(), "select", a)
	return query
}
