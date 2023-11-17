package transforms

import "github.com/mholzen/information/triples"

func RowQuery() *triples.Triples {
	// TODO: should return outer most rows, not all nested, which needs more complex queries

	rowQuery := triples.NewTriples()
	x := triples.NewAnonymousNode()
	rowQuery.AddTriple(x, triples.Subject, triples.NodeMatchAnyAnonymous)
	rowQuery.AddTriple(x, triples.Predicate, triples.NodeMatchAnyIndex)
	rowQuery.AddTriple(x, triples.Object, triples.NodeMatchAnyAnonymous)
	return rowQuery
}

func RowTriples(source *triples.Triples) (*triples.Triples, error) {
	queryTripleMatch, err := NewTripleMatchFromTriples(RowQuery())
	if err != nil {
		return nil, err
	}
	rows, err := source.Map(Filter(queryTripleMatch))
	if err != nil {
		return nil, err
	}
	return rows, nil
}
