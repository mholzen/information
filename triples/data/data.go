package data

import (
	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"
)

var Data *triples.Triples = nil

func InitData() error {
	Data = triples.NewTriples()
	rowQuery := Data.AddTriplesAsContainer(transforms.RowQuery())
	Data.AddTripleFromNodes(rowQuery, triples.Name, triples.NewStringNode("rowQuery"))
	return nil
}
