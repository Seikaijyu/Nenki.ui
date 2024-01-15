package theme

import (
	gtext "github.com/Seikaijyu/gio/text"
	"github.com/Seikaijyu/gio/widget/material"
)

var Shaper *gtext.Shaper = &gtext.Shaper{}

func NewTheme() *material.Theme {
	theme := material.NewTheme()
	theme.Shaper = Shaper
	return theme
}
