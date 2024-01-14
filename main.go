package main

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
	"nenki.ui/widget"
	"nenki.ui/widget/edge"
)

func main() {

	//Test()
	app.NewApp("测试").Size(600, 160).MinSize(600, 160).MaxSize(600, 160).Title("测试窗口").Decorated(false).DragFiles(true).
		Then(func(self *app.App, root *widget.ContainerLayout) {
			cloumn := widget.NewColumnLayout()
			root.AppendChild(cloumn).Margin(edge.All(10))
			editor := widget.NewEditor("请随便输入什么文字").Then(func(self *widget.Editor) {
				self.SingleLine(true).Submit(true).FontSize(20).Margin(edge.All(10)).Then(func(self *widget.Editor) {
					self.OnSubmit(func(text string) {
						self.Text("")
						self.Focus()
					})
				})
			})
			cloumn.AppendRigidChild(widget.NewBorder(editor).Margin(edge.Bottom.FromDirection(5)))
			cloumn.AppendRigidChild(widget.NewButton("提交").CornerRadius(0).Then(func(self *widget.Button) {
				self.OnClicked(func(b *widget.Button) {
					editor.Text("")
					editor.Focus()
				})
			}))
			cloumn.AppendRigidChild(widget.NewButton("关闭").Margin(edge.Top.FromDirection(5)).CornerRadius(0).Then(func(self *widget.Button) {
				self.OnClicked(func(b *widget.Button) {
					go func() {
						time.Sleep(time.Second / 3)
						app.Exit()
					}()
				})
			}))
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
