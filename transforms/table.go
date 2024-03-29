package transforms

import (
	"github.com/mholzen/information/triples"
)

func NewTable(headers *triples.Triples) triples.Mapper {
	return func(source *triples.Triples) (*triples.Triples, error) {
		if headers == nil {
			headers = PredicatesSortedLexical(source)
		}

		rows, err := RowTriples(source)
		if err != nil {
			return nil, err
		}

		res := triples.NewTriples()
		root := triples.NewAnonymousNode()
		for i, row := range rows.GetTripleList().Sort() {
			resultRow := triples.NewAnonymousNode()
			for j, header := range headers.GetTripleList().Sort() {

				query := triples.NewTriples()
				query.AddTriple(triples.NewAnonymousNode(), triples.Subject, row.Object)
				query.AddTriple(triples.NewAnonymousNode(), triples.Predicate, header.Object)
				queryTripleMatch, err := NewTripleMatchFromTriples(query)
				if err != nil {
					return nil, err
				}
				cellTriples, err := source.Map(Filter(queryTripleMatch))
				if err != nil {
					return nil, err
				}
				resultCell := triples.NewAnonymousNode()
				for k, cellTriple := range cellTriples.GetTripleList() {
					res.AddTriple(resultCell, triples.NewIndexNode(k), cellTriple.Object)
				}
				res.AddTriple(resultRow, triples.NewIndexNode(j), resultCell)
			}
			res.AddTriple(root, triples.NewIndexNode(i), resultRow)
		}
		return res, nil
	}
}

var TableMapper = NewTable(nil)
