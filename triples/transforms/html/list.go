package transforms

import (
	"fmt"
	"strings"

	"github.com/mholzen/information/triples"
)

func classes(triple triples.Triple, all *triples.Triples) string {
	if classes := all.GetObjectsBySubjectPredicate(triple.Object, triples.NewStringNode("class")); classes != nil {
		return strings.Join(triples.MakeStringNodes(classes).GetStringList(), " ")
	} else {
		return ""
	}

}
func listrow(triple triples.Triple, all *triples.Triples) string {
	classes := classes(triple, all)
	return fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class=\"%s\">%s</td></tr>\n",
		triple.Subject.String(),
		triple.Predicate.String(),
		classes,
		triple.Object.String(),
	)
}

func NewHtmlListGenerator() *triples.TransformerWithResult {
	var r triples.Node
	return &triples.TransformerWithResult{
		Transformer: func(source *triples.Triples) error {
			res := triples.NewStringNode("")

			res.Value = "<table>\n"
			for _, triple := range source.GetTripleList().Sort() {
				res.Value += listrow(triple, source)
			}
			res.Value += "\n</table>\n"
			r = res
			return nil
		},
		Result: &r,
	}
}
