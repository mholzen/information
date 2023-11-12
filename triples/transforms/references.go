package transforms

import (
	t "github.com/mholzen/information/triples"
)

func References(source *t.Triples) *t.Triples {
	references := make(map[t.Node]t.Triple)

	for _, triple := range source.TripleSet {
		if triple.Predicate == t.Subject {
			ref := references[triple.Subject]
			ref.Subject = triple.Object
			references[triple.Subject] = ref
		}
		if triple.Predicate == t.Predicate {
			ref := references[triple.Subject]
			ref.Predicate = triple.Object
			references[triple.Subject] = ref
		}
		if triple.Predicate == t.Object {
			ref := references[triple.Subject]
			ref.Object = triple.Object
			references[triple.Subject] = ref
		}
	}

	res := t.NewTriples()
	for _, t := range references {
		if t.Subject != nil && t.Predicate != nil && t.Object != nil {
			res.Add(t)
		}
	}
	return res
}
func ReferencesMapper(source *t.Triples) (*t.Triples, error) {
	return References(source), nil
}
