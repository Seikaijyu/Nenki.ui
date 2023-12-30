package layout

import "nenki.ui/monad"

var _ monad.MonadInterface[*Layout, int] = &Layout{}

type Layout struct {
	value *int
}

func (p *Layout) Bind(fn func(*Layout) *Layout) *Layout {
	return fn(p)
}
func (p *Layout) Unit(value int) *Layout {
	return NewLayout(value)
}

func NewLayout(v int) *Layout {
	return &Layout{value: &v}
}
