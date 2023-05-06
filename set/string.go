package set

type StringShape string

func (s StringShape) Equal(i StringShape) bool {
	return s == i
}

func StringShapes(items ...string) []StringShape {
	res := make([]StringShape, 0)
	for _, item := range items {
		res = append(res, StringShape(item))
	}
	return res
}
