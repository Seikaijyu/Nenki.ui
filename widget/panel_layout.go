package widget

// import (
// 	"image"

// 	glayout "github.com/Seikaijyu/gio/layout"
// 	gclip "github.com/Seikaijyu/gio/op/clip"
// 	gunit "github.com/Seikaijyu/gio/unit"
// )

// // 校验接口是否实现
// var _ WidgetInterface = &Panel{}

// type panelConfig struct {
// 	// 是否更新组件
// 	update bool
// 	// 删除事件
// 	_destroy func()
// }

// type Panel struct {
// 	// 配置
// 	config *panelConfig
// 	// 外边距
// 	margin *glayout.Inset
// 	child  WidgetInterface
// }

// // 绑定函数
// func (p *Panel) Then(fn func(self *Panel)) *Panel {
// 	fn(p)
// 	return p
// }

// // 注销自身，清理所有引用
// func (p *Panel) Destroy() {
// 	p.config.update = false
// 	if p.config._destroy != nil {
// 		p.config._destroy()
// 	}
// 	p.config._destroy = nil
// }

// // 注册删除事件
// func (p *Panel) OnDestroy(fn func()) {
// 	p.config._destroy = fn
// }

// // 是否更新组件
// func (p *Panel) Update(update bool) {
// 	p.config.update = update
// }

// // 外边距
// func (p *Panel) Margin(Top, Left, Bottom, Right float32) *Panel {
// 	p.margin.Top = gunit.Dp(Top)
// 	p.margin.Left = gunit.Dp(Left)
// 	p.margin.Bottom = gunit.Dp(Bottom)
// 	p.margin.Right = gunit.Dp(Right)
// 	return p
// }
// func (p *Panel) Layout(gtx glayout.Context) glayout.Dimensions {
// 	if !p.config.update {
// 		p.config.update = false
// 	}
// 	gtx.Constraints.Max = image.Pt(500, 100)
// 	gtx.Constraints.Min = gtx.Constraints.Max
// 	defer gclip.Rect(image.Rect(0, 0, 500, 100)).Push(gtx.Ops).Pop()
// 	p.child.Layout(gtx)
// 	return glayout.Dimensions{Size: image.Pt(0, 0), Baseline: 0}

// }

// func NewPanel(widget WidgetInterface) *Panel {
// 	return &Panel{
// 		child:  widget,
// 		config: &panelConfig{},
// 	}
// }
