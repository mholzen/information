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

func (q Query) GetMatchesMap(source *t.Triples) (MatchesMap, error) {
	res := NewMatchesMap()

	for _, query := range q.Query.TripleSet {
		tripleFilter := NewTripleQueryMatchMapper(query)
		matches, err := source.Map(tripleFilter)
		if err != nil {
			return res, err
		}
		res[query] = matches
	}
	return res, nil
}

func (q Query) GetSolutions(matches MatchesMap) SolutionList {
	selectTriples := q.Query.GetTripleList()

	res := make(t.TriplesList, 0)
	for _, query := range selectTriples {
		res = append(res, matches[query])
	}
	products := res.Cartesian()

	solutions := make(SolutionList, 0)
	for _, product := range products {
		// log.Printf("\n> evaluating solution:\n%s\n", product)
		solution := NewSolution(q.Query)
		for i, query := range selectTriples {
			err := solution.Add(query, product[i])
			if err != nil {
				// log.Printf("%s", err)
				continue
			}
		}
		if solution.IsComplete() && len(solution.SolutionMap) == len(q.Query.TripleSet) {
			// log.Printf("solution passes matches")
			solutions = append(solutions, solution)
		}
	}
	return solutions
}

func (q Query) Apply(source *t.Triples) (SolutionList, error) {
	matches, err := q.GetMatchesMap(source)
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
