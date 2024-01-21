package transforms

import (
	t "github.com/mholzen/information/triples"
)

func Predicates(source *t.Triples) (*t.Triples, error) {
	res := t.NewTriples()
	container := t.NewAnonymousNode()

	for _, row := range source.GetTripleList() {
		res.AddTripleFromAny(container, t.Predicate, row.Predicate)
	}
	return res, nil
}

func PredicatesSortedByString(source *t.Triples) (*t.Triples, error) {
	predicates, err := source.Map(Predicates)
	if err != nil {
		return nil, err
	}
	predicateList := predicates.GetTripleList()
	predicateList.SortBy(func(i, j int) bool {
		return predicateList[i].Object.String() < predicateList[j].Object.String()
	})

	res := t.NewTriples()
	container := t.NewAnonymousNode()
	for i, triple := range predicateList {
		res.AddTripleFromAny(container, t.NewIndexNode(i), triple.Object)
	}
	return res, nil
}
