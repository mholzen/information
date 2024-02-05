package triples

import "fmt"

type NodeSort struct {
	NodeList NodeList
	LessFunc func(i, j int) bool
}

func (s NodeSort) Len() int           { return len(s.NodeList) }
func (s NodeSort) Swap(i, j int)      { s.NodeList[i], s.NodeList[j] = s.NodeList[j], s.NodeList[i] }
func (s NodeSort) Less(i, j int) bool { return s.LessFunc(i, j) }

func NodeFloatValue(node Node) (float64, error) {
	switch node := node.(type) {
	case NumberNode:
		return node.Value, nil
	case FloatNode:
		return float64(node.Value), nil
	case IndexNode:
		return float64(node.Value), nil
	default:
		return 0, fmt.Errorf("expected NumberNode, FloatNode, or IndexNode, got %T", node)
	}
}

// Sort by types
// 1. BoolNode
// 2. IndexNode, NumberNode, FloatNode
// 3. StringNode
// 4. AnonymousNode

func LexicalTypeRank(node Node) (int, func(Node, Node) bool) {
	switch node.(type) {
	case BoolNode:
		return 0, func(i, j Node) bool {
			return i.(BoolNode).LessThan(j.(BoolNode))
		}
	case IndexNode, NumberNode, FloatNode:
		return 1, func(i, j Node) bool {
			iValue, err := NodeFloatValue(i)
			if err != nil {
				return true
			}
			jValue, err := NodeFloatValue(j)
			if err != nil {
				return true
			}
			return iValue < jValue
		}
	case StringNode:
		return 2, func(i, j Node) bool {
			return i.(StringNode).LessThan(j.(StringNode))
		}
	case AnonymousNode:
		return 3, func(i, j Node) bool {
			return i.(AnonymousNode).LessThan(j.(AnonymousNode))
		}
	default:
		return 4, func(i, j Node) bool {
			return i.LessThan(j)
		}
	}
}

type NodeSortFunc func(i, j Node) bool

func NodeLessLexical(i, j Node) bool {
	iRank, iLess := LexicalTypeRank(i)
	jRank, _ := LexicalTypeRank(j)
	if iRank < jRank {
		return true
	}
	if iRank > jRank {
		return false
	}
	return iLess(i, j)
}
