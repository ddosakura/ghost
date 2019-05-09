package sh

// Val for shell
type Val struct {
	v interface{}
}

var (
	val0 = val(0)
)

func val(v interface{}) *Val {
	// TODO:
	return &Val{
		v: v,
	}
}

func (v *Val) String() string {
	// TODO:
	return "Test"
}

// Int value
func (v *Val) Int() int {
	// TODO:
	return 0
}
