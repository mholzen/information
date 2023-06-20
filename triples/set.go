package triples

type NodeSet map[string]struct{}

func NewNodeSet() NodeSet {
	return make(NodeSet)
}

func (set NodeSet) Add(n Node) Node {
	if _, ok := set[n.String()]; ok {
		return n
	} else {
		set[n.String()] = struct{}{}
		return n
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
