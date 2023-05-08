package triples

import (
	"fmt"
	"sort"
)

type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

func (t Triple) String() string {
	return fmt.Sprintf("(%s %s %s)", t.Subject, t.Predicate, t.Object)
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

			triples.NewTriple(subject, predicate, object)
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

			triples.NewTriple(container, predicate, object)
		}
	}
	return res, nil
}

type TripleList []Triple
type TripleSet map[Triple]struct{}

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

var triples *Triples = NewTriples()

func (source *Triples) NewTriple(subject Node, predicate Node, object Node) Triple {
	triple := Triple{subject, predicate, object}
	source.Add(triple)
	return triple
}

func (source *Triples) Add(triple Triple) {
	if _, ok := source.TripleSet[triple]; ok {
		return
	}
	source.TripleSet[triple] = struct{}{}

	source.Nodes.Add(triple.Subject)
	source.Nodes.Add(triple.Predicate)
	if triple.Object != nil {
		source.Nodes.Add(triple.Object)
	}
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

func (source *Triples) AddReachableTriples(node Node, triples *Triples) *Triples {
	if triples == nil {
		triples = NewTriples()
	}
	if triples.Nodes.ContainsOrAdd(node) {
		return triples
	}

	// log.Printf("len(source.Triples): %d", len(source.Triples))
	for triple := range source.TripleSet {
		// log.Printf("searching triple: %+v", triple)
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
	for triple := range source.TripleSet {
		tripleList = append(tripleList, triple)
	}
	return tripleList
}

func (s TripleList) Sort() {
	sort.Sort(TripleSort{s, func(i, j int) bool {
		return s[i].String() < s[j].String()
	}})
}

type TripleSort struct {
	data     TripleList
	lessFunc func(i, j int) bool
}

func (s TripleSort) Len() int           { return len(s.data) }
func (s TripleSort) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
func (s TripleSort) Less(i, j int) bool { return s.lessFunc(i, j) }

func (source *Triples) GetTriplesForSubject(node Node, triples *Triples) TripleList {
	res := make(TripleList, 0)
	for triple := range source.TripleSet {
		if triple.Subject.String() == node.String() {
			res = append(res, triple)
		}
	}
	sort.Sort(TripleSort{res, func(i, j int) bool {
		return res[i].Predicate.String() < res[j].Predicate.String()
	}})
	return res
}

func (source *Triples) String(triples TripleList, prefix string, depth int) string {
	res := ""
	for _, triple := range triples {
		res += fmt.Sprintf("%s%s %s\n", prefix, triple.Predicate, triple.Object)
		if _, ok := triple.Object.(AnonymousNode); ok {
			if depth > 0 {
				r := source.GetTriplesForSubject(triple.Object, nil)
				res += source.String(r, prefix+"    ", depth-1)
			}
		}
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

	for triple := range source.TripleSet {
		if triple.Object != nil && triple.Object.String() == node.String() {
			triples.Add(triple)
		}
	}

	return triples
}
