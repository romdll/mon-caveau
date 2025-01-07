package injector

type Injectable interface {
	Preload() error
	Process() error
	Display()
}

type GenericWrapper[T any] struct {
	injector *Injector[T]
}

func NewGenericWrapper[T any](data []T, keyFunc EntityToMapKeyFunc[T]) *GenericWrapper[T] {
	return &GenericWrapper[T]{
		injector: NewInjector(data, keyFunc),
	}
}

func (gw *GenericWrapper[T]) Preload() error {
	return gw.injector.Preload()
}

func (gw *GenericWrapper[T]) Process() error {
	return gw.injector.Process()
}

func (gw *GenericWrapper[T]) Display() {
	gw.injector.logInjectionStats()
}
