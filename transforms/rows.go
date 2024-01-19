package transforms

import (
	"github.com/mholzen/information/triples"
)

func RowQuery() Query {
	root := Var()
	rootIsAnon := NewComputation(root, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	// Consider
	// NewComputation(root, triples.ObjectCountNode, triples.NewIndexNode(0))

	predicates := Var()
	predicatesAreIndices := NewComputation(predicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))

	rows := Var()
	rowsAreAnon := NewComputation(rows, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	selectTriples := triples.NewTriples()
	selectTriples.AddTriple(root, predicates, rows)

	return NewQuery(selectTriples, NewComputations(rootIsAnon, predicatesAreIndices, rowsAreAnon))
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	query := RowQuery()
	res, err := query.Apply(source)
	if err != nil {
		return nil, err
	}
	return res.GetAllTriples(), nil
}

func MatrixQuery() Query {
	// Computations
	root := Var()
	rootIsAnon := NewComputation(root, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	rowPredicates := Var()
	rowPredicatesAreIndices := NewComputation(rowPredicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))

	rows := Var()
	rowsAreAnon := NewComputation(rows, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	cellPredicates := Var()
	cellPredicatesAreIndices := NewComputation(cellPredicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))

	// Selected Triples
	selectTriples := triples.NewTriples()
	selectTriples.AddTriple(root, rowPredicates, rows)
	selectTriples.AddTriple(rows, cellPredicates, Var())

	return NewQuery(selectTriples, NewComputations(
		rootIsAnon,
		rowPredicatesAreIndices,
		rowsAreAnon,
		cellPredicatesAreIndices))
}
