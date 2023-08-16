package triples

type Mapper func(source *Triples) (*Triples, error)

func (source *Triples) Map(mapper Mapper) (*Triples, error) {
	return mapper(source)
}

type TripleTransform func(source *Triples, triple Triple, i int, root Node) (Triple, error)

type TripleMapper func(triple Triple) (*Triples, error)
