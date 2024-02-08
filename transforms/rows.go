package transforms

import (
	"github.com/mholzen/information/triples"
)

func RowQuery() Query {
	root, predicates, rows := Var(), Var(), Var()

	matching := triples.NewTriples()
	matching.AddTripleFromAny(root, predicates, rows)

	query := NewQuery(matching)
	query.Selected = query.Matching

	rootIsAnon := NewComputation(root, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))
	predicatesAreIndices := NewComputation(predicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))
	rowsAreAnon := NewComputation(rows, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))

	query.Computations = NewComputations(
		rootIsAnon,
		predicatesAreIndices,
		rowsAreAnon,
	)

	rootIsRoot := ComputationGenerator{
		Variable:          root,
		FunctionGenerator: NewObjectCountFunction,
		Expected:          triples.NewIndexNode(0),
	}

	query.ComputationGenerators = ComputationGenerators{rootIsRoot}
	return query
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	return RowQuery().SearchForSelected(source)
}

func MatrixQuery() Query {
	root, rowPredicates, rows, cellPredicates := Var(), Var(), Var(), Var()

	// Match Triples
	matching := triples.NewTriples()
	matching.AddTriple(root, rowPredicates, rows)
	matching.AddTriple(rows, cellPredicates, Var())

	query := NewQuery(matching)

	// Select Triples
	query.Selected = triples.NewTriples()
	query.Selected.AddTriple(root, rowPredicates, rows)

	// Computations
	rootIsAnon := NewComputation(root, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))
	rowPredicatesAreIndices := NewComputation(rowPredicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))
	rowsAreAnon := NewComputation(rows, triples.TypeFunctionNode, triples.Str("triples.AnonymousNode"))
	cellPredicatesAreIndices := NewComputation(cellPredicates, triples.TypeFunctionNode, triples.Str("triples.IndexNode"))

	query.Computations = NewComputations(
		rootIsAnon,
		rowPredicatesAreIndices,
		rowsAreAnon,
		cellPredicatesAreIndices,
	)

	// Computation Generators
	rootIsRoot := ComputationGenerator{
		Variable:          root,
		FunctionGenerator: NewObjectCountFunction,
		Expected:          triples.NewIndexNode(0),
	}
	query.ComputationGenerators = ComputationGenerators{rootIsRoot}

	return query
}
