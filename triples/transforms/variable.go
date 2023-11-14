package transforms

import (
	"log"

	t "github.com/mholzen/information/triples"
	"storj.io/common/uuid"
)

type VariableNode struct {
	t.CreatedNode[uuid.UUID]
}

func NewVariableNode() VariableNode {
	value, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	return VariableNode{
		CreatedNode: t.NewCreatedNode(value),
	}
}

type VariableList []VariableNode

func (v VariableList) Traverse(nodes t.NodeList) []t.NodeList {
	res := make([]t.NodeList, 0)
	indices := make([]int, len(v))
	i := 0
	for {
		if i < len(v) && indices[i] < len(nodes) {
			// output
			row := make(t.NodeList, 0)
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

func GetVariableList(nodes t.NodeList) VariableList {
	res := make(VariableList, 0)
	for _, node := range nodes {
		if variable, ok := node.(VariableNode); ok {
			res = append(res, variable)
		}
	}
	return res
}

func (v VariableList) GetNodeList() t.NodeList {
	res := make(t.NodeList, 0)
	for _, variable := range v {
		res = append(res, variable)
	}
	return res
}
