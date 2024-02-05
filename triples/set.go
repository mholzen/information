package triples

import (
	"sort"
)

type NodeSet map[string]Node

func NewNodeSet() NodeSet {
	return make(NodeSet)
}

func (set NodeSet) Add(n Node) Node {
	if _, ok := set[n.String()]; ok {
		return n
	} else {
		set[n.String()] = n
		return n
	}
}

func (set NodeSet) Get(n string) Node {
	if v, ok := set[n]; ok {
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
		res = append(res, node)
	}
	return res
}

func (set NodeSet) GetSortedNodeList() NodeList {
	nodes := set.GetNodeList()
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].LessThan(nodes[j])
	})

	return nodes
}

func (set NodeSet) Intersect(with NodeSet) NodeSet {
	res := NewNodeSet()
	for _, node := range set {
		if with.Contains(node) {
			res.Add(node)
		}
	}
	return res
}

type NodeList []Node

func NewNodeList(values ...any) (NodeList, error) {
	res := make(NodeList, 0)
	for _, value := range values {
		node, err := NewNode(value)
		if err != nil {
			return nil, err
		}
		res = append(res, node)
	}
	return res, nil
}

func (list NodeList) Strings() []string {
	res := make([]string, 0)
	for _, node := range list {
		res = append(res, node.String())
	}
	return res
}

func (list NodeList) Sort() NodeList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].LessThan(list[j])
	})
	return list
}

func (list NodeList) SortBy(less NodeSortFunc) NodeList {
	sort.Slice(list, func(i, j int) bool {
		return less(list[i], list[j])
	})
	return list
}

func (list NodeList) SortLexical() NodeList {
	return list.SortBy(NodeLessLexical)
}
