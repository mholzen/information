package transforms

import (
	"log"

	. "github.com/mholzen/information/triples"
	"storj.io/common/uuid"
)

type VariableNode struct {
	CreatedNode[uuid.UUID]
	Value *Node
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

func NewQuery(query, dest *Triples) Transformer {
	return func(source *Triples) error {
		// find list of variables from query
		variables := GetVariableList(query.Nodes.GetNodeList())

		valuesList := variables.Traverse(source.Nodes.GetNodeList())

		// log.Printf("source is:\n%v", source.String())

		for solutionNo, values := range valuesList {
			nodeMap := NewNodeMap(variables.GetNodeList(), values)

			// log.Printf("evaluation solution %d with nodemap: %+v", solutionNo, nodeMap)

			res, err := query.Map(NewReplaceNodesMapper(nodeMap))

			if err != nil {
				return err
			}
			// log.Printf("solution %d is:\n%+v", solutionNo, res)

			if source.ContainsTriples(res) {
				container := NewAnonymousNode()
				dest.AddTriple(container, NewStringNode("solution"), NewIndexNode(solutionNo))
				for _, triple := range res.TripleSet {
					node := dest.AddTripleReference(triple)
					dest.AddTriple(container, NewStringNode("contains"), node)
				}
			}
		}
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
