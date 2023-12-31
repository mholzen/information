package transforms

import (
	"regexp"

	"github.com/mholzen/information/triples"
)

func IdLinesTransformer() triples.Transformer {
	re := regexp.MustCompile(`^\d+/\d+/\d+ \d+:\d+:`)
	tr := NewObjectAugmenter(
		NewObjectRegexpTripleMatch(re),
		triples.NewStringNode("class"),
		triples.NewStringNode("line"),
	)
	return tr
}

var IdLinesMapper = triples.MapperFromTransfomer(IdLinesTransformer())
