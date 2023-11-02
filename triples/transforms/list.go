package transforms

import (
	. "github.com/mholzen/information/triples"
)

func tablerow(triple Triple) string {
	return "<tr><td>" + triple.Subject.String() + "</td><td>" + triple.Predicate.String() + "</td><td>" + triple.Object.String() + "</td></tr>\n"
}

func NewListGenerator() *TransformerWithResult {
	var r Node
	return &TransformerWithResult{
		Transformer: func(source *Triples) error {
			res := NewStringNode("")

			res.Value = "<table>\n"
			s := source.GetTripleList().Sort()
			for _, triple := range s {
				res.Value += tablerow(triple)
			}
			res.Value += "\n</table>\n"
			r = res
			return nil
		},
		Result: &r,
	}
}
