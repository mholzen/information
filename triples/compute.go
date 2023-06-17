package triples

type TriplesFunction func(*Triples) error

var Compute TriplesFunction = func(source *Triples) error {
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

// TODO: can we use Compute as a predicate against a node?  can it be expressed as a reducer?
