package paths

type Path struct {
	Val *string
}

func (p *Path) Set(val string) {
	p.Val = &val
}

func (p *Path) String() string {
	if p == nil || p.Val == nil {
		return "<empty>"
	}
	return *p.Val
}

func Empty() *Path {
	return &Path{}
}
