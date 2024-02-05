package triples

type Mapper func(source *Triples) (*Triples, error)

type SafeMapper func(source *Triples) *Triples

func (source *Triples) Map(mapper Mapper) (*Triples, error) {
	return mapper(source)
}

func (source *Triples) FlatMap(generator TriplesGenerator) (*Triples, error) {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		triples, err := generator(triple)
		if err != nil {
			return nil, err
		}
		if triples != nil {
			res.AddTriples(triples)
		}
	}
	return res, nil
}

func (source *Triples) MapTriples(mapper TripleMapper) (*Triples, error) {
	res := NewTriples()
	for _, triple := range source.TripleSet {
		newTriple, err := mapper(triple)
		if err != nil {
			return nil, err
		}
		res.Add(newTriple)
	}
	return res, nil
}

func (source *Triples) MapToNode(mapper MapperToNode) (Node, error) {
	return mapper(source)
}

type MapperToNode func(source *Triples) (Node, error)

type TripleTransform func(source *Triples, triple Triple, i int, root Node) (Triple, error)

type TripleMapper func(triple Triple) (Triple, error)

type TriplesGenerator func(triple Triple) (*Triples, error)

func MapperFromTransfomer(transformer Transformer) Mapper {
	return func(source *Triples) (*Triples, error) {
		res := source.Clone()
		err := res.Transform(transformer)
		return res, err
	}
}

func TransformerFromMapper(mapper Mapper) Transformer {
	return func(source *Triples) error {
		res, err := source.Map(mapper)
		if err != nil {
			return err
		}
		*source = *res
		return nil
	}
}
