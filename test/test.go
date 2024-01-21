package main

import (
	_ "embed"
	"fmt"

	"github.com/Seikaijyu/nenki.ui/app"
	"github.com/Seikaijyu/nenki.ui/utils"
	"github.com/Seikaijyu/nenki.ui/widget"
)

func main() {
	app.NewApp("测试").Title("Layout").
		Then(func(self *app.App, root *widget.ContainerLayout) {
			self.Background(utils.HexToRGBA("#00ffac0a"))
			root.AppendChild(widget.NewSwitch("测试").OnChange(func(p *widget.Switch, value bool) {
				fmt.Println(value)
			}))
		})

	// 阻塞
	app.Run()
}
