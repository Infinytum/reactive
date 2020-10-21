package reactive

// replaySubject is a special implementation of a subject
// It will always keep the last submitted value and new subscribers
// will receive that value immediately.
type replaySubject struct {
	*bufferSubject
}

// Pipe decorates an observable with one or multiple middlewares
// and returns a new observable with the decoration applied
func (subject *replaySubject) Pipe(fns ...func(Observable, Subjectable)) Observable {
	parent := subject
	for _, fn := range fns {
		sub := NewReplaySubject().(*replaySubject)
		fn(parent, sub)
		parent = sub
	}
	return parent
}

// NewReplaySubject returns a pointer
// to an empty instance of replaySubject
func NewReplaySubject() Subjectable {
	return &replaySubject{
		bufferSubject: NewBufferSubject(1).(*bufferSubject),
	}
}
