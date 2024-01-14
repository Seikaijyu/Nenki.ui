// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image"
	"image/color"
	"math"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

// FromListPosition将一个layout.Position转换为两个浮点数，这两个浮点数表示视口在基础内容上的位置。它需要知道列表中的元素个数和列表的主轴大小才能做到这一点。返回的值将在 [0,1] 的范围内，并且start将小于或等于end。
func FromListPosition(lp layout.Position, elements int, majorAxisSize int) (start, end float32) {
	return fromListPosition(lp, elements, majorAxisSize)
}

// fromListPosition将一个layout.Position转换为两个浮点数，这两个浮点数表示视口在基础内容上的位置。它需要知道列表中的元素个数和列表的主轴大小才能做到这一点。返回的值将在 [0,1] 的范围内，并且start将小于或等于end。
func fromListPosition(lp layout.Position, elements int, majorAxisSize int) (start, end float32) {
	// 估算可滚动内容的大小。
	lengthEstPx := float32(lp.Length)
	elementLenEstPx := lengthEstPx / float32(elements)

	// 确定可见内容的比例。
	listOffsetF := float32(lp.Offset)
	listOffsetL := float32(lp.OffsetLast)

	// 使用估计的元素大小和已知的像素偏移计算视口开始位置。
	viewportStart := clamp1((float32(lp.First)*elementLenEstPx + listOffsetF) / lengthEstPx)
	viewportEnd := clamp1((float32(lp.First+lp.Count)*elementLenEstPx + listOffsetL) / lengthEstPx)
	viewportFraction := viewportEnd - viewportStart

	// 仅根据可见大小和估计的总大小之比，计算列表内容的预期可见比例。
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthEstPx

	// 计算确定视口的两种方法之间的误差，并根据我们接近每个端点的程度，在视口的两端扩散误差。
	err := visibleFraction - viewportFraction
	adjStart := viewportStart
	adjEnd := viewportEnd
	if viewportFraction < 1 {
		startShare := viewportStart / (1 - viewportFraction)
		endShare := (1 - viewportEnd) / (1 - viewportFraction)
		startErr := startShare * err
		endErr := endShare * err

		adjStart -= startErr
		adjEnd += endErr
	}
	return adjStart, adjEnd
}

// rangeIsScrollable returns whether the viewport described by start and end
// is smaller than the underlying content (such that it can be scrolled).
// start and end are expected to each be in the range [0,1], and start
// must be less than or equal to end.
func rangeIsScrollable(start, end float32) bool {
	return end-start < 1
}

// ScrollTrackStyle configures the presentation of a track for a scroll area.
type ScrollTrackStyle struct {
	// MajorPadding and MinorPadding along the major and minor axis of the
	// scrollbar's track. This is used to keep the scrollbar from touching
	// the edges of the content area.
	MajorPadding, MinorPadding unit.Dp
	// Color of the track background.
	Color color.NRGBA
}

// ScrollIndicatorStyle configures the presentation of a scroll indicator.
type ScrollIndicatorStyle struct {
	// MajorMinLen is the smallest that the scroll indicator is allowed to
	// be along the major axis.
	MajorMinLen unit.Dp
	// MinorWidth is the width of the scroll indicator across the minor axis.
	MinorWidth unit.Dp
	// Color and HoverColor are the normal and hovered colors of the scroll
	// indicator.
	Color, HoverColor color.NRGBA
	// CornerRadius is the corner radius of the rectangular indicator. 0
	// will produce square corners. 0.5*MinorWidth will produce perfectly
	// round corners.
	CornerRadius unit.Dp
}

// ScrollbarStyle configures the presentation of a scrollbar.
type ScrollbarStyle struct {
	Scrollbar *widget.Scrollbar
	Track     ScrollTrackStyle
	Indicator ScrollIndicatorStyle
}

// Scrollbar configures the presentation of a scrollbar using the provided
// theme and state.
func Scrollbar(th *Theme, state *widget.Scrollbar) ScrollbarStyle {
	lightFg := th.Palette.Fg
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200

	return ScrollbarStyle{
		Scrollbar: state,
		Track: ScrollTrackStyle{
			MajorPadding: 2,
			MinorPadding: 2,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  th.FingerSize,
			MinorWidth:   6,
			CornerRadius: 3,
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
}

// Width 函数返回当前配置下滚动条的次要轴宽度（考虑到滚动轨道的填充）。
func (s ScrollbarStyle) Width() unit.Dp {
	return s.Indicator.MinorWidth + s.Track.MinorPadding + s.Track.MinorPadding
}

// Layout 函数布局滚动条。
func (s ScrollbarStyle) Layout(gtx layout.Context, axis layout.Axis, viewportStart, viewportEnd float32) layout.Dimensions {
	if !rangeIsScrollable(viewportStart, viewportEnd) {
		return layout.Dimensions{}
	}

	// 以与轴无关的方式设置最小约束，然后转换为当前轴的正确表示。
	convert := axis.Convert
	maxMajorAxis := convert(gtx.Constraints.Max).X
	gtx.Constraints.Min.X = maxMajorAxis
	gtx.Constraints.Min.Y = gtx.Dp(s.Width())
	gtx.Constraints.Min = convert(gtx.Constraints.Min)
	gtx.Constraints.Max = gtx.Constraints.Min

	s.Scrollbar.Update(gtx, axis, viewportStart, viewportEnd)

	// 如果鼠标悬停，则变暗指示器。
	if s.Scrollbar.IndicatorHovered() {
		s.Indicator.Color = s.Indicator.HoverColor
	}

	return s.layout(gtx, axis, viewportStart, viewportEnd)
}

// layout the scroll track and indicator.
func (s ScrollbarStyle) layout(gtx layout.Context, axis layout.Axis, viewportStart, viewportEnd float32) layout.Dimensions {
	inset := layout.Inset{
		Top:    s.Track.MajorPadding,
		Bottom: s.Track.MajorPadding,
		Left:   s.Track.MinorPadding,
		Right:  s.Track.MinorPadding,
	}
	if axis == layout.Horizontal {
		inset.Top, inset.Bottom, inset.Left, inset.Right = inset.Left, inset.Right, inset.Top, inset.Bottom
	}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			// Lay out the draggable track underneath the scroll indicator.
			area := image.Rectangle{
				Max: gtx.Constraints.Min,
			}
			pointerArea := clip.Rect(area)
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddDrag(gtx.Ops)

			// Stack a normal clickable area on top of the draggable area
			// to capture non-dragging clicks.
			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddTrack(gtx.Ops)

			paint.FillShape(gtx.Ops, s.Track.Color, clip.Rect(area).Op())
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Use axis-independent constraints.
				gtx.Constraints.Min = axis.Convert(gtx.Constraints.Min)
				gtx.Constraints.Max = axis.Convert(gtx.Constraints.Max)

				// Compute the pixel size and position of the scroll indicator within
				// the track.
				trackLen := gtx.Constraints.Min.X
				viewStart := int(math.Round(float64(viewportStart) * float64(trackLen)))
				viewEnd := int(math.Round(float64(viewportEnd) * float64(trackLen)))
				indicatorLen := max(viewEnd-viewStart, gtx.Dp(s.Indicator.MajorMinLen))
				if viewStart+indicatorLen > trackLen {
					viewStart = trackLen - indicatorLen
				}
				indicatorDims := axis.Convert(image.Point{
					X: indicatorLen,
					Y: gtx.Dp(s.Indicator.MinorWidth),
				})
				radius := gtx.Dp(s.Indicator.CornerRadius)

				// Lay out the indicator.
				offset := axis.Convert(image.Pt(viewStart, 0))
				defer op.Offset(offset).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, s.Indicator.Color, clip.RRect{
					Rect: image.Rectangle{
						Max: indicatorDims,
					},
					SW: radius,
					NW: radius,
					NE: radius,
					SE: radius,
				}.Op(gtx.Ops))

				// Add the indicator pointer hit area.
				area := clip.Rect(image.Rectangle{Max: indicatorDims})
				defer pointer.PassOp{}.Push(gtx.Ops).Pop()
				defer area.Push(gtx.Ops).Pop()
				s.Scrollbar.AddIndicator(gtx.Ops)

				return layout.Dimensions{Size: axis.Convert(gtx.Constraints.Min)}
			})
		},
	)
}

// AnchorStrategy defines a means of attaching a scrollbar to content.
type AnchorStrategy uint8

const (
	// Occupy 预留空间给滚动条，使得下面的内容区域在一个轴上变小。
	Occupy AnchorStrategy = iota
	// Overlay 让滚动条浮动在内容上方，不占用任何空间。位于下面的内容可能会被滚动条遮挡。
	Overlay
)

// ListStyle configures the presentation of a layout.List with a scrollbar.
type ListStyle struct {
	state *widget.List
	ScrollbarStyle
	AnchorStrategy
}

// List constructs a ListStyle using the provided theme and state.
func List(th *Theme, state *widget.List) ListStyle {
	return ListStyle{
		state:          state,
		ScrollbarStyle: Scrollbar(th, &state.Scrollbar),
	}
}

// 布局列表和滚动条。
func (l ListStyle) Layout(gtx layout.Context, length int, w layout.ListElement) layout.Dimensions {
	originalConstraints := gtx.Constraints

	// 确定滚动条占用多少空间。
	barWidth := gtx.Dp(l.Width())

	if l.AnchorStrategy == Occupy {
		// 使用gtx约束为滚动条预留空间。
		max := l.state.Axis.Convert(gtx.Constraints.Max)
		min := l.state.Axis.Convert(gtx.Constraints.Min)
		max.Y -= barWidth
		if max.Y < 0 {
			max.Y = 0
		}
		min.Y -= barWidth
		if min.Y < 0 {
			min.Y = 0
		}
		gtx.Constraints.Max = l.state.Axis.Convert(max)
		gtx.Constraints.Min = l.state.Axis.Convert(min)
	}

	listDims := l.state.List.Layout(gtx, length, w)
	gtx.Constraints = originalConstraints

	// 绘制滚动条
	anchoring := layout.E // layout.Right
	if l.state.Axis == layout.Horizontal {
		anchoring = layout.S // layout.Bottom
	}
	majorAxisSize := l.state.Axis.Convert(listDims.Size).X
	start, end := fromListPosition(l.state.Position, length, majorAxisSize)
	// layout.Direction尊重最小值，所以要确保即使提供的layout.Context有零最小约束，滚动条也会在正确的边缘绘制。
	gtx.Constraints.Min = listDims.Size
	if l.AnchorStrategy == Occupy {
		min := l.state.Axis.Convert(gtx.Constraints.Min)
		min.Y += barWidth
		gtx.Constraints.Min = l.state.Axis.Convert(min)
	}
	anchoring.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return l.ScrollbarStyle.Layout(gtx, l.state.Axis, start, end)
	})

	if delta := l.state.ScrollDistance(); delta != 0 {
		// 处理用户与滚动条交互导致的列表位置变化。
		l.state.List.ScrollBy(delta * float32(length))
	}

	if l.AnchorStrategy == Occupy {
		// 增加宽度以计算滚动条占用的空间。
		cross := l.state.Axis.Convert(listDims.Size)
		cross.Y += barWidth
		listDims.Size = l.state.Axis.Convert(cross)
	}

	return listDims
}
