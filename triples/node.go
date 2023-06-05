package triples

import (
	"time"
)

type CreatedNode[T any] struct {
	Value   T
	Created time.Time
}
