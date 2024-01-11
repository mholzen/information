package transforms

import (
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	t "github.com/mholzen/information/triples"
)

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
