package transforms

import (
	"fmt"
	"log"

	t "github.com/mholzen/information/triples"
	"storj.io/common/uuid"
)

type VariableNode struct {
	t.CreatedNode[uuid.UUID]
}

func NewVariableNode() VariableNode { // CONSIDER: can named variables in a query be handled by a triple?
	value, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	return VariableNode{
		CreatedNode: t.NewCreatedNode(value),
	}
}

var Var = NewVariableNode

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

func NewVariableList(nodes t.NodeList) VariableList {
	res := make(VariableList, 0)
	for _, node := range nodes {
		if variable, ok := node.(VariableNode); ok {
			res = append(res, variable)
		}
	}
	return res
}

func NewVariableListFromTriples(triples *t.Triples) VariableList {
	return NewVariableList(triples.Nodes.GetNodeList())
}

func (v VariableList) GetNodeList() t.NodeList {
	res := make(t.NodeList, 0)
	for _, variable := range v {
		res = append(res, variable)
	}
	return res
}

type VariableMap map[VariableNode]t.Node

func (v VariableMap) IsComplete() bool {
	for _, value := range v {
		if value == nil {
			return false
		}
	}
	return true
}

func NewVariableMap(variables VariableList) VariableMap {
	res := make(VariableMap)
	for _, variable := range variables {
		res[variable] = nil
	}
	return res
}

func NewVariableMapFromTripleList(queryTriples t.TripleList) VariableMap {
	variables := NewVariableList(queryTriples.GetNodes().GetNodeList())
	return NewVariableMap(variables)
}

func (v VariableMap) Clear() {
	for k := range v {
		v[k] = nil
	}
}

func (m VariableMap) TestOrSet(variable VariableNode, value t.Node) error {
	if v, ok := m[variable]; ok && v != nil {
		if v.String() == value.String() {
			return nil
		}
		return fmt.Errorf("variable '%s' already set to '%s'", value.String(), v.String())
	}
	m[variable] = value
	return nil
}

func wrapError(position t.NodePosition, err error) error {
	return fmt.Errorf("error setting %s: %s", string(position), err)
}

func (m VariableMap) TestOrSetTriple(query t.Triple, value t.Triple) error {
	if v, ok := query.Subject.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Subject); err != nil {
			return wrapError(t.Subject1, err)
		}
	}
	if v, ok := query.Predicate.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Predicate); err != nil {
			return wrapError(t.Predicate1, err)
		}
	}
	if v, ok := query.Object.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Object); err != nil {
			return wrapError(t.Object1, err)
		}
	}
	return nil
}

func (v VariableMap) MeetsComputation(computations Computations) bool {
	return computations.Test(v)
}
