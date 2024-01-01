// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"os"

	appx "gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	gwidget "gioui.org/widget"
	"gioui.org/widget/material"
	"nenki.ui/app"
	"nenki.ui/utils"
	"nenki.ui/widget"
)

func aa(a *app.App, i int, btn *widget.Button) {
	a.Then(func(this *app.App, al *widget.AnchorLayout) {
		btn.SetText(fmt.Sprintf("按钮%d", i))
	})

}
func main() {
	//Test()
	app.NewApp("测试").SetMinSize(1024, 1024).SetSize(1024, 1024).SetTitle("你好").
		Then(func(self *app.App, root *widget.AnchorLayout) {
			root.SetDirection(widget.Center)
			self.SetBackground(utils.HexToRGBA("#efaf00"))
			btn := widget.NewButton("你好").SetFontSize(50).SetBackground(utils.HexToRGBA("#ef00fa")).SetCornerRadius(30)
			root.AppendChild(btn)
			for i := 0; i < 256; i++ {

				aa(self, i, btn)
			}
		})
	// 阻塞
	app.Run()
}

// 点击事件
func click(b *widget.Button) {
	fmt.Println("点击了按钮")
}
func Test() {
	go func() {
		w := appx.NewWindow(appx.Decorated(false))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	appx.Main()
}

func loop(w *appx.Window) error {
	var (
		b    gwidget.Clickable
		deco gwidget.Decorations
	)
	var (
		toggle    bool
		decorated bool
		title     string
	)
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {

		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case appx.ConfigEvent:
			decorated = e.Config.Decorated
			title = e.Config.Title
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			for b.Clicked(gtx) {
				toggle = !toggle
				w.Option(appx.Decorated(toggle))
			}

			cl := clip.Rect{Max: e.Size}.Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{A: 0xff, G: 0xff}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, &b, "Toggle decorations")
						btn.TextSize = unit.Sp(42)
						return btn.Layout(gtx)
					}),
					layout.Rigid(material.Body1(th, fmt.Sprintf("Decorated: %v", decorated)).Layout),
				)
			})
			cl.Pop()
			if !decorated {
				w.Perform(deco.Update(gtx))
				material.Decorations(th, &deco, ^system.Action(0), title).Layout(gtx)
			}
			e.Frame(gtx.Ops)
		}
	}
}
