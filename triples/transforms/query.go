package transforms

import (
	"log"

	t "github.com/mholzen/information/triples"
)

type Query struct {
	Query        *t.Triples
	Computations Computations // computations that need to match
}

func NewQuery(query *t.Triples, computations Computations) Query {
	return Query{
		Query:        query,
		Computations: computations,
	}
}

func NewQueryFromTriples(source *t.Triples) (Query, error) {
	computations := Computations{}
	query := t.NewTriples()
	res := Query{}
	for _, triple := range source.TripleSet {
		if _, ok := triple.Predicate.(t.UnaryFunctionNode); ok {
			c, err := NewComputationFromTriple(triple)
			if err != nil {
				return res, err
			}
			computations = append(computations, c)
		} else {
			query.Add(triple)
		}
	}
	return NewQuery(query, computations), nil
}

func (q Query) GetMatches(source *t.Triples) (Matches, error) {
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

func (q Query) GetSolutions(matches Matches) SolutionList {
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

func (q Query) Apply(source *t.Triples) (SolutionList, error) {
	matches, err := q.GetMatches(source)
	if err != nil {
		return SolutionList{}, err
	}

	solutions := q.GetSolutions(matches)
	solutions = solutions.FilterByComputations(q.Computations)

	return solutions, nil
}

func (q Query) GetTriples() *t.Triples {
	res := t.NewTriples()
	res.AddTriples(q.Query)
	res.AddTriples(q.Computations.GetTriples())
	return res
}

func NewQueryMapper(triples *t.Triples) (t.Mapper, error) {
	query := NewQuery(triples, Computations{})
	return func(source *t.Triples) (*t.Triples, error) {
		solutions, err := query.Apply(source)
		if err != nil {
			return nil, err
		}
		return solutions.GetAllTriples(), nil
	}, nil
}

func NewTripleQueryMatchMapper(query t.Triple) t.Mapper {
	matcher := NewTripleMatch(query)

	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, triple := range source.TripleSet {
			if matcher(triple) {
				log.Printf("triple %s matches %s", triple, query)
				res.Add(triple)
			}
		}
		return res, nil
	}
}

func NewNodeTester(node t.Node) t.NodeBoolFunction {
	switch n := node.(type) {
	case t.NodeBoolFunction:
		return n
	case VariableNode:
		return t.NodeMatchAny
	default:
		return func(n t.Node) bool {
			return node == n
		}
	}
}
func NewTripleMatch(query t.Triple) t.TripleMatch {
	subjectTester := NewNodeTester(query.Subject)
	predicateTester := NewNodeTester(query.Predicate)
	objectTester := NewNodeTester(query.Object)
	return func(triple t.Triple) bool {
		return subjectTester(triple.Subject) &&
			predicateTester(triple.Predicate) &&
			objectTester(triple.Object)
	}
}
