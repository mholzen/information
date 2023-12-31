package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NewObjectAugmenter(match t.TripleMatch, predicate, object t.Node) t.Transformer {
	return func(source *t.Triples) error {
		for _, triple := range source.TripleSet {
			if match(triple) {
				source.AddTripleFromNodes(triple.Object, predicate, object)
			}
		}
		return nil
	}
}
