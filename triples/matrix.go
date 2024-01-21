package triples

import "fmt"

type TriplesList []*Triples

func (tl TriplesList) String() string {
	res := "[\n"

	for i, triples := range tl {
		res += fmt.Sprintf("%d: %s\n", i, triples.StringLine())
	}
	res += "]\n"
	return res
}

func (tl TriplesList) Cartesian() TripleMatrix {
	res := make(TripleMatrix, 0)

	if len(tl) == 0 {
		return res
	}

	first := tl[0]
	rest := tl[1:]

	restProducts := rest.Cartesian()
	if len(first.TripleSet) == 0 {
		return restProducts
	}
	for _, triple := range first.GetTripleList().Sort() {
		if len(restProducts) == 0 {
			restTriples := make(TripleList, 0)
			restTriples = append(restTriples, triple)
			res = append(res, restTriples)
		} else {
			for _, restProduct := range restProducts {
				restTriples := make(TripleList, 0)
				restTriples = append(restTriples, triple)
				restTriples = append(restTriples, restProduct...)
				res = append(res, restTriples)
			}
		}
	}
	return res
}

type TripleMatrix []TripleList

func (ta TripleMatrix) String() string {
	res := "[\n"

	for _, triples := range ta {
		res += fmt.Sprintf("[%s]\n", triples.StringLine())
	}
	res += "]\n"
	return res
}

func (matrix TripleMatrix) Triples() *Triples {
	res := NewTriples()
	root := NewAnonymousNode()

	for i, row := range matrix {
		node := res.AddTripleReferences(NewTriplesFromList(row))
		res.AddTripleFromAny(root, NewIndexNode(i), node)
	}
	return res
}
