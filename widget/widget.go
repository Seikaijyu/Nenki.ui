package widget

import (
	glayout "gioui.org/layout"
)

// 定义一个通用的布局接口
type WidgetInterface interface {
	// 渲染UI
	Layout(gtx glayout.Context) glayout.Dimensions
	// 是否被删除
	IsDestroy() bool
	// 注销自身，清理所有引用
	Destroy()
}

// 多子节点布局接口
type MultiChildLayoutInterface[T any] interface {
	WidgetInterface
	// 外边距
	Margin(Top, Left, Bottom, Right float32) T
	// 添加子节点
	AppendChild(child ...WidgetInterface) T
	// 从指定索引删除子节点
	RemoveChildAt(index int) T
	// 删除所有子节点
	RemoveChildAll() T
	// 获取子节点
	ChildAll() []WidgetInterface
	// 获取指定索引的子节点
	ChildAt(index int) WidgetInterface
	// 获取子节点数量
	ChildCount() int
}

// 单子节点布局接口
type SingleChildLayoutInterface[T any] interface {
	WidgetInterface
	// 外边距
	Margin(Top, Left, Bottom, Right float32) T
	// 设置子节点
	AppendChild(childs WidgetInterface) T
	// 获取子节点
	Child() WidgetInterface
	// 删除子节点
	RemoveChild() T
}
