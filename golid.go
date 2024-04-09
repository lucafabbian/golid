package golid

type SignalContext struct {
	context []*Effect
}

type Signal[T any] struct {
	sc            *SignalContext
	value         T
	subscriptions map[*Effect]bool
}

type Effect struct {
	execute func()
}

func NewSignalContext() *SignalContext {
	return &SignalContext{
		context: make([]*Effect, 0),
	}
}

func (sc *SignalContext) Computed(fn func()) {
	var effect *Effect
	effect = &Effect{execute: func() {
		sc.context = append(sc.context, effect)
		fn()
		sc.context = sc.context[:len(sc.context)-1]
	}}

	effect.execute()
}

func NewSignal[T any](sc *SignalContext, value T) *Signal[T] {
	return &Signal[T]{
		sc:            sc,
		value:         value,
		subscriptions: make(map[*Effect]bool),
	}
}

func (s *Signal[T]) Get() T {
	if len(s.sc.context) > 0 {
		observer := s.sc.context[len(s.sc.context)-1]
		s.subscriptions[observer] = true
	}
	return s.value
}

func (s *Signal[T]) Set(newVal T) {
	s.value = newVal
	for observer := range s.subscriptions {
		observer.execute()
	}
}

func Extract[T any](x interface{}, fn func()) (T, func()) {
	switch v := x.(type) {
	case *Signal[T]:
		if fn != nil {
			effect := &Effect{
				execute: fn,
			}
			v.subscriptions[effect] = true
			return v.Get(), func() {
				delete(v.subscriptions, effect)
			}
		}
		return v.Get(), func() {}
	default:
		return v.(T), func() {}
	}
}
