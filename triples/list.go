package triples

func (l TripleList) GetObjects() []Node {
	var objects []Node
	for _, triple := range l {
		objects = append(objects, triple.Object)
	}
	return objects
}

func (l TripleList) GetObjectStrings() []string {
	var objects []string
	for _, triple := range l {
		objects = append(objects, triple.Object.String())
	}
	return objects
}
