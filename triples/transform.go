package triples

type Transformer func(target *Triples) error

func (source *Triples) Transform(transformer Transformer) error {
	return transformer(source)
}

func NewObjectFilter(target *Triples, objectFn UnaryFunctionNode) Transformer {
	return func(source *Triples) error {
		for _, triple := range source.TripleSet {
			value, err := objectFn(triple.Object)
			if err != nil {
				return err
			}

			if value.(NumberNode).Value != 0 {
				target.Add(triple)
			}
		}
		return nil
	}
}

type TransformerWithResult struct {
	Transformer
	Result *Node
}
