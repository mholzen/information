package transforms

import (
	"fmt"
	"log"

	t "github.com/mholzen/information/triples"
	"storj.io/common/uuid"
)

type VariableNode struct {
	t.CreatedNode[uuid.UUID]
	Value *t.Node
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

func (v *VariableNode) TestOrSet(n t.Node) error {
	if v.Value != nil {
		if (*v.Value).String() == n.String() {
			return nil
		}
		return fmt.Errorf("variable already set to '%s'", (*v.Value).String())
	}
	v.Value = &n
	return nil
}

func (v *VariableNode) Clear() {
	v.Value = nil
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

func (v VariableList) Clear() {
	for _, variable := range v {
		variable.Clear()
	}
}

func TestOrSet(query, triple t.Triple) error {
	if v, ok := query.Subject.(VariableNode); ok {
		if err := v.TestOrSet(triple.Subject); err != nil {
			return err
		}
	}
	if v, ok := query.Predicate.(VariableNode); ok {
		if err := v.TestOrSet(triple.Predicate); err != nil {
			return err
		}
	}
	if v, ok := query.Object.(VariableNode); ok {
		if err := v.TestOrSet(triple.Object); err != nil {
			return err
		}
	}
	return nil
}

type VariableMap map[VariableNode]t.Node

func NewVariableMap(queryTriples t.TripleList) VariableMap {
	variables := GetVariableList(queryTriples.GetNodes().GetNodeList())
	res := make(VariableMap)
	for _, variable := range variables {
		res[variable] = nil
	}
	return res
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
		return fmt.Errorf("variable already set to '%s'", v.String())
	}
	m[variable] = value
	return nil
}

func (m VariableMap) TestOrSetTriple(query t.Triple, value t.Triple) error {
	if v, ok := query.Subject.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Subject); err != nil {
			return err
		}
	}
	if v, ok := query.Predicate.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Predicate); err != nil {
			return err
		}
	}
	if v, ok := query.Object.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Object); err != nil {
			return err
		}
	}
	return nil

}
