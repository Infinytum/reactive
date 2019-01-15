package reactive

import "github.com/infinytum/reactive/internal"

// bufferSubject is a special implementation of a subject
// It will always keep the submitted values until the buffer is full
// and new subscribers will receive those values immediately.
//
// New values will remove the oldest value from the buffer
type bufferSubject struct {
	LastValues [][]interface{}
	subject
}

// AsChannel returns a channel which will receive all
// further updates of this observable
func (subject *bufferSubject) AsChannel() chan []interface{} {
	channel := subject.subject.AsChannel()
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
func (subject *bufferSubject) Close() {
	subject.LastValues = nil
	subject.Subscriptions = make(map[Subscription]interface{})
}

// Next takes an undefined amount of parameters
// which will be passed to subscribed functions
func (subject *bufferSubject) Next(values ...interface{}) {
	data := append([][]interface{}{values}, subject.LastValues...)
	subject.LastValues = append([][]interface{}(nil), data[:len(subject.LastValues)]...)
	subject.subject.Next(values...)
}

// Pipe decorates an observable with one or multiple middlewares
// and returns a new observable with the decoration applied
func (subject *bufferSubject) Pipe(fns ...func(Observable, Subjectable)) Observable {
	parent := subject
	for _, fn := range fns {
		sub := NewBufferSubject(len(subject.LastValues)).(*bufferSubject)
		fn(parent, sub)
		parent = sub
	}
	return parent
}

// Subscribe registers a function for further updates of
// this observable and returns a subscription token which can
// be used to unsubscribe from it at any time
func (subject *bufferSubject) Subscribe(fn interface{}) (Subscription, error) {
	subscription, err := subject.subject.Subscribe(fn)

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
// to an empty instance of bufferSubject
func NewBufferSubject(bufferSize int) Subjectable {
	return &bufferSubject{
		LastValues: make([][]interface{}, bufferSize),
		subject: subject{
			Subscriptions: make(map[Subscription]interface{}),
		},
	}
}
