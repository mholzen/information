package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type Solutions struct {
	Query *Query
	// Rows         TripleMatrix
	SolutionList SolutionList
}

// func (s Solutions) Triples() *t.Triples {
// 	return s.Rows.Triples()
// }

// func (s Solutions) SelectTriples() *t.Triples {
// 	res := t.NewTriples()
// 	for i, triple := range s.Query.QueryTriples {
// 		if s.Query.IsSelect(triple) {
// 			for _, row := range s.Rows {
// 				res.Add(row[i])
// 			}
// 		}
// 	}
// 	return res
// }

// func (q Query) GetSolutions(source *t.Triples) (Solutions, error) {
// 	res := Solutions{}
// 	res.Query = &q
// 	solutions, err := q.Apply(source)
// 	if err != nil {
// 		return res, err
// 	}
// 	res.Rows = solutions
// 	return res, nil

// }

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
