package transforms

import (
	"fmt"

	t "github.com/mholzen/information/triples"
)

type NodeMatch func(node t.Node) (bool, error)

func NewNodeMatch(node t.Node) NodeMatch {
	switch node := node.(type) {
	case t.UnaryFunctionNode:
		return func(n t.Node) (bool, error) {
			v, err := node(n)
			if err != nil {
				return false, err
			}
			return v == t.NewNumberNode(1), nil
		}
	default:
		return func(n t.Node) (bool, error) {
			return n == node, nil
		}
	}
}

type TripleFMatch func(triple t.Triple) (bool, error)

func NewTripleFMatch(triple t.Triple) TripleFMatch {
	tripleFunction, tripleFunctionOk := triple.Predicate.(t.TripleFunctionNode)

	subjectMatch := NewNodeMatch(triple.Subject)
	predicateMatch := NewNodeMatch(triple.Predicate)
	objectMatch := NewNodeMatch(triple.Object)
	return func(t t.Triple) (bool, error) {
		// log.Printf("testing for %v (is function: %v)", triple, tripleFunctionOk)

		if tripleFunctionOk {
			object, err := tripleFunction(t)
			if err != nil {
				return false, err
			}
			// log.Printf("function %v returned %v", tripleFunction, object)
			return (object == triple.Object), nil
		}

		// log.Printf("does %v match %v", t, triple)
		s, err := subjectMatch(t.Subject)
		if err != nil {
			return false, err
		}
		if !s {
			return false, nil
		}

		p, err := predicateMatch(t.Predicate)
		if err != nil {
			return false, err
		}
		if !p {
			return false, nil
		}

		o, err := objectMatch(t.Object)
		if err != nil {
			return false, err
		}
		if !o {
			return false, nil
		}
		// log.Printf("yes")
		return true, nil
	}
}

func NewContains(triple t.Triple, dest *t.Triples) t.Transformer {
	match := NewTripleFMatch(triple)
	return func(source *t.Triples) error {
		for _, t := range source.TripleSet {
			m, err := match(t)
			if err != nil {
				return err
			}
			if !m {
				continue
			}
			// if we get here, we have a match
			// n := dest.AddTripleReference(t)
			dest.AddTripleReference(t)
			// dest.AddTriple(n, "contained", NewNumberNode(1))
		}
		return nil
	}

}

func NewContainsTriples(needles *t.Triples) t.Mapper {
	matches := make([]TripleFMatch, 0)
	for _, triple := range needles.TripleSet {
		matches = append(matches, NewTripleFMatch(triple))
	}
	return func(haystack *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, match := range matches {
			matched := false
			for _, t := range haystack.TripleSet {
				m, err := match(t)
				if err != nil {
					return nil, err
				}
				if !m {
					continue
				}
				res.Add(t)
				matched = true
				break
			}
			if !matched {
				return res, nil
			}
		}
		return res, nil
	}

}

func ContainsTriples(needles, haystack *t.Triples) (bool, error) {
	res, err := NewContainsTriples(needles)(haystack)
	if err != nil {
		return false, err
	}
	return len(res.TripleSet) == len(needles.TripleSet)*3, nil
}

func NewContainsMapper(triple t.Triple) t.Mapper {

	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		transformer := NewContains(triple, res)
		err := transformer(source)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func NewContainsOrComputeMapper(triple t.Triple, functions *t.Triples) t.Mapper {

	// do we have a function that computes this label?
	functionFinder := Filter(
		And(
			NewPredicateTripleMatch(ComputeNode),
			NewObjectTripleMatch(triple.Predicate),
		),
	)
	functions, functionFinderErr := functionFinder(functions)
	subjects := t.GetUnaryNodes(functions.GetSubjectList())

	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		finder := NewContains(triple, res)
		// logrus.Debugf("searching for %v", triple)
		err := finder(source)
		if err != nil {
			return nil, err
		}

		// if we found it, we're done
		if len(res.TripleSet) > 0 {
			// logrus.Debugf("found it immediately")
			return res, nil
		}

		// otherwise, compute it
		if functionFinderErr != nil {
			return nil, functionFinderErr
		}
		if len(subjects) == 0 {
			// log.Printf("no functions found for %v -- returning not found", triple.Predicate)
			return res, nil
		}
		if len(subjects) > 1 {
			return nil, fmt.Errorf("too many (%d) functions found for %v", len(subjects), triple.Predicate)
		}

		// logrus.Debugf("computing for it using %v", triple.Predicate)
		computer := ComputeTripleTransformer(triple.Subject, subjects[0], triple.Predicate)
		err = computer(source)
		if err != nil {
			return nil, err
		}

		// and try again
		err = finder(source)
		if err != nil {
			return nil, err
		}

		// if len(res.TripleSet) > 0 {
		// 	logrus.Debugf("found it after compute")
		// } else {
		// 	logrus.Debugf("not found after compute")
		// }

		return res, nil
	}
}

func NewMultiContainsOrComputeMapper(toFind *t.Triples, functions *t.Triples) t.Mapper {
	mappers := make([]t.Mapper, 0)
	for _, triple := range toFind.TripleSet {
		mappers = append(mappers, NewContainsOrComputeMapper(triple, functions))
	}
	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, mapper := range mappers {
			m, err := mapper(source)
			if err != nil {
				return nil, err
			}
			// if len(m.TripleSet) > 0 {
			// 	logrus.Debugf("matcher for %d found %d", i, len(m.TripleSet))
			// }
			res.AddTriples(m)
		}
		return res, nil
	}
}
