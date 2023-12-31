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
	// 添加子节点
	AppendChild(child ...WidgetInterface) T
	// 添加子节点并设置ID
	AppendChildWithID(id string, child ...WidgetInterface) T
	// 从ID删除子节点
	RemoveChildAtID(id string) T
	// 删除所有子节点
	RemoveChildAll() T
	// 获取子节点
	Childs() []WidgetInterface
	// 获取子节点数量
	ChildCount() int
	// 从索引获取子节点
	ChildAt(index int) WidgetInterface
	// 从ID获取子节点
	ChildAtID(id string) []WidgetInterface
}

// 单子节点布局接口
type SingleChildLayoutInterface[T any] interface {
	WidgetInterface
	// 设置子节点
	AppendChild(childs WidgetInterface) T
	// 获取子节点
	Child() WidgetInterface
	// 删除子节点
	RemoveChild() T
}
