package triples

type Transformer func(target *Triples) error

func (source *Triples) Transform(transformer Transformer) error {
	return transformer(source)
}
