package transforms

import "github.com/mholzen/information/triples"

func RowQuery() *triples.Triples {
	// TODO: should return outer most rows, not all nested, which needs more complex queries

	rowQuery := triples.NewTriples()
	rowQuery.AddTriple(triples.NewAnonymousNode(), triples.Predicate, triples.NewNodeMatchAnyIndex1())
	rowQuery.AddTriple(triples.NewAnonymousNode(), triples.Object, triples.NewNodeMatchAnyAnonymous1())
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
