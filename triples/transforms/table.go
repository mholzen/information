package transforms

import (
	"strings"

	"github.com/mholzen/information/triples"
)

type Rows [][][]triples.Node

type TableGenerator struct {
	Transformer triples.Transformer
	Definition  TableDefinition
	Rows        Rows
}

type TableDefinition struct {
	Columns       [][]triples.Node
	ColumnFilters []triples.TripleMatch
	// SortColumn    int // TODO: should become a function
}

func NewTableDefinition(definition *triples.Triples) TableDefinition {
	columnNo := 0
	res := TableDefinition{
		Columns:       make([][]triples.Node, 0),
		ColumnFilters: make([]triples.TripleMatch, 0),
	}
	for {
		column := triples.NewTriples()
		f := triples.NewPredicateFilter(column, triples.NewIndexNode(columnNo))
		err := definition.Transform(f) // TODO: why can a transform fail?
		if err != nil {
			panic(err)
		}
		if len(column.TripleSet) == 0 {
			break
		}
		col := make([]triples.Node, 0)
		for _, triple := range column.GetTripleList() {
			col = append(col, triple.Object)
		}
		res.Columns = append(res.Columns, col)
		res.ColumnFilters = append(res.ColumnFilters, triples.NewPredicateOrMatch(col...))
		columnNo++
	}
	return res
}

func NewTableGenerator(definition *triples.Triples) *TableGenerator {
	def := NewTableDefinition(definition)
	res := TableGenerator{
		Definition: def,
	}

	res.Transformer = func(source *triples.Triples) error {
		for _, triple := range source.GetTripleList() {
			row := make([][]triples.Node, len(def.Columns))
			for j, filter := range def.ColumnFilters {
				cell := make([]triples.Node, 0)
				if filter(triple) {
					cell = append(cell, triple.Object)
				}
				row[j] = cell
			}
			res.Rows = append(res.Rows, row)
		}
		return nil
	}
	return &res
}

func (g TableGenerator) Html() string {
	res := make([]string, 0)
	for _, row := range g.Rows {
		cells := make([]string, 0)
		for _, cell := range row {
			nodes := make([]string, 0)
			for _, node := range cell {
				nodes = append(nodes, node.String())
			}
			cells = append(cells, "<td>"+strings.Join(nodes, "<br>")+"</td>")
		}
		res = append(res, "<tr>"+strings.Join(cells, "\n")+"</tr>")
	}
	return "<table>\n" +
		g.Definition.Html() +
		strings.Join(res, "\n") +
		"\n</table>"
}

func (d TableDefinition) Html() string {
	res := make([]string, 0)
	for _, col := range d.Columns {
		nodes := make([]string, 0)
		for _, node := range col {
			nodes = append(nodes, node.String())
		}
		res = append(res, "<th>"+strings.Join(nodes, "<br>")+"</th>")
	}
	return "<tr>\n" + strings.Join(res, "\n") + "\n</tr>"
}

func NewDefaultTableDefinition(source *triples.Triples) *triples.Triples {
	res := triples.NewTriples()
	container := triples.NewAnonymousNode()
	for i, predicate := range source.GetPredicates() {
		res.AddTriple(container, triples.NewIndexNode(i), predicate)
	}
	return res
}
