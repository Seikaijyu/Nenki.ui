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
	"gioui.org/widget"
	"gioui.org/widget/material"
	"nenki.ui/app"
)

func main() {
	//Test()
	app.NewApp("测试").SetNavigationColor(255, 255, 39, 255).SetMinSize(1024, 1024).SetSize(1024, 1024).SetMaxSize(768, 768)
	// 阻塞
	app.Run()
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
		b    widget.Clickable
		deco widget.Decorations
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
					layout.Rigid(material.Button(th, &b, "Toggle decorations").Layout),
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
