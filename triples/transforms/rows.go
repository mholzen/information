package transforms

import "github.com/mholzen/information/triples"

func RowQuery() *triples.Triples {
	// TODO: should return outer most rows, not all nested, which needs more complex queries

	rowQuery := triples.NewTriples()
	rowQuery.AddTriple(triples.NodeMatchAnyAnonymous, triples.NodeMatchAnyIndex, triples.NodeMatchAnyAnonymous)
	return rowQuery
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	res, err := source.Map(NewQueryMapper(RowQuery()))
	if err != nil {
		return nil, err
	}
	return References(res), nil
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

// func MatrixQuery2() *triples.Triples {
// 	query := triples.NewTriples()
// 	root := triples.NewAnonymousNode()
// 	rows := triples.NewAnonymousNode()
// 	cells := triples.NewAnonymousNode()
// 	query.AddTriple(root, triples.NodeMatchAnyIndex, rows)
// 	query.AddTriple(rows, triples.NodeMatchAnyIndex, cells)
// 	return query
// }
