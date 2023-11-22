package transforms

import (
	"fmt"
	"log"

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

type TripleMatrix []t.TripleList

func (ta TripleMatrix) String() string {
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
	log.Printf("queryTriples: %s", queryTriples)

	return func(source *t.Triples) (*t.Triples, error) {
		// TODO: refactor into a pipe (to help debug)
		solutionsPerQueryTriple, err := SolutionsPerQueryTriple(queryTriples, source)
		if err != nil {
			return nil, err
		}
		log.Printf("solutionsPerQueryTriple: %s", solutionsPerQueryTriple)

		products := Cartesian(solutionsPerQueryTriple)
		log.Printf("products: %s", products)

		solutions := FilterByVariables(queryTriples, products)
		log.Printf("solutions: %s", solutions)

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

func FilterByVariables(queryTriples t.TripleList, products TripleMatrix) *t.Triples {
	res := t.NewTriples()
	root := t.NewAnonymousNode()

	variables := NewVariableMap(queryTriples)

	for i, solution := range products {

		variables.Clear()
		breakOuter := false
		for j, triple := range solution {
			if err := variables.TestOrSetTriple(queryTriples[j], triple); err != nil {
				log.Printf("skipping solution %d: %s", i, err)
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
func NewTripleMatch(query t.Triple) TripleMatch {
	subjectTester := NewNodeTester(query.Subject)
	predicateTester := NewNodeTester(query.Predicate)
	objectTester := NewNodeTester(query.Object)
	return func(triple t.Triple) bool {
		return subjectTester(triple.Subject) &&
			predicateTester(triple.Predicate) &&
			objectTester(triple.Object)
	}
}
