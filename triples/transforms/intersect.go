package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NewIntersectMapper(difference *t.Triples) t.Mapper {
	contained := NotMatch(NewContainedTripleMatch(difference))
	return func(source *t.Triples) (*t.Triples, error) {
		return source.Filter(contained), nil
	}
}
