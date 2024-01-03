package widget

// 容器布局
type ContainerLayout = AnchorLayout

// 创建容器
func NewContainerLayout() *ContainerLayout {
	return NewAnchorLayout(Center)
}
