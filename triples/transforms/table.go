package transforms

import "github.com/mholzen/information/triples"

type TableGenerator struct {
	Transformer triples.Transformer
	Result      [][][]triples.Node
}

type TableDefinition struct {
	Columns       [][]triples.Node
	ColumnFilters []triples.TripleMatch
	SortColumn    int // TODO: should become a function
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
	res := TableGenerator{}
	def := NewTableDefinition(definition)
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
			res.Result = append(res.Result, row)
		}
		return nil
	}
	return &res
}
