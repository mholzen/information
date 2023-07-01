package triples

type TripleTransform func(source *Triples, triple Triple, i int, root Node) (Triple, error)

func NewMap(start Node, tripleTransform TripleTransform, output *Triples) Transformer {

	return func(source *Triples) error {
		triples := source.GetTriplesForSubject(start)
		triples.Sort()
		for i, triple := range triples {
			triple, err := tripleTransform(source, triple, i, start)
			if err != nil {
				return err
			}
			output.Add(triple)
		}
		return nil
	}
}

func NewTripleObjectTransformer(target Node, dest *Triples) TripleTransform {
	return func(set *Triples, triple Triple, i int, root Node) (Triple, error) {

		if _, ok := triple.Object.(AnonymousNode); ok {
			var objectNode Node
			l := set.GetTriplesForSubject(triple.Object)
			for _, t := range l {
				if t.Predicate == Object {
					objectNode = t.Object
					break
				}
			}
			if objectNode != nil {
				newTriple := set.NewTripleFromNodes(target, NewIndexNode(i), objectNode)
				return newTriple, nil
			}
		}
		return Triple{}, nil
	}
}
