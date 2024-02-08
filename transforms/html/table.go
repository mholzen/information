package html

import (
	"github.com/mholzen/information/transforms"
	"github.com/mholzen/information/triples"
)

type html string

func FromTriples(data *triples.Triples) (html, error) {
	tpls := data.GetTriplesForPredicate(triples.NewStringNode("html")).GetTripleList()
	if len(tpls) > 0 {
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
	htmlTable := table(source)
	res.AddTriple(container, triples.Str("html"), triples.NewStringNode(htmlTable))
	return res, nil
}

func table(source *triples.Triples) string {
	rows, err := transforms.RowTriples(source)
	if err != nil {
		return ""
	}
	// exclude rows triples from the data
	// to avoid the predice list containing row indices
	data := source.Clone()
	data.DeleteTriples(rows)

	headers := data.GetPredicateList().SortLexical()

	html := "<table>\n"
	html += headerrow(headers)
	for _, triple := range rows.GetTripleList().Sort() {
		html += tablerow(triple, headers, data)
	}
	html += "\n</table>\n"
	return html
}

func headerrow(headers triples.NodeList) string {
	res := "<tr>\n"
	for _, node := range headers {
		res += "<th>" + node.String() + "</th>\n"
	}
	res += "</tr>\n"
	return res
}

func tablerow(triple triples.Triple, headers triples.NodeList, source *triples.Triples) string {
	res := "<tr>\n"

	for _, header := range headers {
		cellTriples := source.GetTriplesForSubjectPredicate(triple.Object, header)
		res += tablecell(cellTriples, source)
	}

	res += "</tr>\n"
	return res
}

func tablecell(cellTriples, source *triples.Triples) string {
	res := "<td>"
	for _, cell := range cellTriples.GetTripleList().Sort() {
		res += cell.Object.String() + "<br>\n"
	}
	res += "</td>\n"
	return res
}
