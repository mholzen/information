package transforms

import (
	t "github.com/mholzen/information/triples"
)

func Cartesian(sets TriplesList) TripleMatrix {
	res := make(TripleMatrix, 0)

	if len(sets) == 0 {
		return res
	}

	first := sets[0]
	rest := sets[1:]

	restProducts := Cartesian(rest)
	if len(first.TripleSet) == 0 {
		return restProducts
	}
	for _, triple := range first.GetTripleList().Sort() {
		if len(restProducts) == 0 {
			restTriples := make(t.TripleList, 0)
			restTriples = append(restTriples, triple)
			res = append(res, restTriples)
		} else {
			for _, restProduct := range restProducts {
				restTriples := make(t.TripleList, 0)
				restTriples = append(restTriples, triple)
				restTriples = append(restTriples, restProduct...)
				res = append(res, restTriples)
			}
		}
	}
	return res
}
