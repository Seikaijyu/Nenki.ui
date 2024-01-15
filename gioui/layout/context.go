// SPDX-License-Identifier: Unlicense OR MIT

package layout

import (
	"image"
	"time"

	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
)

// Context carries the state needed by almost all layouts and widgets.
// A zero value Context never returns events, map units to pixels
// with a scale of 1.0, and returns the zero time from Now.
type Context struct {
	// Constraints track the constraints for the active widget or
	// layout.
	Constraints Constraints

	Metric unit.Metric
	// By convention, a nil Queue is a signal to widgets to draw themselves
	// in a disabled state.
	Queue event.Queue
	// Now is the animation time.
	Now time.Time

	// Locale provides information on the system's language preferences.
	// BUG(whereswaldon): this field is not currently populated automatically.
	// Interested users must look up and populate these values manually.
	Locale system.Locale

	*op.Ops
}

// NewContext is a shorthand for
//
//	Context{
//	  Ops: ops,
//	  Now: e.Now,
//	  Queue: e.Queue,
//	  Config: e.Config,
//	  Constraints: Exact(e.Size),
//	}
//
// NewContext calls ops.Reset and adjusts ops for e.Insets.
func NewContext(ops *op.Ops, e system.FrameEvent) Context {
	ops.Reset()

	size := e.Size

	if e.Insets != (system.Insets{}) {
		left := e.Metric.Dp(e.Insets.Left)
		top := e.Metric.Dp(e.Insets.Top)
		op.Offset(image.Point{
			X: left,
			Y: top,
		}).Add(ops)

		size.X -= left + e.Metric.Dp(e.Insets.Right)
		size.Y -= top + e.Metric.Dp(e.Insets.Bottom)
	}

	return Context{
		Ops:         ops,
		Now:         e.Now,
		Queue:       e.Queue,
		Metric:      e.Metric,
		Constraints: Exact(size),
	}
}

// Dp 函数将单位 Dp 转换为像素。
func (c Context) Dp(v unit.Dp) int {
	return c.Metric.Dp(v)
}

// Sp 函数将单位 Sp 转换为像素。
func (c Context) Sp(v unit.Sp) int {
	return c.Metric.Sp(v)
}

// Events 返回可用于键的事件。如果没有配置队列，Events 返回 nil。
func (c Context) Events(k event.Tag) []event.Event {
	if c.Queue == nil {
		return nil
	}
	return c.Queue.Events(k)
}

// Disabled 返回此上下文的一个副本，副本中的队列为 nil，可以阻止事件传递到使用它的小部件。
//
// 按照惯例，nil 队列是指示小部件以禁用状态绘制自身的信号。
func (c Context) Disabled() Context {
	c.Queue = nil
	return c
}
