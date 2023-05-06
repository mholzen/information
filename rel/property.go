package rel

// type Property interface {

// 	// A relationship between identities,
// 	Aggregate[Identity]

// 	// indicating that they share the same property
// 	// a property probably needs at least two identities
// 	Share(Identity, Identity) Truth

// 	// with an expectation that the set is infinite.
// 	/*
// 		Eg: these two apples are red:
// 			apple1, apple2 Identity
// 			color Property
// 			color.Share(apple1, apple2)
// 	*/
// }

// func NewProperty() Property {
// 	p := make(IdentityList, 0)
// 	return &p
// }

// type PropertyPair[T Property] struct {
// 	Property
// 	left, right T
// }

// func (p PropertyPair[T]) Share(i1 Identity, i2 Identity) Truth {
// 	leftT := p.left.Share(i1, i2)
// 	rightT := p.right.Share(i1, i2)

// 	return Boolean{Bool: leftT.Boolean() && rightT.Boolean()}
// }

// type PropertyIdentitySet struct { // This implementation specifically does support infinite size
// 	Identities MapSet[Identity]
// }

// func (s PropertyIdentitySet) Contains(i Identity) Truth {
// 	return s.Identities.Contains(i).(Boolean)
// }
// func (s PropertyIdentitySet) Add(i Identity, t Truth) {
// }

// func (s PropertyIdentitySet) Share(i1 Identity, i2 Identity) Truth {
// 	res := And{
// 		Property1: s.Identities.Contains(i1).(Boolean),
// 		Property2: s.Identities.Contains(i2).(Boolean),
// 	}
// 	return res
// }
// func NewPropertyIdentitySet() PropertyIdentitySet {
// 	return PropertyIdentitySet{Identities: NewMapSet()}
// }

// type PropertySet map[Property]struct{} // TODO: why can't I use MapSet[Identity]?

// func (s PropertySet) Contains(p Property) Truth {
// 	_, ok := s[p]
// 	return Boolean{ok}
// }

// func (s PropertySet) Add(p Property, t Truth) {
// 	if t.Boolean() {
// 		s[p] = struct{}{}
// 	}
// }
// func (s PropertySet) Share(p1 Property, p2 Property) Truth {
// 	return Boolean{s.Contains(p1).Boolean() && s.Contains(p2).Boolean()}
// }

// // Property List
// type PropertyList []struct {
// 	Property Property
// 	Truth    Truth
// }

// func (s PropertyList) Contains(p Property) Truth {

// }

// func (s PropertyList) Add(p Property, t Truth) {
// 	panic("nyi")
// }

// func (s PropertyList) Share(p1 Property, p2 Property) Truth {
// 	return Boolean{s.Contains(p1).Boolean() && s.Contains(p2).Boolean()}
// }

// type PropertyEven struct {
// }

// type Truth interface {
// 	Property
// 	Boolean() bool
// 	// Probability() float32
// }

// type Boolean struct {
// 	Bool bool
// }

// func (b Boolean) Contains(i1 Identity) Truth {
// 	return Boolean{Bool: false}
// }
// func (b Boolean) Add(i1 Identity, t Truth) {
// }
// func (b Boolean) Get(t Truth) []Identity {
// 	return nil
// }

// func (b Boolean) Share(i1 Identity, i2 Identity) Truth {
// 	return Boolean{Bool: false}
// }
// func (b Boolean) Boolean() bool {
// 	return b.Bool
// }

// type And struct {
// 	Truth
// 	Property1 Boolean
// 	Property2 Boolean
// }

// func (b And) Boolean() bool {
// 	return b.Property1.Bool && b.Property2.Bool
// }

// func NewTextProperty(text string) Property {
// 	return Text(text)
// }

// type Text string

// func (t Text) Contains(i Identity) Truth {
// 	return Boolean{Bool: false}
// }
// func (text Text) Add(i Identity, t Truth) {
// }
// func (text Text) Get(t Truth) []Identity {
// 	return nil
// }

// func (t Text) Share(i1 Identity, i2 Identity) Truth {
// 	return Boolean{Bool: false}
// }
