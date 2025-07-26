package event

type Notifier[T Event] interface {
	AddObserver(Observer[T])
	RemoveObserver(Observer[T])
	Notify(EventWrapper[T])
}

type BaseNotifier[T Event] struct {
	observers []Observer[T]
}

func NewBaseNotifier[T Event]() *BaseNotifier[T] {
	return &BaseNotifier[T]{
		observers: make([]Observer[T], 0),
	}
}

func (s *BaseNotifier[T]) AddObserver(observer Observer[T]) {
	s.observers = append(s.observers, observer)
}

func (s *BaseNotifier[T]) RemoveObserver(observer Observer[T]) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *BaseNotifier[T]) Notify(event EventWrapper[T]) {
	for _, obs := range s.observers {
		obs.Update(event)
	}
}
