package transforms

import "github.com/mholzen/information/triples"

func ObjectsToStrings(source *triples.Triples) (*triples.Triples, error) {
	res := triples.NewTriples()
	for _, row := range source.TripleSet {
		res.AddTriple(row.Subject, row.Predicate, triples.Str(row.Object.String()))
	}
	return res, nil
}
