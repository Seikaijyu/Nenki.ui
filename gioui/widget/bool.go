// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
)

type Bool struct {
	Value    bool
	_checked func(bool)
	clk      Clickable
}

// Update the widget state and report whether Value was changed.
func (b *Bool) Update(gtx layout.Context) bool {
	changed := false
	for b.clk.Clicked(gtx) {
		b.Value = !b.Value
		changed = true
		if b._checked != nil {
			b._checked(b.Value)
		}
	}
	return changed
}

// Hovered reports whether pointer is over the element.
func (b *Bool) Hovered() bool {
	return b.clk.Hovered()
}

// Pressed reports whether pointer is pressing the element.
func (b *Bool) OnChecked(fn func(bool)) {
	b._checked = fn
}

// Focused reports whether b has focus.
func (b *Bool) Focused() bool {
	return b.clk.Focused()
}

func (b *Bool) History() []Press {
	return b.clk.History()
}

func (b *Bool) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	b.Update(gtx)
	dims := b.clk.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.SelectedOp(b.Value).Add(gtx.Ops)
		semantic.EnabledOp(gtx.Queue != nil).Add(gtx.Ops)
		return w(gtx)
	})
	return dims
}
