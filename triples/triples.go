package triples

import (
	"fmt"
)

type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

func (t Triple) String() string {
	return fmt.Sprintf("(%s %s %s)", NodeString(t.Subject), NodeString(t.Predicate), NodeString(t.Object))
}

type TripleList []Triple
type TripleSet map[string]Triple

type Triples struct {
	TripleSet TripleSet
	Nodes     NodeSet
}

func NewTriples() *Triples {
	return &Triples{
		TripleSet: make(TripleSet),
		Nodes:     NewNodeSet(),
	}
}

func (source *Triples) NewTripleFromNodes(subject Node, predicate Node, object Node) Triple {
	triple := Triple{subject, predicate, object}
	source.Add(triple)
	return triple
}

func (source *Triples) NewTriple(subject, predicate, object any) (Triple, error) {
	var err error
	var s, p, o Node
	s, err = source.NewNode(subject)
	if err != nil {
		return Triple{}, err
	}
	p, err = source.NewNode(predicate)
	if err != nil {
		return Triple{}, err
	}
	o, err = source.NewNode(object)
	if err != nil {
		return Triple{}, err
	}

	triple := Triple{s, p, o}
	source.Add(triple)
	return triple, nil
}

func (source *Triples) NewTripleString(subject string, predicate string, object string) Triple {
	triple := Triple{NewStringNode(subject), NewStringNode(predicate), NewStringNode(object)}
	source.Add(triple)
	return triple
}

func (source *Triples) Add(triple Triple) {
	if _, ok := source.TripleSet[triple.String()]; ok {
		return
	}
	source.TripleSet[triple.String()] = triple

	source.Nodes.Add(triple.Subject)
	source.Nodes.Add(triple.Predicate)
	if triple.Object != nil {
		source.Nodes.Add(triple.Object)
	}
}

func (source *Triples) NewTriplesFromMap(m map[string]interface{}) (TripleList, error) {
	res := make(TripleList, 0)
	is_spo_form := false
	is_po_form := false
	container := NewAnonymousNode()
	for key, val := range m {
		if key == "s" || key == "p" || key == "o" {
			if is_spo_form {
				continue
			}
			if is_po_form {
				return res, fmt.Errorf("cannot mix spo and predicate/object form (m: %+v)", m)
			}
			subject, err := source.NewNode(m["s"])
			if err != nil {
				return res, err
			}
			predicate, err := source.NewNode(m["p"])
			if err != nil {
				return res, err
			}
			object, err := source.NewNode(m["o"])
			if err != nil {
				return res, err
			}

			source.NewTripleFromNodes(subject, predicate, object)
			is_spo_form = true
		} else {
			is_po_form = true
			if is_spo_form {
				return res, fmt.Errorf("cannot mix spo and predicate/object form (m: %+v)", m)
			}
			// create triples dependant on type
			predicate, err := source.NewNode(key)
			if err != nil {
				return res, err
			}
			object, err := source.NewNode(val)
			if err != nil {
				return res, err
			}

			source.NewTripleFromNodes(container, predicate, object)
		}
	}
	return res, nil
}

func (source *Triples) NewTriplesFromSlice(slice []interface{}) (TripleList, error) {
	triples := make(TripleList, 0)
	triple := Triple{}
	var err error
	triple.Subject, err = source.NewNode(slice[0])
	if err != nil {
		return triples, err
	}
	slice = slice[1:]

	triple.Predicate, err = source.NewNode(slice[0])
	if err != nil {
		return triples, err
	}
	slice = slice[1:]

	if len(slice) == 0 {
		source.Add(triple)
		triples = append(triples, triple)
	} else {
		for len(slice) > 0 {
			triple.Object, err = source.NewNode(slice[0])
			if err != nil {
				return triples, err
			}
			source.Add(triple)
			triples = append(triples, triple)
			slice = slice[1:]
		}
	}
	return triples, nil
}

func (source *Triples) Contains(triple Triple) bool {
	return source.TripleSet[triple.String()] == triple
}

func (source *Triples) GetTriplesForSubject(node Node, triples *Triples) TripleList {
	res := make(TripleList, 0)
	for _, triple := range source.TripleSet {
		if triple.Subject.String() == node.String() {
			res = append(res, triple)
		}
	}
	res.Sort()
	return res
}

func (source *Triples) String() string {
	res := ""
	for _, triple := range source.TripleSet {
		res += fmt.Sprintf("%s\n", triple)
	}
	return res
}

func (source *Triples) GetTriplesForObject(node Node, triples *Triples) *Triples {
	if triples == nil {
		triples = NewTriples()
	}
	if triples.Nodes.ContainsOrAdd(node) {
		return triples
	}

	for _, triple := range source.TripleSet {
		if triple.Object != nil && triple.Object.String() == node.String() {
			triples.Add(triple)
		}
	}

	return triples
}

func (source *Triples) GetTriples(subject, predicate, object Node) *Triples {
	return nil
}

func (source *Triples) AddReachableTriples(node Node, triples *Triples) *Triples {
	if triples == nil {
		triples = NewTriples()
	}
	if triples.Nodes.ContainsOrAdd(node) {
		return triples
	}

	for _, triple := range source.TripleSet {
		if triple.Subject.String() == node.String() ||
			triple.Predicate.String() == node.String() ||
			(triple.Object != nil && triple.Object.String() == node.String()) {

			if _, ok := triple.Subject.(AnonymousNode); ok {
				// log.Printf("found subject nil node: %+v", n)
				source.AddReachableTriples(triple.Subject, triples)
				if triple.Object != nil {
					source.AddReachableTriples(triple.Object, triples)
				}
			}
			if _, ok := triple.Object.(AnonymousNode); ok {
				// log.Printf("found object nil node: %+v", n)
				source.AddReachableTriples(triple.Subject, triples)
				if triple.Object != nil {
					source.AddReachableTriples(triple.Object, triples)
				}
			}
			triples.Add(triple)
		}
	}

	return triples
}
func (source *Triples) GetTripleList() TripleList {
	tripleList := make(TripleList, 0)
	for _, triple := range source.TripleSet {
		tripleList = append(tripleList, triple)
	}
	return tripleList
}

func (source *Triples) Compute() error {
	// generate new triples based on predicates that are functions
	for _, triple := range source.TripleSet {
		if p, ok := triple.Predicate.(UnaryFunctionNode); ok {
			// find triples describing this node

			object, err := p(triple.Object)
			if err != nil {
				return err
			}
			predicate := triple.Predicate.String()

			source.NewTripleFromNodes(triple.Subject, NewStringNode(predicate), object)
		}
	}
	return nil
}
