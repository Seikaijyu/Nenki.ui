package edge

type EdgeDirection uint8

const (
	Top EdgeDirection = iota
	Left
	Bottom
	Right
)

// 所有方向都设置为同一个值
func All(v float32) (float32, float32, float32, float32) {
	return v, v, v, v
}

// 所有方向的值都为零
func Zero() (float32, float32, float32, float32) {
	return 0, 0, 0, 0
}

// 从方向设置值
func (p EdgeDirection) FromDirection(v float32) (float32, float32, float32, float32) {
	switch p {
	case Top:
		return v, 0, 0, 0
	case Left:
		return 0, v, 0, 0
	case Bottom:
		return 0, 0, v, 0
	case Right:
		return 0, 0, 0, v
	}
	return 0, 0, 0, 0
}
