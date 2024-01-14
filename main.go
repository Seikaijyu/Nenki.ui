// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

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
	"nenki.ui/widget/edge"
)

func main() {

	//Test()
	app.NewApp("测试").Size(600, 1000).MinSize(600, 1000).MaxSize(600, 1000).Title("测试窗口").DragFiles(true).
		Then(func(self *app.App, root *widget.ContainerLayout) {

			cloumn := widget.NewColumnLayout()
			root.AppendChild(cloumn)
			editor := widget.NewEditor("请随便输入什么文字").Then(func(self *widget.Editor) {
				self.Submit(true).MaxLines(1).FontSize(20).Margin(edge.All(10))
				self.OnSubmit(func(text string) {
					fmt.Println("回车了", text)
				})
			})
			submit := widget.NewButton("提交").Margin(edge.All(10)).CornerRadius(0)
			cloumn.AppendFlexChild(0.1, widget.NewContainerLayout().Background(utils.HexToRGBA("#0fc0ff")))
			cloumn.AppendFlexChild(1, widget.NewBorder(editor).Margin(edge.All(1)))
			cloumn.AppendFlexChild(0.1, submit)
		})
	// 阻塞
	app.Run()
}

func longclick(b *widget.Button) {
	fmt.Println("双击了按钮")
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
	var a gwidget.Clickable = gwidget.Clickable{}
	th := material.NewTheme()

	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {

		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						return material.Button(th, &a, "测试").Layout(gtx)
					}),
				)
			}))
			e.Frame(gtx.Ops)
		}
	}
}
