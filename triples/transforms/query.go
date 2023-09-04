package transforms

import (
	"log"

	. "github.com/mholzen/information/triples"
	"golang.org/x/exp/maps"
	"storj.io/common/uuid"
)

type VariableNode struct {
	CreatedNode[uuid.UUID]
}

func NewVariableNode() VariableNode {
	value, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	return VariableNode{
		CreatedNode: NewCreatedNode(value),
	}
}

type VariableList []VariableNode

func (v VariableList) Traverse(nodes NodeList) []NodeList {
	res := make([]NodeList, 0)
	indices := make([]int, len(v))
	i := 0
	for {
		if i < len(v) && indices[i] < len(nodes) {
			// output
			row := make(NodeList, 0)
			for _, j := range indices {
				row = append(row, nodes[j])
			}
			res = append(res, row)
		} else {
			break
		}
		// advance
		{
			j := i
			for {
				if j == len(v) {
					break
				}
				indices[j]++
				if indices[j] == len(nodes) {
					indices[j] = 0
					j++
				} else {
					break
				}
			}
			if j == len(v) {
				break
			}
		}
	}
	return res
}

func GetVariableList(nodes NodeList) VariableList {
	res := make(VariableList, 0)
	for _, node := range nodes {
		if variable, ok := node.(VariableNode); ok {
			res = append(res, variable)
		}
	}
	return res
}

// The results can contain multiple solutions, each identified by the predicate "solution" and the object being an index node.
func NewQueryTransformer(query, dest, definitions *Triples) Transformer {
	return func(source *Triples) error {
		// find list of variables from query
		variables := GetVariableList(query.Nodes.GetNodeList())

		valuesList := variables.Traverse(source.Nodes.GetNodeList())

		// log.Printf("source is:\n%v", source.String())

		for solutionNo, values := range valuesList {
			nodeMap := NewNodeMap(variables.GetNodeList(), values)

			log.Printf("evaluation solution %d with nodemap: %+v", solutionNo, maps.Values(nodeMap))

			instantiatedQuery, err := query.Map(NewReplaceNodesMapper(nodeMap))
			if err != nil {
				return err
			}

			// compute any needed triples
			// err = NewComputeWithDefinitions(definitions)(instantiatedQuery)
			// if err != nil {
			// 	return err
			// }

			// filter := NewContainsTriples(instantiatedQuery)
			filter := NewMultiContainsOrComputeMapper(instantiatedQuery, definitions) // TODO: must support multiple triples

			matches, err := filter(source)
			if err != nil {
				return err
			}
			// logrus.Debugf("=== matches:\n%v", matches)
			// log.Printf("len(matches):\n%v", len(matches.TripleSet))
			// log.Printf("len(query):\n%v", len(query.TripleSet))
			if len(matches.TripleSet) == (len(query.TripleSet) * 3) { // TODO: magic number
				// log.Printf("found match")
				container := NewAnonymousNode()
				dest.AddTriple(container, NewStringNode("solution"), NewIndexNode(solutionNo))

				// Add NodeMap to dest
				// TODO: refactor to method of NodeMap
				for variable, value := range nodeMap {
					dest.AddTriple(container, variable, value)
				}

				log.Printf("solution %d contains:\n%s", solutionNo, matches.String())
				// dest.AddTriples(matches)
				for _, triple := range matches.TripleSet {
					dest.Add(triple) // should already be a reference
					dest.AddTriple(container, NewStringNode("contains"), triple.Subject)
				}
			}
		}
		log.Printf("results:\n%s", dest.String())
		return nil
	}
}

func (v VariableList) GetNodeList() NodeList {
	res := make(NodeList, 0)
	for _, variable := range v {
		res = append(res, variable)
	}
	return res
}
