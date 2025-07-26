package event

type Observer[T Event] interface {
	Update(EventWrapper[T]) error
}
