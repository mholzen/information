package transforms

import (
	"github.com/mholzen/information/triples"
)

func listrow(triple triples.Triple) string {
	return "<tr><td>" + triple.Subject.String() + "</td><td>" + triple.Predicate.String() + "</td><td>" + triple.Object.String() + "</td></tr>\n"
}

func NewListGenerator() *triples.TransformerWithResult {
	var r triples.Node
	return &triples.TransformerWithResult{
		Transformer: func(source *triples.Triples) error {
			res := triples.NewStringNode("")

			res.Value = "<table>\n"
			for _, triple := range source.GetTripleList().Sort() {
				res.Value += listrow(triple)
			}
			res.Value += "\n</table>\n"
			r = res
			return nil
		},
		Result: &r,
	}
}
