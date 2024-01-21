package transforms

import "github.com/mholzen/information/triples"

func ObjectsToStrings(source *triples.Triples) (*triples.Triples, error) {
	res := triples.NewTriples()
	for _, row := range source.TripleSet {
		res.AddTripleFromAny(row.Subject, row.Predicate, triples.NewStringNode(row.Object.String()))
	}
	return res, nil
}
