package main

import (
	_ "embed"
	"fmt"

	"github.com/Seikaijyu/nenki.ui/app"
	"github.com/Seikaijyu/nenki.ui/context"
	"github.com/Seikaijyu/nenki.ui/widget"
	"github.com/Seikaijyu/nenki.ui/widget/axis"
	"github.com/Seikaijyu/nenki.ui/widget/edge"
)

func main() {
	app.NewApp("测试").Title("配置界面").
		Then(func(app *app.App, root *widget.ContainerLayout) {
			app.Background(23, 23, 24, 255)
			root.AppendChild(widget.NewSlider(axis.Vertical).Margin(edge.All(20)).OnDragging(func(p *widget.Slider, value float32) {
				fmt.Println(value)
			}))
		})

	// 阻塞
	app.Run()
}
func GenericInputBorder(parent *app.App, text string) *widget.Border {
	InputWidget := widget.NewEditor(text)
	OutInput := widget.NewBorder(InputWidget).Color(255, 255, 255, 255).Margin(0, 0, 5, 0).CornerRadius(10)
	InputWidget.
		HintColor(255, 255, 255, 165).
		OnFocused(func(p *widget.Editor, focus bool, text string) {
			parent.Then(func(self *app.App, root *context.Root) {
				if focus {
					OutInput.Color(255, 255, 152, 255)
				} else {
					OutInput.Color(255, 255, 255, 165)
				}
			})
		}).
		TextColor(255, 255, 255, 255).Margin(5, 5, 0, 0).LineHeight(50)
	return OutInput
}
