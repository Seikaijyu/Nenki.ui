## Nenki.UI
基于gio 0.4.0（固定版本）实现的更好的golang gui框架

暂无文档，可以根据此代码自行摸索，文档敬请期待
```go
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
						for i := 0; i < 1000000; i++ {
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
```