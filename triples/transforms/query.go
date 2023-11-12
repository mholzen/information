package transforms

import (
	"log"

	t "github.com/mholzen/information/triples"
	"storj.io/common/uuid"
)

type VariableNode struct {
	t.CreatedNode[uuid.UUID]
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
func NewQueryMapper(query *t.Triples) t.Mapper {
	// for each triple in the query
	// find the set of triples that match it
	// then generate the cartesian product of those sets

	return func(source *t.Triples) (*t.Triples, error) {
		solutions := make([]*t.Triples, 0)
		for _, triple := range query.TripleSet {
			tripleFilter := NewTripleQueryMapper(triple)
			matches, err := source.Map(tripleFilter)
			if err != nil {
				return nil, err
			}
			solutions = append(solutions, matches)
		}
		products := Cartesian(solutions)

		res := t.NewTriples()
		root := t.NewAnonymousNode()
		for i, triples := range products {
			node := res.AddTripleReferences(triples)
			res.AddTriple(root, t.NewIndexNode(i), node)
		}
		return res, nil
	}
}

func NewTripleQueryMapper(query t.Triple) t.Mapper {
	return func(source *t.Triples) (*t.Triples, error) {
		res := t.NewTriples()
		for _, triple := range source.TripleSet {
			if query == triple {
				res.Add(triple)
			}
		}
		return res, nil
	}
}

// The results can contain multiple solutions, each identified by the predicate "solution" and the object being an index node.
func NewQueryTransformerWithDefinitions(query, dest, definitions *t.Triples) t.Transformer {
	return func(source *t.Triples) error {
		// find list of variables from query
		variables := GetVariableList(query.Nodes.GetNodeList())

		valuesList := variables.Traverse(source.Nodes.GetNodeList())

		// log.Printf("source is:\n%v", source.String())

		for solutionNo, values := range valuesList {
			nodeMap := NewNodeMap(variables.GetNodeList(), values)

			// log.Printf("evaluation solution %d with nodemap: %+v", solutionNo, maps.Values(nodeMap))

			instantiatedQuery, err := query.Map(NewReplaceNodesMapper(nodeMap))
			if err != nil {
				return err
			}

			// compute any needed triples
			// err = NewComputeWithDefinitions(definitions)(instantiatedQuery)
			// if err != nil {
			// 	return err
			// }

			// filter := NewContainsTriples(instantiatedQuery)
			filter := NewMultiContainsOrComputeMapper(instantiatedQuery, definitions) // TODO: must support multiple triples

			matches, err := filter(source)
			if err != nil {
				return err
			}
			// logrus.Debugf("=== matches:\n%v", matches)
			// log.Printf("len(matches):\n%v", len(matches.TripleSet))
			// log.Printf("len(query):\n%v", len(query.TripleSet))
			if len(matches.TripleSet) == (len(query.TripleSet) * 3) { // TODO: magic number
				// log.Printf("found match")
				container := t.NewAnonymousNode()
				dest.AddTriple(container, t.NewStringNode("solution"), t.NewIndexNode(solutionNo))

				// Add NodeMap to dest
				// TODO: refactor to method of NodeMap
				for variable, value := range nodeMap {
					dest.AddTriple(container, variable, value)
				}

				log.Printf("solution %d contains:\n%s", solutionNo, matches.String())
				// dest.AddTriples(matches)
				for _, triple := range matches.TripleSet {
					dest.Add(triple) // should already be a reference
					dest.AddTriple(container, t.NewStringNode("contains"), triple.Subject)
				}
			}
		}
		log.Printf("results:\n%s", dest.String())
		return nil
	}
}

func (v VariableList) GetNodeList() t.NodeList {
	res := make(t.NodeList, 0)
	for _, variable := range v {
		res = append(res, variable)
	}
	return res
}
