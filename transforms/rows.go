package transforms

import (
	"github.com/mholzen/information/triples"
)

func RowMapper() (triples.Mapper, error) {
	query := RowQuery()
	return func(source *triples.Triples) (*triples.Triples, error) {
		res, err := query.Apply(source)
		if err != nil {
			return nil, err
		}
		return res.GetAllTriples(), nil
	}, nil
}

func RowQuery() Query {
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

	return NewQuery(rowQuery, NewComputations(cRoot, cIndices, cRows))
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	query := RowQuery()
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
