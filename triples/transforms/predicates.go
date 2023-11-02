package transforms

import (
	"github.com/mholzen/information/triples"
	. "github.com/mholzen/information/triples"
)

func Predicates(source *Triples) (*Triples, error) {
	res := NewTriples()
	container := NewAnonymousNode()

	for _, row := range source.GetTripleList() {
		res.AddTriple(container, triples.Predicate, row.Predicate)
	}
	return res, nil
}

func PredicatesSortedByString(source *Triples) (*Triples, error) {
	predicates, err := source.Map(Predicates)
	if err != nil {
		return nil, err
	}
	predicateList := predicates.GetTripleList()
	predicateList.SortBy(func(i, j int) bool {
		return predicateList[i].Object.String() < predicateList[j].Object.String()
	})

	res := NewTriples()
	container := NewAnonymousNode()
	for i, triple := range predicateList {
		res.AddTriple(container, NewIndexNode(i), triple.Object)
	}
	return res, nil
}
