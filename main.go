package main

import (
	_ "embed"
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
)

func main() {
	Test()
	// app.NewApp("测试").Size(1024, 1024).Title("你好").DragFiles(true).
	// 	Then(func(self *app.App, root *widget.ContainerLayout) {
	// 		h := widget.NewColumnLayout()
	// 		root.AppendChild(h)
	// 		h.AppendRigidChild(widget.NewCheckBox("确实是这样的"))
	// 		h.AppendRigidChild(widget.NewCheckBox("并不是这样的"))

	// 	})

	// // 阻塞
	// app.Run()
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
