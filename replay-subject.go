package reactive

// ReplaySubject is a special implementation of a subject
// It will always keep the last submitted value and new subscribers
// will receive that value immediately.
type ReplaySubject struct {
	BufferSubject
}

// Pipe decorates an observable with one or multiple middlewares
// and returns a new observable with the decoration applied
func (subject *ReplaySubject) Pipe(fns ...func(Observable, Subjectable)) Observable {
	parent := subject
	for _, fn := range fns {
		sub := NewReplaySubject()
		fn(parent, sub)
		parent = sub
	}
	return parent
}

// NewReplaySubject returns a pointer
// to an empty instance of ReplaySubject
func NewReplaySubject() *ReplaySubject {
	return &ReplaySubject{
		BufferSubject: *NewBufferSubject(1),
	}
}
