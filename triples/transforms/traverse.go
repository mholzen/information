package transforms

import (
	. "github.com/mholzen/information/triples"
	"github.com/sirupsen/logrus"
)

func NewTraverse(start Node, filter TripleMatch, dest Node, output *Triples) Transformer {
	visitedNodes := make(NodeSet)
	nodeQueue := make([]Node, 0)
	nodeQueue = append(nodeQueue, start)
	resultIndex := 0

	return func(source *Triples) error {
		for len(nodeQueue) > 0 {
			node := nodeQueue[0]
			nodeQueue = nodeQueue[1:]

			for _, triple := range source.GetTriplesForSubject(node) {
				if !filter(triple) {
					logrus.Debugf("%s fail", triple)
					continue
				}
				logrus.Debugf("%s pass", triple)
				tripleReference := output.AddTripleReference(triple)
				output.NewTripleFromNodes(dest, NewIndexNode(resultIndex), tripleReference)
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
