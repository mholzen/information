package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type TriplesList []*t.Triples

func (tl TriplesList) String() string {
	res := "[\n"

	for i, triples := range tl {
		res += fmt.Sprintf("%d: %s\n", i, triples.StringLine())
	}
	res += "]\n"
	return res
}

type TripleArray []t.TripleList

func (ta TripleArray) String() string {
	res := "[\n"

	for _, triples := range ta {
		res += fmt.Sprintf("[%s]\n", triples.StringLine())
	}
	res += "]\n"
	return res
}

func NewQueryMapper(query *t.Triples) t.Mapper {
	// for each triple in the query
	// find the set of triples that match it
	// then generate the cartesian product of those sets
	// then filter solutions to those where variables match
	queryTriples := query.GetTripleList()

	return func(source *t.Triples) (*t.Triples, error) {
		solutionsPerQueryTriple, err := SolutionsPerQueryTriple(queryTriples, source)
		if err != nil {
			return nil, err
		}

		products := Cartesian(solutionsPerQueryTriple)

		solutions := FilterByVariables(queryTriples, products)

		return solutions, nil
	}
}

func SolutionsPerQueryTriple(queryTriples t.TripleList, source *t.Triples) (TriplesList, error) {
	solutionsPerQueryTriple := make(TriplesList, 0)

	for _, triple := range queryTriples {
		tripleFilter := NewTripleQueryMatchMapper(triple)
		matches, err := source.Map(tripleFilter)
		if err != nil {
			return nil, err
		}
		solutionsPerQueryTriple = append(solutionsPerQueryTriple, matches)
	}
	return solutionsPerQueryTriple, nil
}

func FilterByVariables(queryTriples t.TripleList, products TripleArray) *t.Triples {
	res := t.NewTriples()
	root := t.NewAnonymousNode()

	variables := NewVariableMap(queryTriples)

	for i, solution := range products {

		variables.Clear()
		breakOuter := false
		for j, triple := range solution {
			if variables.TestOrSetTriple(queryTriples[j], triple) != nil {
				breakOuter = true
				break
			}
		}
		if breakOuter {
			continue
		}
		node := res.AddTripleReferences(t.NewTriplesFromList(solution))
		res.AddTriple(root, t.NewIndexNode(i), node)
	}
	return res
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
