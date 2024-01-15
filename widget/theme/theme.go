package theme

import (
	gtext "gioui.org/text"
	"gioui.org/widget/material"
)

var Shaper *gtext.Shaper = &gtext.Shaper{}

func NewTheme() *material.Theme {
	theme := material.NewTheme()
	theme.Shaper = Shaper
	return theme
}
