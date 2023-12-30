package monad

type MonadInterface[T any, V any] interface {
	// Monad模式必要的绑定，接受一个函数，翻译执行后得到的Monad结构体
	Bind(func(T) T) T
	// Monad模式必要的包装，应该等同于NewMonad函数，用于把外部值包装为Monad
	Unit(value V) T
}

// Monad接口实现例子，所有结构体都必须实现MonadInterface
var _ MonadInterface[*Monad, int] = &Monad{}

// Monad结构体例子
type Monad struct {
	value *int
}

func (p *Monad) Bind(fn func(*Monad) *Monad) *Monad {
	return fn(p)
}

func (p *Monad) SetValue(value int) *Monad {
	p.value = &value
	return p
}

func (p *Monad) Unit(value int) *Monad {
	return NewMonad(value)
}
func (p *Monad) Value() int {
	return *p.value
}
func NewMonad(v int) *Monad {
	return &Monad{value: &v}
}
