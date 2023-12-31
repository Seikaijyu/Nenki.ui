package widget

// 组件池，用于储存全局组件并索引组件ID
type widgetPool map[string][]WidgetInterface

// widgetPool 是一个全局的组件池，用于储存全局组件并索引组件ID
var pool = &widgetPool{}

// AddWidget 将一个小部件添加到具有指定ID的小部件池中。
// 如果ID在池中不存在，则创建一个新条目。
// 小部件将追加到与该ID关联的小部件列表中。
func (p widgetPool) AddWidget(id string, widget WidgetInterface) int {
	if _, ok := p[id]; !ok {
		p[id] = []WidgetInterface{}
	}
	p[id] = append(p[id], widget)
	// 返回对应的索引
	return len(p[id]) - 1
}

// GetWidgets 返回与指定ID关联的小部件列表。
// 如果ID在池中不存在，则返回nil。
func (p widgetPool) GetWidgets(id string) []WidgetInterface {
	if _, ok := p[id]; !ok {
		return nil
	}
	return p[id]
}

// GetWidgetAtIndex 从索引获取小部件
func (p widgetPool) GetWidgetAtIndex(id string, index int) WidgetInterface {
	if _, ok := p[id]; !ok {
		return nil
	}
	if index < 0 || index >= len(p[id]) {
		return nil
	}
	return p[id][index]
}

// Remove 从与指定ID关联的小部件列表中删除所有小部件。
// 如果ID在池中不存在，则不执行任何操作。
// 注意：该函数会将与指定ID关联的小部件列表置为nil
func (p widgetPool) Remove(id string) {
	if _, ok := p[id]; !ok {
		return
	}
	p[id] = nil
}

// RemoveAtIndex 从与指定ID关联的小部件列表中删除具有指定索引的小部件。
// 如果ID在池中不存在，则不执行任何操作。
// 如果索引超出范围，则不执行任何操作。
// 注意：该函数只会将指定索引位置的小部件置为nil，并不会移动其他小部件的位置。
func (p widgetPool) RemoveAtIndex(id string, index int) {
	if _, ok := p[id]; !ok {
		return
	}
	if index < 0 || index >= len(p[id]) {
		return
	}
	p[id][index] = nil
}

// RemoveAll 从 WidgetPool 中移除所有的小部件。
func (p widgetPool) RemoveAll() {
	for id := range p {
		delete(p, id)
	}
}
