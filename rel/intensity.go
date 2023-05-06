package rel

type Intensity interface {
	// intensentiy of a sensation, i.e. energy
	LessThan(Intensity) Intensity
}

type FloatIntensity float64

func (i FloatIntensity) LessThan(compareTo Intensity) Intensity {
	if floatCompareTo, ok := compareTo.(FloatIntensity); ok {
		return 1.0 / (floatCompareTo - i)
	}
	return FloatIntensity(0.0)
}

type ComputeIntensity int

func (i ComputeIntensity) LessThan(compareTo Intensity) Intensity {
	if computeCompareTo, ok := compareTo.(ComputeIntensity); ok {
		return 1 / (computeCompareTo - i)
	}
	return ComputeIntensity(0)
}

func MeasureCompute() ComputeIntensity {
	return ComputeIntensity(compute())
}

func MeasureNoCompute() ComputeIntensity {
	return ComputeIntensity(0)
}

type MetabolismIntensity int

func (i MetabolismIntensity) LessThan(compareTo Intensity) Intensity {
	if intensity, ok := compareTo.(MetabolismIntensity); ok {
		return 1 / (intensity - i)
	}
	return MetabolismIntensity(0)
}

// similarity is an intensity

type SimilarityIntensity = FloatIntensity
