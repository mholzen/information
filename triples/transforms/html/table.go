package html

import (
	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"
)

type html string

func FromTriples(data *triples.Triples) (html, error) {
	tpls := data.GetTriplesForPredicate(triples.NewStringNode("html")).GetTripleList()
	if len(tpls) >= 0 {
		res := ""
		for _, tpl := range tpls {
			res += tpl.Object.String()
		}
		return html(res), nil
	} else {
		return html(table(data)), nil
	}
}

func HtmlTableMapper(source *triples.Triples) (*triples.Triples, error) {
	res := triples.NewTriples()
	container := triples.NewAnonymousNode()
	htmlTable := table((source))
	res.AddTriple(container, "html", triples.NewStringNode(htmlTable))
	return res, nil
}

func table(source *triples.Triples) string {
	rows, err := transforms.RowTriples(source)
	if err != nil {
		return ""
	}
	html := "<table>\n"
	for _, triple := range rows.GetTripleList().Sort() {
		html += tablerow(triple, source)
	}
	html += "\n</table>\n"
	return html
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
