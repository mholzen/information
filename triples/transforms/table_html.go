package transforms

import (
	"github.com/mholzen/information/triples"
)

func HtmlTable(source *triples.Triples) (*triples.Triples, error) {

	rows, err := RowTriples(source)
	if err != nil {
		return nil, err
	}
	html := "<table>\n"
	for _, triple := range rows.GetTripleList().Sort() {
		html += tablerow(triple, source)
	}
	html += "\n</table>\n"

	res := triples.NewTriples()
	table := triples.NewAnonymousNode()
	res.AddTriple(table, "html", html)

	return res, nil
}

func tablerow(triple triples.Triple, source *triples.Triples) string {
	rows := source.GetTriplesForSubject(triple.Object).GetTripleList().Sort()
	if len(rows) == 0 {
		return ""
	}
	res := "<tr>\n"
	for _, row := range rows {
		res += tablecell(row.Object, source)
	}
	res += "</tr>\n"
	return res
}

func tablecell(node triples.Node, source *triples.Triples) string {
	res := "<td>\n"
	for _, cell := range source.GetTriplesForSubject(node).GetTripleList().Sort() {
		res += cell.Object.String() + "<br>\n"
	}
	res += "</td>\n"
	return res
}
