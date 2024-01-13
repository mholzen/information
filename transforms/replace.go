package transforms

import (
	"github.com/mholzen/information/triples"
)

func NewReplaceTransform(from, to triples.Node) triples.Transformer {
	return func(source *triples.Triples) error {
		for _, triple := range source.TripleSet {
			if triple.Subject == from {
				triple.Subject = to
			}
			if triple.Predicate == from {
				triple.Predicate = to
			}
			if triple.Object == from {
				triple.Object = to
			}
		}
		return nil
	}
}

func NewReplaceMapper(from, to triples.Node) triples.Mapper {
	return func(source *triples.Triples) (*triples.Triples, error) {
		res := triples.NewTriples()
		for _, triple := range source.TripleSet {
			if triple.Subject == from {
				triple.Subject = to
			}
			if triple.Predicate == from {
				triple.Predicate = to
			}
			if triple.Object == from {
				triple.Object = to
			}
			res.Add(triple)
		}
		return res, nil
	}
}

func NewReplaceNodesMapper(replacement NodeMap) triples.Mapper {
	return func(source *triples.Triples) (*triples.Triples, error) {
		res := triples.NewTriples()
		for _, triple := range source.TripleSet {
			if v, ok := replacement[triple.Subject.String()]; ok {
				triple.Subject = v
			}
			if v, ok := replacement[triple.Predicate.String()]; ok {
				triple.Predicate = v
			}
			if v, ok := replacement[triple.Object.String()]; ok {
				triple.Object = v
			}
			res.Add(triple)
		}
		return res, nil
	}
}

type NodeMap map[string]triples.Node

func NewNodeMap(from, to triples.NodeList) NodeMap {
	res := make(NodeMap)
	for i, node := range from {
		res[node.String()] = to[i]
	}
	return res
}
