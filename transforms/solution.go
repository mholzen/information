package transforms

import (
	t "github.com/mholzen/information/triples"
)

type Solution struct {
	VariableMap VariableMap
	SolutionMap SolutionMap
}

func NewSolution(match *t.Triples) Solution {
	return Solution{
		VariableMap: NewVariableMap(NewVariableListFromTriples(match)),
		SolutionMap: make(SolutionMap),
	}
}

func (s *Solution) Add(match t.Triple, solution t.Triple) error {
	if err := s.VariableMap.TestOrSetTriple(match, solution); err != nil {
		return err
	}
	s.SolutionMap[match] = solution
	return nil
}

func (s *Solution) IsComplete() bool {
	return s.VariableMap.IsComplete()
}

func (s *Solution) GetMatching(triple t.Triple) t.Triple {
	return s.SolutionMap[triple]
}

func (s *Solution) GetSelected(triple t.Triple) (t.Triple, error) {
	return s.VariableMap.GetTriple(triple)
}

func (s *Solution) GetAllMatching() *t.Triples {
	res := t.NewTriples()
	for query := range s.SolutionMap {
		res.Add(s.SolutionMap[query])
	}
	return res
}

func (s *Solution) GetSelectTriples(selected *t.Triples) (*t.Triples, error) {
	res := t.NewTriples()
	for _, selected := range selected.TripleSet {
		triple, err := s.GetSelected(selected)
		if err != nil {
			return nil, err
		}
		res.Add(triple)
	}
	return res, nil
}

func (s *Solution) MeetsComputation(computations Computations) bool {
	return s.VariableMap.MeetsComputation(computations)
}
