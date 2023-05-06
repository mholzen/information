package rel

type Relationship interface {
	// A collection, finite or infinite, comparable or not
	// Contains(T) Truth
}

// func (s IdentityList) Get(t Truth) []Identity {
// 	panic("nyi")
// }

// func (s IdentityList) Share(i1 Identity, i2 Identity) Truth {
// 	return Boolean{s.Contains(i1).Boolean() && s.Contains(i2).Boolean()}
// }

// Set
// type MapSet[T comparable] struct {
// 	Relationship
// 	Identities map[T]bool
// }

// func NewMapSet() MapSet[Identity] {
// 	return MapSet{
// 		Identities: make(map[Identity]bool),
// 	}
// }

// func (set *MapSet[T]) Contains(i T) Truth {
// 	_, ok := set.Identities[i]
// 	return Boolean{Bool: ok}
// }

// type Set[T any] map[T]struct{}

// type IdentitySet2 Set[Identity]

// type SubjectVerb struct {
// 	/* This relationship indicates that
// 	- this property should be included in the aggregate that defines the identity
// 	- this identity should be included in the aggregate of idententies that defines this property
// 	*/
// 	Relationship
// 	Identity Identity
// 	Property Property
// }

// SubjectVerb

// func NewSubjectVerb(subject Identity, verb Property, t Truth) SubjectVerb {
// 	subject.Add(verb, t)
// 	verb.Add(subject, t)
// 	return SubjectVerb{Identity: subject, Property: verb}
// }

// type SubjectVerbObject struct {
// 	SubjectVerb SubjectVerb
// 	Object      interface{}
// }

// func (sv SubjectVerbObject) Equal(i Identity) Truth {
// 	sv2 := i.(SubjectVerbObject)
// 	b := sv.SubjectVerb.Equal(sv2.SubjectVerb).Boolean() && sv.Object == sv2.Object
// 	return Boolean{Bool: b}
// }

// func NewSubjectVerbObject(subject Identity, verb Identity, object interface{}) SubjectVerbObject {
// 	return SubjectVerbObject{
// 		SubjectVerb: NewSubjectVerb(subject, verb),
// 		Object:      object,
// 	}
// }
