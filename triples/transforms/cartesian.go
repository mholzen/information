package transforms

import (
	t "github.com/mholzen/information/triples"
)

func Cartesian(sets []*t.Triples) []*t.Triples {
	res := make([]*t.Triples, 0)

	if len(sets) == 0 {
		return res
	}

	first := sets[0]
	rest := sets[1:]

	restProducts := Cartesian(rest)

	for _, triple := range first.TripleSet {
		restTriples := t.NewTriples()
		restTriples.Add(triple)
		if len(restProducts) == 0 {
			res = append(res, restTriples)
		} else {
			for _, restProduct := range restProducts {
				restTriples.AddTriples(restProduct)
				res = append(res, restTriples)
			}
		}
	}
	return res
}
