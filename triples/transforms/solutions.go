package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type TriplesList []*t.Triples

func (tl TriplesList) String() string {
	res := "[\n"

	for i, triples := range tl {
		res += fmt.Sprintf("%d: %s\n", i, triples.StringLine())
	}
	res += "]\n"
	return res
}

type TripleMatrix []t.TripleList

func (ta TripleMatrix) String() string {
	res := "[\n"

	for _, triples := range ta {
		res += fmt.Sprintf("[%s]\n", triples.StringLine())
	}
	res += "]\n"
	return res
}

func (matrix TripleMatrix) Triples() *t.Triples {
	res := t.NewTriples()
	root := t.NewAnonymousNode()

	for i, row := range matrix {
		node := res.AddTripleReferences(t.NewTriplesFromList(row))
		res.AddTriple(root, t.NewIndexNode(i), node)
	}
	return res
}
