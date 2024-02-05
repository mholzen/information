package transforms

import (
	"log"

	t "github.com/mholzen/information/triples"
)

type Query struct {
	Matching              *t.Triples
	Selected              *t.Triples
	Computations          Computations
	ComputationGenerators ComputationGenerators
}

func NewQuery(matching *t.Triples) Query {
	return Query{
		Matching:              matching,
		Selected:              t.NewTriples(),
		Computations:          Computations{},
		ComputationGenerators: ComputationGenerators{},
	}
}

func NewQueryFromTriples(source *t.Triples) (Query, error) {
	query := Query{}
	computations := Computations{}
	matching := t.NewTriples()
	for _, triple := range source.TripleSet {
		if _, ok := triple.Predicate.(t.UnaryFunctionNode); ok {
			c, err := NewComputationFromTriple(triple)
			if err != nil {
				return query, err
			}
			computations = append(computations, c)
		} else {
			matching.Add(triple)
		}
	}
	query.Computations = computations
	query.Matching = matching
	return query, nil
}

func (q Query) SearchForMatches(source *t.Triples) (MatchesMap, error) {
	res := NewMatchesMap()

	for _, query := range q.Matching.TripleSet {
		tripleFilter := NewTripleQueryMatchMapper(query)
		matches, err := source.Map(tripleFilter)
		if err != nil {
			return res, err
		}
		res[query] = matches
	}
	return res, nil
}

func (q Query) ComputeSolutions(matches MatchesMap) SolutionList {
	matching := q.Matching.GetTripleList()

	res := make(t.TriplesList, 0)
	for _, query := range matching {
		res = append(res, matches[query])
	}
	products := res.Cartesian()

	solutions := make(SolutionList, 0)
	for _, product := range products {
		// log.Printf("\n> evaluating solution:\n%s\n", product)
		solution := NewSolution(q.Matching)
		for i, match := range matching {
			err := solution.Add(match, product[i])
			if err != nil {
				// log.Printf("%s", err)
				continue
			}
		}
		if solution.IsComplete() && len(solution.SolutionMap) == len(q.Matching.TripleSet) {
			// log.Printf("solution passes matches")
			solutions = append(solutions, solution)
		}
	}
	return solutions
}

func (q Query) SearchForSolutions(source *t.Triples) (SolutionList, error) {
	matches, err := q.SearchForMatches(source)
	if err != nil {
		return SolutionList{}, err
	}
	log.Printf("matches count: %v", len(matches))

	solutions := q.ComputeSolutions(matches)
	log.Printf("solutions count: %v", len(solutions))

	computations := q.Computations.AugmentWithGenerators(q.ComputationGenerators, source)
	solutions = solutions.FilterByComputations(computations)
	log.Printf("solutions count after filtering: %v", len(solutions))

	return solutions, nil
}

func (q Query) SearchForSelected(source *t.Triples) (*t.Triples, error) {
	solutions, err := q.SearchForSolutions(source)
	if err != nil {
		return nil, err
	}
	return solutions.GetSelectTriples(q.Selected)
}

func (q Query) GetTriples() *t.Triples {
	res := t.NewTriples()
	res.AddTriples(q.Matching)
	res.AddTriples(q.Computations.GetTriples())
	return res
}

func (q Query) GetMapper() t.Mapper {
	return func(source *t.Triples) (*t.Triples, error) {
		res, err := q.SearchForSelected(source)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func NewQueryMapper(triples *t.Triples) (t.Mapper, error) {
	query := NewQuery(triples)
	return func(source *t.Triples) (*t.Triples, error) {
		solutions, err := query.SearchForSolutions(source)
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
