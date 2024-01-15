package widget

import (
	glayout "github.com/Seikaijyu/gio/layout"
)

// 定义一个通用的布局接口
type WidgetInterface interface {
	// 渲染UI
	Layout(gtx glayout.Context) glayout.Dimensions
	//  删除函数注册
	OnDestroy(fn func())
	// 删除组件
	Destroy()

	// 是否更新组件
	Update(update bool)
}

// 多子节点布局接口
type MultiChildLayoutInterface[T any] interface {
	WidgetInterface
	// 重新设置父节点
	ResetParent(child WidgetInterface)
	// 外边距
	Margin(Top, Left, Bottom, Right float32) T
	// 从指定索引删除子节点
	RemoveChildAt(index int) T
	// 从组件删除子节点
	RemoveChild(child WidgetInterface) T
	// 删除所有子节点
	RemoveChildAll() T
	// 获取子节点
	GetChildAll() []WidgetInterface
	// 获取指定索引的子节点
	GetChildAt(index int) WidgetInterface
	// 获取子节点数量
	GetChildCount() int
}

// 单子节点布局接口
type SingleChildLayoutInterface[T any] interface {
	WidgetInterface
	// 重新设置父节点
	ResetParent(child WidgetInterface)
	// 外边距
	Margin(Top, Left, Bottom, Right float32) T
	// 设置子节点
	AppendChild(childs WidgetInterface) T
	// 获取子节点
	GetChild() WidgetInterface
	// 删除子节点
	RemoveChild() T
}
