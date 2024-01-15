package context

// Queue 表示任意类型 T 的队列。
type Queue[T any] struct {
	data []T
}

// Enqueue 将元素添加到队列的末尾。
func (q *Queue[T]) Enqueue(v T) {
	q.data = append(q.data, v)
}

// Dequeue 从队列中移除并返回第一个元素。
// 如果队列为空，则返回 false。
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.data) == 0 {
		var zero T // 创建类型 T 的零值
		return zero, false
	}
	v := q.data[0]
	q.data = q.data[1:]
	return v, true
}
