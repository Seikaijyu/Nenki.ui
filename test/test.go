package main

import (
	_ "embed"

	"github.com/Seikaijyu/nenki.ui/app"
	"github.com/Seikaijyu/nenki.ui/context"
	"github.com/Seikaijyu/nenki.ui/widget"
	"github.com/Seikaijyu/nenki.ui/widget/anchor"
)

func main() {
	app.NewApp("测试").Title("配置界面").
		Then(func(app *app.App, root *widget.ContainerLayout) {

			app.Background(23, 23, 24, 255)

			root.AppendChild(widget.NewColumnLayout().Then(
				func(self *widget.ColumnLayout) {

					OpenFanUrlButton := widget.NewSwitch()
					self.AppendFlexChild(10, GenericInputBorder(app, "请输入排队关键词"))
					self.AppendFlexAnchorChild(10, anchor.Top, OpenFanUrlButton)
					self.AppendFlexChild(10, GenericInputBorder(app, "请输入您的身份码"))
				},
			),
			)
		})

	// 阻塞
	app.Run()
}
func GenericInputBorder(parent *app.App, text string) *widget.Border {
	InputWidget := widget.NewEditor(text)
	OutInput := widget.NewBorder(InputWidget).Color(255, 255, 255, 255).Margin(0, 0, 5, 0).CornerRadius(5)
	InputWidget.
		HintColor(255, 255, 255, 165).
		OnFocused(func(p *widget.Editor, focus bool, text string) {
			parent.Then(func(self *app.App, root *context.Root) {
				if focus {
					p.HintColor(255, 255, 152, 255)
				} else {
					p.HintColor(255, 255, 255, 165)
				}
			})
		}).
		TextColor(255, 255, 255, 255).Margin(5, 5, 0, 0).LineHeight(50)
	return OutInput
}
