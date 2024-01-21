package transforms

import (
	t "github.com/mholzen/information/triples"
)

func Root(source *t.Triples) (t.Node, error) {
	reverse := func(node t.Node) t.NodeList {
		return source.GetTripleListForObject(node).GetSubjects()
	}
	start := source.GetSubjectList()[0]
	rootMapper := NewNodeTraverse(start, reverse)
	res, err := source.Map(rootMapper)
	if err != nil {
		return nil, err
	}
	list := res.GetTripleList()
	if len(list) == 0 {
		return nil, nil
	}
	return list[len(list)-1].Object, nil
}

func RootMapper(source *t.Triples) (*t.Triples, error) {
	root, err := Root(source)
	if err != nil {
		return nil, err
	}
	res := t.NewTriples()
	res.AddTripleFromAny(t.NewAnonymousNode(), t.NewStringNode("root"), root)
	return res, nil
}

func RootNode(source *t.Triples) (t.Node, error) {

	traversed := make(t.NodeSet)
	if len(source.GetTripleList()) == 0 {
		return nil, nil
	}
	root := source.GetTripleList()[0].Subject
	for {
		sources := source.GetTriplesForObject(root)
		if len(sources.GetTripleList()) == 0 {
			break
		}

		traversed.Add(root)
	}
	return root, nil

}
