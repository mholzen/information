package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NewTraverse(start t.Node, filter t.TripleMatch, output *t.Triples) t.Transformer {
	visitedNodes := make(t.NodeSet)
	nodeQueue := make([]t.Node, 0)
	nodeQueue = append(nodeQueue, start)
	resultIndex := 0

	return func(source *t.Triples) error {
		dest := t.NewAnonymousNode()
		for len(nodeQueue) > 0 {
			node := nodeQueue[0]
			nodeQueue = nodeQueue[1:]

			for _, triple := range source.GetTripleListForSubject(node) {
				if !filter(triple) {
					continue
				}
				tripleReference := output.AddTripleReference(triple)
				output.NewTripleFromNodes(dest, t.NewIndexNode(resultIndex), tripleReference)
				resultIndex++

				if !visitedNodes.Contains(triple.Object) {
					visitedNodes.Add(triple.Object)
					nodeQueue = append(nodeQueue, triple.Object)
				}
			}
		}
		return nil
	}
}

type NodeToNodeList func(t.Node) t.NodeList

func NewNodeTraverse(start t.Node, next NodeToNodeList) t.Mapper {
	visitedNodes := make(t.NodeSet)
	nodeQueue := make(t.NodeList, 0)
	nodeQueue = append(nodeQueue, start)
	resultIndex := 0

	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		root := t.NewAnonymousNode()
		for len(nodeQueue) > 0 {
			node := nodeQueue[0]
			nodeQueue = nodeQueue[1:]

			// Output node
			res.NewTripleFromNodes(root, t.NewIndexNode(resultIndex), node) // TOTRY: consider using an array
			visitedNodes.Add(node)
			resultIndex++

			// Queue next nodes
			for _, nextNode := range next(node) {
				if !visitedNodes.Contains(nextNode) {
					nodeQueue = append(nodeQueue, nextNode)
				}
			}
		}
		return res, nil
	}
}
