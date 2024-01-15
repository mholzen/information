package transforms

import (
	t "github.com/mholzen/information/triples"
)

func NewMap(start t.Node, tripleTransform t.TripleTransform, output *t.Triples) t.Transformer {

	return func(source *t.Triples) error {
		triples := source.GetTripleListForSubject(start)
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

func NewFlatMap(start t.Node, mapper t.TriplesGenerator, output *t.Triples) t.Transformer {

	return func(source *t.Triples) error {

		triples := source.GetTripleListForSubject(start)
		triples.Sort()
		for _, triple := range triples {
			triples, err := mapper(triple)
			if err != nil {
				return err
			}
			output.AddTriples(triples)
		}
		return nil
	}
}

func GetStringObjectMapper(triple t.Triple) (*t.Triples, error) {
	triples := t.NewTriples()
	if _, ok := triple.Object.(t.StringNode); ok {
		triples.Add(triple)
	}
	return triples, nil
}
