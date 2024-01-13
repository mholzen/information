package transforms

import (
	t "github.com/mholzen/information/triples"
)

type Matches struct {
	Query      *Query
	MatchesMap MatchesMap
}

func NewMatches(query Query) Matches {
	return Matches{
		Query:      &query,
		MatchesMap: make(MatchesMap),
	}
}

func (m Matches) IsComplete() bool {
	for _, query := range m.Query.Query.TripleSet {
		if matches, ok := m.MatchesMap[query]; !ok || matches == nil {
			return false
		} else if matches.IsEmpty() {
			return false
		}
	}
	return true
}

type MatchesMap map[t.Triple]*t.Triples
