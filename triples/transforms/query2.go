package transforms

import (
	t "github.com/mholzen/information/triples"
)

type Query2 struct {
	Query        *t.Triples
	Computations Computations // computations that need to match
}

func NewQuery2(query *t.Triples, computations Computations) Query2 {
	return Query2{
		Query:        query,
		Computations: computations,
	}
}

func (q *Query2) GetMatches(source *t.Triples) (Matches, error) {
	res := NewMatches(q)

	for _, query := range q.Query.TripleSet {
		tripleFilter := NewTripleQueryMatchMapper(query)
		matches, err := source.Map(tripleFilter)
		if err != nil {
			return res, err
		}
		res.MatchesMap[query] = matches
	}
	return res, nil
}

func (q Query2) GetSolutions(matches Matches) SolutionList {
	// get list of keys from mm
	keys := q.Query.GetTripleList()

	// Cartesian product of mm
	res := make(TriplesList, 0)
	for _, query := range keys {
		res = append(res, matches.MatchesMap[query])
	}
	products := Cartesian(res)

	res1 := make(SolutionList, 0)
	for _, product := range products {
		m := NewSolution(q.Query)
		for i, query := range keys {
			err := m.Add(query, product[i])
			if err != nil {
				continue
			}
		}
		if m.IsComplete() && len(m.SolutionMap) == len(q.Query.TripleSet) {
			res1 = append(res1, m)
		}
	}
	return res1
}

func (q *Query2) Apply(source *t.Triples) (SolutionList, error) {
	matches, err := q.GetMatches(source)
	if err != nil {
		return SolutionList{}, err
	}

	solutions := q.GetSolutions(matches)
	solutions = solutions.FilterByComputations(q.Computations)

	return solutions, nil

}

func NewQuery2Mapper(triples *t.Triples) (t.Mapper, error) {
	query := NewQuery2(triples, Computations{})
	return func(source *t.Triples) (*t.Triples, error) {
		solutions, err := query.Apply(source)
		if err != nil {
			return nil, err
		}
		return solutions.GetAllTriples(), nil
	}, nil

}
