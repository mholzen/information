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
	return fmt.Sprintf("(%s, %s, %s)", shortString(t.Subject), shortString(t.Predicate), shortString(t.Object))
}

type NodePosition byte

const (
	Subject1 NodePosition = iota
	Predicate1
	Object1
)

func (t Triple) GetNode(position Node) (Node, error) {
	switch position {
	case Subject:
		return t.Subject, nil
	case Predicate:
		return t.Predicate, nil
	case Object:
		return t.Object, nil
	default:
		return nil, fmt.Errorf("invalid node position %s", position)
	}
}

func GetNodeFunction(position Node) (func(Triple) Node, error) {
	switch position {
	case Subject:
		return func(t Triple) Node { return t.Subject }, nil
	case Predicate:
		return func(t Triple) Node { return t.Predicate }, nil
	case Object:
		return func(t Triple) Node { return t.Object }, nil
	}
	return nil, fmt.Errorf("invalid node position %d", position)
}

func shortString(n Node) string {
	if a, ok := n.(AnonymousNode); ok {
		return a.String()[0:8]
	} else {
		return n.String()
	}
}

func NewTriple(subject, predicate, object any) (Triple, error) {
	var err error
	var s, p, o Node
	s, err = NewNode(subject)
	if err != nil {
		return Triple{}, err
	}
	p, err = NewNode(predicate)
	if err != nil {
		return Triple{}, err
	}
	o, err = NewNode(object)
	if err != nil {
		return Triple{}, err
	}

	triple := Triple{s, p, o}
	return triple, nil
}

func NewTripleFromNodes(subject, predicate, object Node) Triple {
	return Triple{subject, predicate, object}
}

type TripleSet map[string]Triple

type Triples struct {
	TripleSet
	Nodes NodeSet
}

func NewTriples() *Triples {
	return &Triples{
		TripleSet: make(TripleSet),
		Nodes:     NewNodeSet(),
	}
}

func (source *Triples) NewNode(value interface{}) (Node, error) {
	return NewNode(value)
}

func (source *Triples) NewTripleFromNodes(subject Node, predicate Node, object Node) Triple {
	triple := Triple{subject, predicate, object}
	source.Add(triple)
	return triple
}

func (source *Triples) AddTriple(subject, predicate, object any) (Triple, error) {
	triple, err := NewTriple(subject, predicate, object)
	if err != nil {
		return Triple{}, err
	}
	source.Add(triple)
	return triple, nil
}

func (source *Triples) AddTripleString(subject string, predicate string, object string) Triple {
	triple := Triple{NewStringNode(subject), NewStringNode(predicate), NewStringNode(object)}
	source.Add(triple)
	return triple
}

func (source *Triples) AddTripleReference(triple Triple) Node {
	container := NewAnonymousNode()
	source.NewTripleFromNodes(container, Subject, triple.Subject)
	source.NewTripleFromNodes(container, Predicate, triple.Predicate)
	source.NewTripleFromNodes(container, Object, triple.Object)
	return container
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

func (source *Triples) AddTriples(triples *Triples) {
	for _, triple := range source.TripleSet {
		source.Add(triple)
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
	_, ok := source.TripleSet[triple.String()]
	return ok
}

func (source *Triples) ContainsTriples(triples *Triples) bool {
	for _, triple := range triples.TripleSet {
		if !source.Contains(triple) {
			return false
		}
	}
	return true
}

func (source *Triples) GetTriplesForSubject(subject Node) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if triple.Subject == subject {
			res.Add(triple)
		}
	}
	return res
}

func (source *Triples) GetTriplesForSubjects(subjects NodeSet) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if subjects.Contains(triple.Subject) {
			res.Add(triple)
		}
	}
	return res
}

func (source *Triples) GetTripleListForSubject(node Node) TripleList {
	return source.GetTriplesForSubject(node).GetTripleList()
}

func (source *Triples) GetTriplesForPredicate(predicate Node) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if triple.Predicate == predicate {
			res.Add(triple)
		}
	}
	return res
}

func (source *Triples) GetTriplesForObject(object Node) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if triple.Object == object {
			res.Add(triple)
		}
	}
	return res
}

func (source *Triples) GetSubjects() NodeSet {
	res := NewNodeSet()
	for _, triple := range source.TripleSet {
		res.Add(triple.Subject)
	}
	return res
}

func (source *Triples) GetPredicates() NodeSet {
	res := NewNodeSet()
	for _, triple := range source.TripleSet {
		res.Add(triple.Predicate)
	}
	return res
}

func (source *Triples) GetObjects() NodeSet {
	res := NewNodeSet()
	for _, triple := range source.TripleSet {
		res.Add(triple.Object)
	}
	return res
}

func (source *Triples) GetSubjectList() NodeList {
	return source.GetSubjects().GetSortedNodeList()
}

func (source *Triples) GetPredicateList() NodeList { // TODO: refactor to avoid boilerplate with GetSubjectList
	set := source.GetPredicates()
	keys := make([]string, 0)
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	res := make(NodeList, 0)
	for _, key := range keys {
		res = append(res, set[key])
	}
	return res
}

func (source *Triples) GetTripleList() TripleList {
	tripleList := make(TripleList, 0)
	for _, triple := range source.TripleSet {
		tripleList = append(tripleList, triple)
	}
	return tripleList
}

func (source Triples) String() string {
	res := ""
	for _, triple := range source.GetTripleList().Sort() {
		res += fmt.Sprintf("%s\n", triple)
	}
	return res
}

func (l TripleList) NewTriples() *Triples {
	triples := NewTriples()
	for _, triple := range l {
		triples.Add(triple)
	}
	return triples
}

func (source TripleSet) String() string {
	res := ""
	for _, triple := range source {
		res += fmt.Sprintf("%s\n", triple)
	}
	return res
}

func (source TripleList) String() string {
	res := ""
	for _, triple := range source {
		res += fmt.Sprintf("%s\n", triple)
	}
	return res
}
