package transforms

import (
	t "github.com/mholzen/information/triples"
)

// type QueryMatches struct {
// 	Queries *t.Triples
// }

// func NewQueryMatches(queries *t.Triples) *QueryMatches {
// 	return &QueryMatches{
// 		Queries: queries,
// 	}
// }

type Matches struct {
	Query      *Query2
	MatchesMap MatchesMap
}

func NewMatches(query *Query2) Matches {
	return Matches{
		Query:      query,
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

// func (mm MatchesMap) Query() *t.Triples {
// 	res := t.NewTriples()
// 	for query := range mm {
// 		res.Add(query)
// 	}
// 	return res
// }

// func (mm MatchesMap) Solutions() SolutionList {
// 	// get list of keys from mm
// 	keys := make([]t.Triple, 0)
// 	for k := range mm {
// 		keys = append(keys, k)
// 	}

// 	// Cartesian product of mm
// 	res := make(TriplesList, 0)
// 	for _, query := range keys {
// 		res = append(res, mm[query])
// 	}
// 	products := Cartesian(res)

// 	res1 := make(SolutionList, 0)
// 	for _, product := range products {
// 		m := NewSolution(mm.Query())
// 		for i, query := range keys {
// 			err := m.Add(query, product[i])
// 			if err != nil {
// 				continue
// 			}
// 		}
// 		if m.IsComplete() {
// 			res1 = append(res1, m)
// 		}
// 	}
// 	return res1
// }
