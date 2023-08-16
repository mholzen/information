package triples

type NodeSet map[string]*Node

func NewNodeSet() NodeSet {
	return make(NodeSet)
}

func (set NodeSet) Add(n Node) Node {
	if _, ok := set[n.String()]; ok {
		return n
	} else {
		set[n.String()] = &n
		return n
	}
}

func (set NodeSet) Get(n Node) *Node {
	if v, ok := set[n.String()]; ok {
		return v
	} else {
		return nil
	}
}

func (set NodeSet) Contains(node Node) bool {
	_, ok := set[node.String()]
	return ok
}

func (set NodeSet) ContainsOrAdd(node Node) bool {
	if set.Contains(node) {
		return true
	} else {
		set.Add(node)
		return false
	}
}

func (set NodeSet) GetNodeList() NodeList {
	res := make(NodeList, 0)
	for _, node := range set {
		res = append(res, *node)
	}
	return res
}

// func (set NodeSet) GetVariableList() VariableList {
// 	res := make(VariableList, 0)
// 	for _, node := range set {
// 		if variable, ok := (*node).(VariableNode); ok {
// 			res = append(res, variable)
// 		}
// 	}
// 	return res
// }

type NodeList []Node
