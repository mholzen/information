package transforms

import (
	t "github.com/mholzen/information/triples"
)

type Solution struct {
	VariableMap VariableMap
	SolutionMap SolutionMap
}

func NewSolution(query *t.Triples) Solution {
	return Solution{
		VariableMap: NewVariableMap(NewVariableListFromTriples(query)),
		SolutionMap: make(SolutionMap),
	}
}

func (s *Solution) Add(query t.Triple, triple t.Triple) error {
	if err := s.VariableMap.TestOrSetTriple(query, triple); err != nil {
		return err
	}
	s.SolutionMap[query] = triple
	return nil
}

func (s *Solution) IsComplete() bool {
	// TODO: should also check all queries have been matched
	return s.VariableMap.IsComplete()
}

func (s *Solution) GetTriple(query t.Triple) t.Triple {
	return s.SolutionMap[query]
}

func (s *Solution) GetAllTriples() *t.Triples {
	res := t.NewTriples()
	for query := range s.SolutionMap {
		res.Add(s.SolutionMap[query])
	}
	return res
}

func (s *Solution) MeetsComputation(computations Computations) bool {
	return s.VariableMap.MeetsComputation(computations)
}
