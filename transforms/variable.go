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
		if t.NodeEquals(v, value) {
			return nil
		}
		return fmt.Errorf("variable '%s' already set to '%s'", value.String(), v.String())
	}
	m[variable] = value
	return nil
}

func (m VariableMap) Get(variable VariableNode) (t.Node, error) {
	if v, ok := m[variable]; !ok {
		return v, fmt.Errorf("variable '%s' not found", variable.String())
	} else {
		return v, nil
	}
}

func (m VariableMap) TestOrSetTriple(query t.Triple, value t.Triple) error {
	if v, ok := query.Subject.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Subject); err != nil {
			return t.SubjectPosition.WrapError(err)
		}
	}
	if v, ok := query.Predicate.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Predicate); err != nil {
			return t.PredicatePosition.WrapError(err)
		}
	}
	if v, ok := query.Object.(VariableNode); ok {
		if err := m.TestOrSet(v, value.Object); err != nil {
			return t.ObjectPosition.WrapError(err)
		}
	}
	return nil
}

func (m VariableMap) GetVariableOrNode(node t.Node) (t.Node, error) {
	if variable, ok := node.(VariableNode); ok {
		value, err := m.Get(variable)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return node, nil
}

func (m VariableMap) GetTriple(triple t.Triple) (t.Triple, error) {
	subject, err := m.GetVariableOrNode(triple.Subject)
	if err != nil {
		return triple, t.SubjectPosition.WrapError(err)
	}
	predicate, err := m.GetVariableOrNode(triple.Predicate)
	if err != nil {
		return triple, t.PredicatePosition.WrapError(err)
	}
	object, err := m.GetVariableOrNode(triple.Object)
	if err != nil {
		return triple, t.ObjectPosition.WrapError(err)
	}

	return t.NewTriple(
		subject,
		predicate,
		object), nil
}

func (m VariableMap) MeetsComputation(computations Computations) bool {
	return computations.Test(m)
}
