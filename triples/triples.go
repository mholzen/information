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
	return fmt.Sprintf("(%s, %s, %s)", ShortString(t.Subject), ShortString(t.Predicate), ShortString(t.Object))
}

type NodePosition byte

const (
	SubjectPosition NodePosition = iota
	PredicatePosition
	ObjectPosition
)

func (p NodePosition) Getter() (func(Triple) Node, error) {
	switch p {
	case SubjectPosition:
		return func(t Triple) Node {
			return t.Subject
		}, nil
	case PredicatePosition:
		return func(t Triple) Node {
			return t.Predicate
		}, nil
	case ObjectPosition:
		return func(t Triple) Node {
			return t.Object
		}, nil
	}
	return nil, fmt.Errorf("invalid node position %d", p)
}

func (p NodePosition) String() string {
	return [...]string{"subject", "predicate", "object"}[p]
}

func (p NodePosition) WrapError(err error) error {
	return fmt.Errorf("error setting %s: %s", p.String(), err)
}

func GetNodeFunction(position Node) (func(Triple) Node, error) {
	switch position {
	case Subject:
		return func(t Triple) Node {
			return t.Subject
		}, nil
	case Predicate:
		return func(t Triple) Node {
			return t.Predicate
		}, nil
	case Object:
		return func(t Triple) Node {
			return t.Object
		}, nil
	}
	return nil, fmt.Errorf("invalid node position %d", position)
}

func ShortString(n Node) string {
	// TODO: handle VariableNode
	if a, ok := n.(AnonymousNode); ok {
		return a.String()[0:8]
	} else if n == nil {
		return "<nil>"
	} else {
		return n.String()
	}
}

func NewTriple(subject, predicate, object Node) Triple {
	return Triple{subject, predicate, object}
}

func NewTripleFromAny(subject, predicate, object any) (Triple, error) {
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

func NewTriplesFromList(triples TripleList) *Triples {
	res := NewTriples()
	for _, triple := range triples {
		res.Add(triple)
	}
	return res
}

func NewTriplesFromAny(node ...any) (*Triples, error) {
	res := NewTriples()
	err := res.AddTripleNodes(node...)
	return res, err
}

func (source *Triples) Length() int {
	return len(source.TripleSet)
}

func (source *Triples) IsEmpty() bool {
	return source.Length() == 0
}

func (source *Triples) NewNode(value interface{}) (Node, error) {
	return NewNode(value)
}

func (source *Triples) Add(triple Triple) *Triples {
	if _, ok := source.TripleSet[triple.String()]; ok {
		return source
	}
	source.TripleSet[triple.String()] = triple

	source.Nodes.Add(triple.Subject)
	source.Nodes.Add(triple.Predicate)
	if triple.Object != nil {
		source.Nodes.Add(triple.Object)
	}
	return source
}

func (source *Triples) AddTriple(subject Node, predicate Node, object Node) Triple {
	triple := Triple{subject, predicate, object}
	source.Add(triple)
	return triple
}

func (source *Triples) AddTripleFromAny(subject, predicate, object any) (Triple, error) {
	triple, err := NewTripleFromAny(subject, predicate, object)
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
	source.AddTriple(container, Subject, triple.Subject)
	source.AddTriple(container, Predicate, triple.Predicate)
	source.AddTriple(container, Object, triple.Object)
	return container
}

func (source *Triples) AddTripleReferences(triples *Triples) Node {
	container := NewAnonymousNode()
	for _, triple := range triples.TripleSet {
		node := source.AddTripleReference(triple)
		source.Add(NewTriple(container, Contains, node))
	}
	return container
}

func (source *Triples) AddTripleList(triple ...Triple) *Triples {
	for _, triple := range triple {
		source.Add(triple)
	}
	return source
}

func (source *Triples) Delete(triple Triple) {
	delete(source.TripleSet, triple.String())
	// should we delete nodes?
}

func (source *Triples) DeleteTriples(triples *Triples) {
	for _, triple := range triples.TripleSet {
		source.Delete(triple)
	}
}

func (source *Triples) AddTriples(triples *Triples) {
	for _, triple := range triples.TripleSet {
		source.Add(triple)
	}
}

func (source *Triples) AddTripleNodes(nodes ...any) error {
	if len(nodes)%3 != 0 {
		return fmt.Errorf("number of nodes (%d) not a multiple of 3", len(nodes))
	}
	for len(nodes) > 0 {
		triple, err := NewTripleFromAny(nodes[0], nodes[1], nodes[2])
		if err != nil {
			return err
		}
		source.Add(triple)
		nodes = nodes[3:]
	}
	return nil
}

func (source *Triples) AddTriplesAsContainer(triples *Triples) Node {
	container := NewAnonymousNode()
	for _, triple := range triples.TripleSet {
		source.Add(triple)
		source.AddTriple(container, Contains, triple.Subject)
	}
	return container
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

			source.AddTriple(subject, predicate, object)
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

			source.AddTriple(container, predicate, object)
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

func (source *Triples) Clone() *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		res.Add(triple)
	}
	return res
}

func (source *Triples) Contains(triple Triple) bool {
	_, ok := source.TripleSet[triple.String()]
	return ok
}

func (source *Triples) ContainsTriple(subject, predicate, object any) bool {
	triple, err := NewTripleFromAny(subject, predicate, object)
	if err != nil {
		return false
	}
	return source.Contains(triple)
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

func (source *Triples) GetTriplesForSubjectPredicate(subject, predicate Node) *Triples {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		if triple.Subject == subject && triple.Predicate == predicate {
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

func (source *Triples) GetTripleListForObject(node Node) TripleList {
	return source.GetTriplesForObject(node).GetTripleList()
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

func (source *Triples) GetTriplesForPredicateString(predicate string) *Triples {
	return source.GetTriplesForPredicate(NewStringNode(predicate))
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

func (source *Triples) GetSubjectsByPredicateObject(predicate, object Node) NodeSet {
	res := make(NodeSet)
	for _, triple := range source.TripleSet {
		if NodeEquals(triple.Predicate, predicate) && NodeEquals(triple.Object, object) {
			res.Add(triple.Subject)
		}
	}
	return res
}

func (source *Triples) GetObjectsBySubjectPredicate(subject, predicate Node) NodeSet {
	res := make(NodeSet)
	for _, triple := range source.TripleSet {
		if NodeEquals(triple.Subject, subject) && NodeEquals(triple.Predicate, predicate) {
			res.Add(triple.Object)
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
	return source.GetSubjects().GetNodeList()
}

func (source *Triples) GetPredicateList() NodeList {
	return source.GetPredicates().GetNodeList()
}

func (source *Triples) GetObjectList() NodeList {
	return source.GetObjects().GetNodeList()
}

func (source *Triples) GetTripleList() TripleList {
	tripleList := make(TripleList, 0)
	for _, triple := range source.TripleSet {
		tripleList = append(tripleList, triple)
	}
	return tripleList
}

func (source *Triples) GetTripleReferences(triple Triple) NodeSet {
	subjects := source.GetSubjectsByPredicateObject(Subject, triple.Subject)
	if len(subjects) == 0 {
		return nil
	}
	predicates := source.GetSubjectsByPredicateObject(Predicate, triple.Predicate)
	if len(predicates) == 0 {
		return nil
	}
	objects := source.GetSubjectsByPredicateObject(Object, triple.Object)
	if len(objects) == 0 {
		return nil
	}
	return subjects.Intersect(predicates).Intersect(objects)
}

func (source Triples) String() string {
	res := ""
	for _, triple := range source.GetTripleList().Sort() {
		res += fmt.Sprintf("%s\n", triple)
	}
	return res
}

func (source Triples) StringLine() string {
	res := ""
	for _, triple := range source.GetTripleList().Sort() {
		res += fmt.Sprintf("%s ", triple)
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
