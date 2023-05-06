package rel

// biological organism learn positive/negative from food/threat through natural selection
// func Test_more_compute_is_better(t *testing.T) {
// 	compute1 := ComputeIntensity(compute())
// 	compute2 := ComputeIntensity(compute() + compute())

// 	assert.True(t, compute1.LessThan(compute2))
// }

// func Test_0_and_1_look_different(t *testing.T) {
// 	one := SymbolShape(1)
// 	zero := SymbolShape(0)

// 	// TODO: need something high intensity to compare against

// 	// then learns differentiation/similarity

// 	// computational organism should learn to differentiate compute resources
// 	// perhaps spent 1 ms counting

// 	compute1 := MeasureCompute()
// 	compute2 := MeasureCompute()

// 	assert.True(t, one.Equal(zero).LessThan(compute1.Equal(compute2)))

// 	// how do we get a boolean from a LessThan that returns an intensity?
// }

// func Test_identity(t *testing.T) {
// 	x1 := NewIdentity()
// 	x2 := NewIdentity()
// 	assert.NotEqual(t, x1, x2)

// 	assert.False(t, x1.Equal(x2).Boolean())

// 	p1 := NewProperty()
// 	p1.Add(x1, Boolean{true})

// 	assert.True(t, p1.Contains(x1).Boolean())
// }

// func Test_define_identity(t *testing.T) {
// 	marc := NewIdentity()
// 	alive := NewProperty()
// 	marc.Add(alive, Boolean{true})

// 	assert.True(t, marc.Contains(alive).Boolean())
// }

// func Test_text_property(t *testing.T) {
// 	marc := NewIdentity()
// 	namedMarc := NewTextProperty("Marc")
// 	marc.Add(namedMarc, Boolean{true})

// 	assert.True(t, marc.Contains(namedMarc).Boolean())
// }

// func _Test_property(t *testing.T) {

// 	p1 := NewProperty()
// 	p2 := NewProperty()

// 	assert.NotEqual(t, p1, p2, "the equality of two empty properties should be unknown")
// }
