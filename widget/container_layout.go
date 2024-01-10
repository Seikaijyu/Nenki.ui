package widget

import glayout "gioui.org/layout"

// 容器布局
type ContainerLayout struct {
	// 居中锚定布局
	centerAnchorLayout *AnchorLayout
}

var _ SingleChildLayoutInterface[*ContainerLayout] = &ContainerLayout{}

func NewContainerLayout() *ContainerLayout {

	return &ContainerLayout{
		centerAnchorLayout: NewAnchorLayout(Center),
	}
}

// 绑定函数
func (p *ContainerLayout) Then(fn func(*ContainerLayout)) *ContainerLayout {
	fn(p)
	return p
}

// 设置子节点
func (p *ContainerLayout) AppendChild(child WidgetInterface) *ContainerLayout {
	p.centerAnchorLayout.AppendChild(child)
	return p
}

// 获取子节点
func (p *ContainerLayout) Child() WidgetInterface {
	return p.centerAnchorLayout.Child()
}

// 删除子节点
func (p *ContainerLayout) RemoveChild() *ContainerLayout {
	p.centerAnchorLayout.RemoveChild()
	return p
}

// 是否被删除
func (p *ContainerLayout) IsDestroy() bool {
	return p.centerAnchorLayout.IsDestroy()
}

// 删除自身
func (p *ContainerLayout) Destroy() {
	p.centerAnchorLayout.Destroy()
}

// 渲染UI
func (p *ContainerLayout) Layout(gtx glayout.Context) glayout.Dimensions {
	return p.centerAnchorLayout.Layout(gtx)
}
