package event

type FakeObserver[T Event] struct {
	events []EventWrapper[T]
}

func NewFakeObserver[T Event]() *FakeObserver[T] {
	return &FakeObserver[T]{
		events: make([]EventWrapper[T], 0),
	}
}

func (f *FakeObserver[T]) Update(event EventWrapper[T]) error {
	f.events = append(f.events, event)
	return nil
}

func (f *FakeObserver[T]) GetEvents() []EventWrapper[T] {
	return f.events
}
