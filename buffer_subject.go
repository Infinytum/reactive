package reactive

import "github.com/infinytum/reactive/internal"

// BufferSubject is a special implementation of a subject
// It will always keep the submitted values until the buffer is full
// and new subscribers will receive those values immediately.
//
// New values will remove the oldest value from the buffer
type BufferSubject struct {
	LastValues [][]interface{}
	Subject
}

// AsChannel returns a channel which will receive all
// further updates of this observable
func (subject *BufferSubject) AsChannel() chan []interface{} {
	channel := subject.Subject.AsChannel()
	go func(channel chan []interface{}) {
		data := subject.LastValues
		internal.Reverse(data)
		for _, valueArray := range data {
			if valueArray != nil {
				channel <- valueArray
			}
		}

	}(channel)
	return channel
}

// Close will remove all subscribers and render
// the subjectable useless
func (subject *BufferSubject) Close() {
	subject.LastValues = nil
	subject.Subscriptions = make(map[Subscription]interface{})
}

// Next takes an undefined amount of parameters
// which will be passed to subscribed functions
func (subject *BufferSubject) Next(values ...interface{}) {
	data := append([][]interface{}{values}, subject.LastValues...)
	subject.LastValues = append([][]interface{}(nil), data[:len(subject.LastValues)]...)
	subject.Subject.Next(values...)
}

// Pipe decorates an observable with one or multiple middlewares
// and returns a new observable with the decoration applied
func (subject *BufferSubject) Pipe(fns ...func(Observable, Subjectable)) Observable {
	parent := subject
	for _, fn := range fns {
		sub := NewBufferSubject(len(subject.LastValues))
		fn(parent, sub)
		parent = sub
	}
	return parent
}

// Subscribe registers a function for further updates of
// this observable and returns a subscription token which can
// be used to unsubscribe from it at any time
func (subject *BufferSubject) Subscribe(fn interface{}) (Subscription, error) {
	subscription, err := subject.Subject.Subscribe(fn)

	if err == nil && subject.LastValues != nil {
		data := subject.LastValues
		internal.Reverse(data)
		for _, valueArray := range data {
			if valueArray != nil {
				subject.notifySubscriber(subscription, valueArray)
			}
		}
	}

	return subscription, err
}

// NewBufferSubject returns a pointer
// to an empty instance of BufferSubject
func NewBufferSubject(bufferSize int) *BufferSubject {
	return &BufferSubject{
		LastValues: make([][]interface{}, bufferSize),
		Subject: Subject{
			Subscriptions: make(map[Subscription]interface{}),
		},
	}
}
