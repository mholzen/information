package triples

type Reducer (func(Node, Node) Node)

func SumReducer(a Node, b Node) Node {
	an, a_ok := a.(NumberNode)
	bn, b_ok := b.(NumberNode)
	if !a_ok || !b_ok {
		return nil
	}
	return NewNumberNode(an.Value + bn.Value)
}
