package ctxfmt

type argstate struct {
	idx  int
	args []interface{}
}

func (a *argstate) next() (arg interface{}, idx int, has bool) {
	if a.idx < len(a.args) {
		arg, idx = a.args[a.idx], a.idx
		a.idx++
		return arg, idx, true
	}
	return nil, len(a.args), false
}
