package rel

type Sensation interface {
	// stimuli of a particular sensory organ
	Equal(Sensation) SimilarityIntensity
}

// Position, Shape, Color
// Symbol, Boolean, Text

type SymbolShape byte

func (s SymbolShape) Equal(to Sensation) SimilarityIntensity {
	if toSymboleShape, ok := to.(SymbolShape); ok {
		if toSymboleShape == s {
			return SimilarityIntensity(1.0)
		}
	}
	return SimilarityIntensity(0)
}

type StringShape string

func (s StringShape) Equal(to Sensation) SimilarityIntensity {
	if sensation, ok := to.(StringShape); ok {
		if sensation == s {
			return SimilarityIntensity(1.0)
		}
	}
	return SimilarityIntensity(0)
}
