package data

import (
	"github.com/mholzen/information/transforms"
	"github.com/mholzen/information/triples"
)

var Data *triples.Triples = nil

func InitData() error {
	Data = triples.NewTriples()
	rowQuery := Data.AddTriplesAsContainer(transforms.RowQuery().GetTriples())
	Data.AddTriple(rowQuery, triples.Name, triples.NewStringNode("rowQuery"))
	return nil
}
