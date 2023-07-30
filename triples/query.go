package triples

func TestOrExclude(variable Node, value Node) {
	if variable, ok := variable.(VariableNode); ok {
		variable.Exclude(value)
	}
}

func Equal(a Node, b Node) bool {
	if a, ok := a.(VariableNode); ok {
		return !a.Excludes.Contains(b)
	}
	if b, ok := b.(VariableNode); ok {
		return !b.Excludes.Contains(a)
	}
	return a == b
}

func Match(triple Triple, query Triple) bool {
	subjectMatch := Equal(triple.Subject, query.Subject)
	predicateMatch := Equal(triple.Predicate, query.Predicate)
	objectMatch := Equal(triple.Object, query.Object)

	tripleMatch := subjectMatch && predicateMatch && objectMatch
	if !tripleMatch {
		TestOrExclude(query.Subject, triple.Subject)
		TestOrExclude(query.Predicate, triple.Predicate)
		TestOrExclude(query.Object, triple.Object)
	}
	return tripleMatch
}

func NewQueryMatch(dest *Triples, query *Triples) Transformer {
	return func(source *Triples) error {
		for _, q := range query.TripleSet {
			for _, t := range source.TripleSet {
				if Match(t, q) {
					dest.Add(t)
				}
			}
		}
		return nil
	}
}
