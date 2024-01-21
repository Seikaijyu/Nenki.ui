package widget

// import (
// 	glayout "github.com/Seikaijyu/gio/layout"
// 	gunit "github.com/Seikaijyu/gio/unit"
// )

// // 校验接口是否实现
// var _ WidgetInterface = &Caption{}

// type captionConfig struct {
// 	// 是否更新组件
// 	update bool
// 	// 删除事件
// 	_destroy func()
// }

// type Caption struct {
// 	// 配置
// 	config *captionConfig
// 	// 外边距
// 	margin *glayout.Inset
// }

// // 绑定函数
// func (p *Caption) Then(fn func(self *Caption)) *Caption {
// 	fn(p)
// 	return p
// }

// // 注销自身，清理所有引用
// func (p *Caption) Destroy() {
// 	p.config.update = false
// 	if p.config._destroy != nil {
// 		p.config._destroy()
// 	}
// 	p.config._destroy = nil
// }

// // 注册删除事件
// func (p *Caption) OnDestroy(fn func()) {
// 	p.config._destroy = fn
// }

// // 是否更新组件
// func (p *Caption) Update(update bool) {
// 	p.config.update = update
// }

// // 外边距
// func (p *Caption) Margin(Top, Left, Bottom, Right float32) *Caption {
// 	p.margin.Top = gunit.Dp(Top)
// 	p.margin.Left = gunit.Dp(Left)
// 	p.margin.Bottom = gunit.Dp(Bottom)
// 	p.margin.Right = gunit.Dp(Right)
// 	return p
// }
// func (p *Caption) Layout(gtx glayout.Context) glayout.Dimensions {
// if !p.config.update {
// 	p.config.update = false
// }
// 	return gmaterial.
// }

// func NewCaption() *Caption {
// 	return &Caption{
// 		config: &captionConfig{},
// 	}
// }
