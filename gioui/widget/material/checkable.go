// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/internal/f32color"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type checkable struct {
	Label              string
	Color              color.NRGBA
	Font               font.Font
	TextSize           unit.Sp
	IconColor          color.NRGBA
	Size               unit.Dp
	shaper             *text.Shaper
	checkedStateIcon   *widget.Icon
	uncheckedStateIcon *widget.Icon
}

func (c *checkable) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
	var icon *widget.Icon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {

			return layout.Stack{Alignment: layout.N}.Layout(gtx,

				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					defer op.Offset(image.Pt(0, 2)).Push(gtx.Ops).Pop()
					size := gtx.Dp(c.Size)
					col := c.IconColor
					if gtx.Queue == nil {
						col = f32color.Disabled(col)
					}
					gtx.Constraints.Min = image.Point{X: size}
					icon.Layout(gtx, col)
					return layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			defer op.Offset(image.Pt(0, 0)).Push(gtx.Ops).Pop()
			colMacro := op.Record(gtx.Ops)
			paint.ColorOp{Color: c.Color}.Add(gtx.Ops)

			return widget.Label{}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label, colMacro.Stop())
		}),
	)
	return dims
}
