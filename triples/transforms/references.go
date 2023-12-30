package transforms

import (
	t "github.com/mholzen/information/triples"
)

func ReferenceTriples(source *t.Triples) *t.Triples {
	return source.Filter(ReferenceTripleMatch)
}

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

func RemoveReferencesMapper(source *t.Triples) (*t.Triples, error) {
	return source.Map(NewIntersectMapper(ReferenceTriples(source)))
}

func ReferenceTriplesConnected(source *t.Triples) *t.Triples {
	referenceTriples := ReferenceTriples(source)

	// TODO: should use traverse
	filter := NewSubjectsTripleMatch(referenceTriples.GetSubjects())
	super := source.Filter(filter)
	referenceTriples.AddTriples(super)

	filter = NewObjectsTripleMatch(referenceTriples.GetSubjects())
	super = source.Filter(filter)
	referenceTriples.AddTriples(super)

	return referenceTriples
}
