// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	appx "gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	gwidget "gioui.org/widget"
	"gioui.org/widget/material"
	"nenki.ui/app"
	"nenki.ui/context"
	"nenki.ui/utils"
	"nenki.ui/widget"
	"nenki.ui/widget/anchor"
)

func main() {

	//Test()
	app.NewApp("测试").Size(1024, 1024).Title("你好").DragFiles(true).
		Then(func(self *app.App, root *widget.ContainerLayout) {
			self.Background(utils.HexToRGBA("#efaf00"))
			v := widget.NewButton("测网速").FontSize(200)
			b := widget.NewBorder(v)
			h := widget.NewAnchorLayout(anchor.Top).AppendChild(b)

			root.AppendChild(h)
			go func() {
				time.Sleep(2 * time.Second)
				self.Then(func(self *app.App, root *context.Root) {
					v.Text("测网速2")
				})
				time.Sleep(2 * time.Second)
				self.Then(func(self *app.App, root *context.Root) {
					v.Text("测网速333")
				})
			}()
		})
	// 阻塞
	app.Run()
}

func longclick(b *widget.Button) {
	fmt.Println("双击了按钮")
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
	btn := material.Button(th, &b, "Toggle decorations")

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

			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

				return layout.Inset{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

					return btn.Layout(gtx)
				})

			})
			if !decorated {
				w.Perform(deco.Update(gtx))
				material.Decorations(th, &deco, ^system.Action(0), title).Layout(gtx)
			}
			e.Frame(gtx.Ops)
		}
	}
}
