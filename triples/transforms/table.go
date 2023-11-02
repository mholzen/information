package transforms

import (
	"strings"

	"github.com/mholzen/information/triples"
	. "github.com/mholzen/information/triples"
)

type Rows [][][]Node

type TableGenerator struct {
	Transformer Transformer
	Definition  TableDefinition
	Rows        Rows
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

func NewTableGenerator(definition *Triples) *TableGenerator {
	def := NewTableDefinition(definition)
	res := TableGenerator{
		Definition: def,
	}

	rowFilter := Filter(def.RowFilter)

	res.Transformer = func(source *Triples) error {
		rows, err := source.Map(rowFilter)
		if err != nil {
			return err
		}
		rowNodes := rows.GetObjects().GetSortedNodeList()

		// for _, subject := range def.GetSubjectList() {
		for _, subject := range rowNodes {
			row := make([][]Node, len(def.Columns))
			for j, filter := range def.ColumnFilters {
				cell := make([]Node, 0)
				for _, triple := range source.GetTripleListForSubject(subject) {
					if filter(triple) {
						cell = append(cell, triple.Object)
					}
				}
				row[j] = cell
			}
			res.Rows = append(res.Rows, row)
		}
		return nil
	}
	return &res
}

type TableDefinition struct {
	Columns       [][]Node
	ColumnFilters []TripleMatch
	RowQuery      *Triples
	RowFilter     TripleMatch
	// SortColumn    int // TODO: should become a function
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

func NewTableDefinition(definition *Triples) TableDefinition {
	columnNo := 0
	res := TableDefinition{
		Columns:       make([][]Node, 0),
		ColumnFilters: make([]TripleMatch, 0),
	}
	for {
		column := NewTriples()
		f := NewPredicateFilter(column, NewIndexNode(columnNo))
		err := definition.Transform(f) // TODO: why can a transform fail?
		if err != nil {
			panic(err)
		}
		if len(column.TripleSet) == 0 {
			break
		}
		col := make([]Node, 0)
		for _, triple := range column.GetTripleList() {
			col = append(col, triple.Object)
		}
		res.Columns = append(res.Columns, col)
		res.ColumnFilters = append(res.ColumnFilters, NewPredicateOrMatch(col...))
		columnNo++
	}

	res.RowQuery = RowQuery()
	filter, err := NewTripleMatchFromTriples(res.RowQuery)
	if err != nil {
		panic(err)
	}
	res.RowFilter = filter

	return res
}

// TODO: make this into a transform
func OrderedPredicateList(rows *Triples) *Triples {
	res := NewTriples()
	container := NewAnonymousNode()

	var columns = make(map[Node]int, 0)

	for _, row := range rows.GetTripleList() { // TODO: Sort
		predicate := row.Predicate
		if _, ok := columns[predicate]; !ok {
			columns[predicate] = len(columns)
			res.AddTriple(container, NewIndexNode(columns[predicate]), predicate)
		}
	}
	return res
}

func RowQuery() *Triples {
	rowQuery := triples.NewTriples()
	rowQuery.AddTriple(triples.NewAnonymousNode(), triples.Predicate, triples.NewNodeMatchAnyIndex1())
	rowQuery.AddTriple(triples.NewAnonymousNode(), triples.Object, triples.NewNodeMatchAnyAnonymous1())
	return rowQuery
}
