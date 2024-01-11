package transforms

import (
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	t "github.com/mholzen/information/triples"
)

type Query struct {
	QueryTriples t.TripleList // triples that need to match
	Computations Computations // computations that need to match

	SelectTriples t.TripleList // triples that are part of the response

	Query *t.Triples // source query
}

func NewQueryMapper(triples *t.Triples) (t.Mapper, error) {
	query, err := NewQueryFromTriples(triples)
	if err != nil {
		return nil, err
	}

	return func(source *t.Triples) (*t.Triples, error) {
		solutions, err := query.Apply(source)
		if err != nil {
			return nil, err
		}
		return solutions.Triples(), nil
	}, nil
}

func (q Query) IsSelect(triple t.Triple) bool {
	// TODO: should probably cache
	nodes := q.Query.GetTripleReferences(triple)
	if len(nodes) == 0 {
		return false
	}
	for _, n := range nodes {
		match := t.NewTripleFromNodes(t.NodeBoolFunction(t.NodeMatchAnyAnonymous), t.NewStringNode("select"), n)
		if len(q.QueryTriples.Filter(NewTripleMatch(match))) > 0 {
			return true
		}
	}
	return false
}

func (q Query) Solutions(source *t.Triples) (TriplesList, error) {
	solutionsPerQueryTriple := make(TriplesList, 0)

	for _, triple := range q.QueryTriples {
		tripleFilter := NewTripleQueryMatchMapper(triple)
		matches, err := source.Map(tripleFilter)
		if err != nil {
			return nil, err
		}
		solutionsPerQueryTriple = append(solutionsPerQueryTriple, matches)
	}
	return solutionsPerQueryTriple, nil
}

func (query Query) FilterByVariables(products TripleMatrix) (TripleMatrix, error) {
	res := make(TripleMatrix, 0)
	variables := NewVariableMapFromTripleList(query.QueryTriples)

	for i, solution := range products {

		variables.Clear()
		breakOuter := false

		for j, triple := range solution {
			if err := variables.TestOrSetTriple(query.QueryTriples[j], triple); err != nil {
				log.Printf("skipping solution %d: %s", i, err)
				breakOuter = true
				break
			}
		}
		if breakOuter {
			continue
		}
		// Passes variables -- now test computations
		log.Printf("variables:\n\t%s", variables)

		// ok, err := query.Computations.Test(variables, solution)
		// if err != nil {
		// 	return nil, err
		// }
		// if !ok {
		// 	continue
		// }

		res = append(res, solution)
	}
	return res, nil
}

func InitLog() {
	logFile, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	// defer logFile.Close()
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	spew.Config.Indent = "\t" // Set the indentation to a tab, for example
}

func (query Query) Apply(source *t.Triples) (TripleMatrix, error) {
	// for each triple in the query
	// find the set of triples that match it
	// then generate the cartesian product of those sets
	// then filter solutions to those where variables match

	log.Printf("query: %s", spew.Sdump(query))

	log.Printf("query\n%+v", query)
	log.Printf("source\n%+v", source)

	// TODO: refactor into a pipe (to help debug)
	solutionsPerQueryTriple, err := query.Solutions(source)
	if err != nil {
		return nil, err
	}
	log.Printf("solutionsPerQueryTriple: %s", solutionsPerQueryTriple)

	products := Cartesian(solutionsPerQueryTriple)
	log.Printf("products: %s", products)

	solutions, err := query.FilterByVariables(products)
	if err != nil {
		return nil, err
	}
	log.Printf("solutions: %s", solutions)

	return solutions, nil
}

func NewQueryFromTriples(query *t.Triples) (Query, error) {
	res := Query{}
	res.Query = query

	referenceConnected := ReferenceTriplesConnected(query)
	queryTriples, err := query.Map(NewIntersectMapper(referenceConnected))
	if err != nil {
		return res, err
	}
	res.QueryTriples = queryTriples.GetTripleList()

	computations, err := NewComputationsFromTriples(queryTriples)
	if err != nil {
		return res, err
	}
	res.Computations = computations
	return res, nil
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
