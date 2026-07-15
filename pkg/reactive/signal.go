// Package reactive contains a minimal implementation of a
// signal-based reactivity system — the same pattern used internally
// by frameworks like SolidJS and Vue 3, implemented here in plain Go,
// with no dependency on any JS runtime.
package reactive

// Signal represents an observable value. Whenever the value changes
// via Set, all registered listeners are notified — this is what
// triggers the UI re-render.
type Signal[T any] struct {
	value     T
	listeners []func()
}

// NewSignal creates a new Signal with an initial value.
func NewSignal[T any](initial T) *Signal[T] {
	return &Signal[T]{value: initial}
}

// Get returns the current value of the signal.
func (s *Signal[T]) Get() T {
	return s.value
}

// Set updates the value and notifies all registered listeners.
func (s *Signal[T]) Set(v T) {
	s.value = v
	for _, listener := range s.listeners {
		listener()
	}
}

// Subscribe registers a function to be called whenever the value
// changes. In this MVP, the most common "listener" is a function
// that marks the UI as needing a redraw.
func (s *Signal[T]) Subscribe(listener func()) {
	s.listeners = append(s.listeners, listener)
}
