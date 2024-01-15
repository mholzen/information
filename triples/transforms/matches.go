package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NewMatchesMap() MatchesMap {
	return make(MatchesMap)
}

type MatchesMap map[t.Triple]*t.Triples
