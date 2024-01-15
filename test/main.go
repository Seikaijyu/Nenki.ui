package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	appx "gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	gwidget "gioui.org/widget"
	"gioui.org/widget/material"
	"nenki.ui/app"
	"nenki.ui/utils"
	"nenki.ui/widget"
	"nenki.ui/widget/axis"
)

func main() {
	//Test()
	app.NewApp("测试").Title("Layout").
		Then(func(self *app.App, root *widget.ContainerLayout) {
			self.Background(utils.HexToRGBA("#00ffac0a"))
			root.AppendChild(widget.NewRowLayout().
				Then(func(row *widget.RowLayout) {
					self.Then(func(self *app.App, root *widget.ContainerLayout) {
						list := widget.NewListLayout(axis.Vertical).ScrollMinLen(30)
						cloumn2 := widget.NewColumnLayout()
						row.AppendFlexChild(2.5, widget.NewBorder(list))
						row.AppendFlexChild(6, widget.NewBorder(cloumn2))
						for i := 0; i < 10000; i++ {
							list.AppendChild(widget.NewBorder(widget.NewButton(fmt.Sprintf("Item %d", i)).
								CornerRadius(0).Background(utils.HexToRGBA("#00fff00f")).FontColor(utils.HexToRGBA("#000000"))))
						}
						cloumn2.AppendFlexChild(1, widget.NewBorder(widget.NewContainerLayout()))
						cloumn2.AppendFlexChild(8, widget.NewBorder(widget.NewContainerLayout()))
					})
				}),
			)
		})

	// 阻塞
	app.Run()
}

func Test() {
	go func() {
		w := appx.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	appx.Main()
}

func loop(w *appx.Window) error {
	var a gwidget.Enum = gwidget.Enum{}
	th := material.NewTheme()

	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {

		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,
					layout.Rigid(material.RadioButton(th, &a, "1", "测试1").Layout),
					layout.Rigid(material.RadioButton(th, &a, "2", "测试2").Layout),
				)
			})
			e.Frame(gtx.Ops)
		}
	}
}
