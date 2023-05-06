package rel

// type Identity interface {
// 	// An aggregate of properties
// 	// indicating that they all define this identity
// 	Aggregate[Property]

// 	// that helps us authenticate signals as emanating from the identity we're familiar with
// 	// IsIdentifiedBy(...Property) Truth
// 	Equal(Identity) Truth
// }

// type Uuid struct {
// 	String     string
// 	Properties PropertyList
// }

// func NewUuid() Uuid {
// 	return Uuid{
// 		String:     uuid.NewString(),
// 		Properties: make(PropertyList, 0),
// 	}
// }
// func (u Uuid) Equal(i Identity) Truth {
// 	v := i.(Uuid)
// 	return Boolean{Bool: v.String == u.String}
// }
// func (u Uuid) Add(p Property, t Truth) {
// 	u.Properties.Add(p, t)
// }
// func (u Uuid) Contains(p Property) Truth {
// 	if text, ok := p.(Text); ok {
// 		if text == Text(u.String) {
// 			return Boolean{true}
// 		}
// 	}
// 	return Boolean{false}
// }
// func (u Uuid) Get(t Truth) []Property {
// 	panic("nyi")
// }

// func (u Uuid) IsIdentifiedBy(ps ...Property) Truth {
// 	for _, p := range ps {
// 		u.Add(p)
// 	}
// 	return Boolean{false}
// }

// func NewIdentity() Identity {
// 	return NewUuid()
// }

// func NewNamedIdentity(name string) Identity {
// 	res := NewUuid()
// 	// x = NewSubjectVerbObject(res, name, "Marc")
// 	return res
// }

// Pair

// type IdentityPair[T Identity] struct {
// 	Relationship
// 	Identity
// 	left, right T
// }

// func (i IdentityPair[T]) Contains(p Property) Truth {
// 	// b := p.left.Equal(i).Boolean() || p.right.Equal(i).Boolean()
// 	return Boolean{Bool: false}
// }

// func (i IdentityPair[T]) Equal(i1 Identity) Truth {
// 	p2, ok := i1.(IdentityPair[T])
// 	if !ok {
// 		return Boolean{Bool: false}
// 	}

// 	b := i.left.Equal(p2.left).Boolean() && i.right.Equal(p2.right).Boolean()
// 	return Boolean{Bool: b}
// }

// func (i IdentityPair[T]) GetMembers() (Identity, Identity) {
// 	return i.left, i.right
// }

// type IdentityPairs[T Identity] []IdentityPair[T] // TODO: should be an Aggregate

// func (s *IdentityPairs[T]) add(x IdentityPair[T]) {
// 	*s = append(*s, x)
// }

// type IdentitySet map[uint32]struct{} // TODO: why can't I use MapSet[T comparable]?

// func (s IdentitySet) Contains(i Identity) Truth {
// 	_, ok := s[i]
// 	return Boolean{ok}
// }

// func (s IdentitySet) Add(i Identity, t Truth) {
// 	properties := i.Get(Boolean{true})
// 	s[properties] = struct{}{}
// }

// func (s IdentitySet) Get(t Truth) []Identity {
// 	return nil
// }

// func (s IdentitySet) Share(i1 Identity, i2 Identity) Truth {
// 	return Boolean{s.Contains(i1).Boolean() && s.Contains(i2).Boolean()}
// }

// type IdentityList []Identity

// func (s IdentityList) Contains(identity Identity) Truth {
// 	for _, i := range s {
// 		if i.Equal(identity).Boolean() {
// 			return Boolean{true}
// 		}
// 	}
// 	return Boolean{false}
// }

// func (s *IdentityList) Add(i Identity, t Truth) {
// 	if s.Contains(i).Boolean() {
// 		return
// 	}
// 	*s = append(*s, i)
// }

// func (s IdentityList) Get(t Truth) []Identity {
// 	panic("nyi")
// }

// func (s IdentityList) Share(i1 Identity, i2 Identity) Truth {
// 	return Boolean{s.Contains(i1).Boolean() && s.Contains(i2).Boolean()}
// }
