package transforms

import t "github.com/mholzen/information/triples"

type SolutionList []Solution

func (sl SolutionList) GetTriples(query t.Triple) *t.Triples {
	res := t.NewTriples()
	for _, solution := range sl {
		res.Add(solution.GetTriple(query))
	}
	return res
}

func (sl SolutionList) GetAllTriples() *t.Triples {
	res := t.NewTriples()
	for _, solution := range sl {
		res.AddTriples(solution.GetAllTriples())
	}
	return res
}

func (sl SolutionList) FilterByComputations(computations Computations) SolutionList {
	res := make(SolutionList, 0)
	for _, solution := range sl {
		if solution.MeetsComputation(computations) {
			res = append(res, solution)
		}
	}
	return res
}

type SolutionMap map[t.Triple]t.Triple

func (sm SolutionMap) Query() *t.Triples {
	res := t.NewTriples()
	for query := range sm {
		res.Add(query)
	}
	return res
}

func (sm SolutionMap) GetVariableMap() VariableMap {
	return NewVariableMap(NewVariableListFromTriples(sm.Query()))
}

func (sm SolutionMap) TestVariables() bool {
	variables := sm.GetVariableMap()
	for query, triple := range sm {
		if err := variables.TestOrSetTriple(query, triple); err != nil {
			return false
		}
	}
	return true
}
