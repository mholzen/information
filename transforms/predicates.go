package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NodeListToTriples(source t.NodeList) *t.Triples {
	res := t.NewTriples()
	container := t.NewAnonymousNode()
	for i, node := range source {
		res.AddTriple(container, t.NewIndexNode(i), node)
	}
	return res
}

func PredicatesSortedLexical(source *t.Triples) *t.Triples {
	list := source.GetPredicateList().SortLexical()
	return NodeListToTriples(list)
}
