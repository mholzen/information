package html

import (
	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"
	"github.com/russross/blackfriday/v2"
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
	}

	return Table(data)
}

func Table(source *triples.Triples) (html, error) {
	rows, err := transforms.RowTriples(source)
	if err != nil {
		return "", err
	}
	res := "<table>\n"
	for _, triple := range rows.GetTripleList().Sort() {
		res += tablerow(triple, source)
	}
	res += "\n</table>\n"

	return html(res), nil
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

type markdown string

func (m markdown) Html() string {
	htmlContent := blackfriday.Run([]byte(m))
	return string(htmlContent)
}
