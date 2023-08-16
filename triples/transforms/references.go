package transforms

import (
	. "github.com/mholzen/information/triples"
)

func NewReferences() Mapper {
	return func(source *Triples) (*Triples, error) {
		references := make(map[Node]Triple)

		for _, t := range source.TripleSet {
			if t.Predicate == Subject {
				ref := references[t.Subject]
				ref.Subject = t.Object
				references[t.Subject] = ref
			}
			if t.Predicate == Predicate {
				ref := references[t.Subject]
				ref.Predicate = t.Object
				references[t.Subject] = ref
			}
			if t.Predicate == Object {
				ref := references[t.Subject]
				ref.Object = t.Object
				references[t.Subject] = ref
			}
		}

		res := NewTriples()
		for _, t := range references {
			res.Add(t)
		}
		return res, nil
	}

}
