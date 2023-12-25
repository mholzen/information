package triples

type TripleMatch func(triple Triple) bool

type TripleMatchError func(triple Triple) (bool, error)

func (source *Triples) Filter(match TripleMatch) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if match(triple) {
			res.Add(triple)
		}
	}
	return res
}
