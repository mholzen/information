package triples

import "time"

type StringNode struct {
	Value   string    `json:"text"`
	Created time.Time `json:"created"`
}

func NewStringNode(value string) StringNode {
	if node, ok := stringNodes[value]; ok {
		return node
	} else {
		stringNodes[value] = StringNode{
			Value:   value,
			Created: time.Now(),
		}
		return stringNodes[value]
	}
}

func (n StringNode) String() string {
	return n.Value
}

func (n StringNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case StringNode:
		return n.String() < other.String()
	case IndexNode:
		return n.Created.Compare(other.Created) < 0
	case AnonymousNode:
		return n.Created.Compare(other.Created) < 0
	default:
		return false
	}
}

var Subject = NewStringNode("1-subject")
var Predicate = NewStringNode("2-predicate")
var Object = NewStringNode("3-object")
var Contains = NewStringNode("contains")
var Name = NewStringNode("name")

type StringNodes map[string]StringNode

var stringNodes StringNodes = make(StringNodes)

func MakeStringNodes(from NodeSet) StringNodes {
	res := make(StringNodes)
	for _, node := range from.GetNodeList() {
		if s, ok := node.(StringNode); ok {
			res[node.String()] = s
		}
	}
	return res
}

func (nodes StringNodes) GetStringList() []string {
	res := make([]string, 0)
	for _, node := range nodes {
		res = append(res, node.String())
	}
	return res
}
