package transforms

import "github.com/mholzen/information/triples"

type Node struct {
	Id   string `json:"id"`
	Node triples.Node
}

func NewNode(n triples.Node) Node {
	return Node{Id: n.String(), Node: n}
}

type Nodes struct {
	Nodes     []Node         `json:"nodes"`
	NodeIndex map[string]int `json:"-"`
}

func NewNodes() Nodes {
	return Nodes{
		Nodes:     make([]Node, 0),
		NodeIndex: make(map[string]int),
	}
}

func (n *Nodes) Id(node triples.Node) int {
	if _, ok := n.NodeIndex[node.String()]; !ok {
		n.Nodes = append(n.Nodes, NewNode(node))
		n.NodeIndex[node.String()] = len(n.Nodes) - 1
	}
	return n.NodeIndex[node.String()]
}

type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

type NodeLink struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type NodeLinkTransformer struct {
	Transformer triples.Transformer
	Result      NodeLink
}

func NewNodeLinkTransformer() *NodeLinkTransformer {
	res := NodeLinkTransformer{}
	res.Transformer = func(source *triples.Triples) error {
		nodes := NewNodes()
		res.Result = NodeLink{
			Nodes: make([]Node, 0),
			Links: make([]Link, 0),
		}
		for _, triple := range source.TripleSet {
			res.Result.Links = append(res.Result.Links, Link{
				nodes.Id(triple.Subject),
				nodes.Id(triple.Object),
			})
		}
		res.Result.Nodes = nodes.Nodes
		return nil
	}
	return &res
}
