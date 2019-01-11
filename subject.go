package reactive

import (
	"errors"
	"reflect"
)

// Subject is the basic implementation of a subjectable
type Subject struct {
	Subscriptions map[Subscription]interface{}
}

// AsChannel returns a channel which will receive all
// further updates of this observable
func (subject *Subject) AsChannel() chan []interface{} {
	channel := make(chan []interface{})
	go subject.Subscribe(func(args ...interface{}) {
		go func(channel chan []interface{}) {
			channel <- args
		}(channel)
	})
	return channel
}

// Subscribe registers a function for further updates of
// this observable and returns a subscription token which can
// be used to unsubscribe from it at any time
func (subject *Subject) Subscribe(fn interface{}) (Subscription, error) {
	if fn != nil && reflect.TypeOf(fn).Kind() == reflect.Func {
		subscription := NewSubscription()
		subject.Subscriptions[subscription] = fn

		return subscription, nil
	}
	return Subscription(""), errors.New("fn is not a function")
}

// Unsubscribe unregisters a previously registered function for all
// further updates of this observable or until re-registering.
func (subject *Subject) Unsubscribe(subscription Subscription) {
	if _, ok := subject.Subscriptions[subscription]; ok {
		delete(subject.Subscriptions, subscription)
	}
}

// Pipe decorates an observable with one or multiple middlewares
// and returns a new observable with the decoration applied
func (subject *Subject) Pipe(fns ...func(Observable, Subjectable)) Observable {
	parent := subject
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		sub := NewSubject()
		fn(parent, sub)
		parent = sub
	}
	return parent
}

// Next takes an undefined amount of parameters
// which will be passed to subscribed functions
func (subject *Subject) Next(values ...interface{}) {
	for subscription := range subject.Subscriptions {
		subject.notifySubscriber(subscription, values)
	}
}

func (subject Subject) notifySubscriber(subscription Subscription, values []interface{}) {
	if fn, ok := subject.Subscriptions[subscription]; ok {
		refFn := reflect.ValueOf(fn)
		fnArgs := make([]reflect.Value, 0)
		if refFn.Type().NumIn() != len(values) {
			return
		}

		for i := 0; i < refFn.Type().NumIn(); i++ {
			val := values[i]

			if val == nil {
				fnArgs = append(fnArgs, reflect.New(refFn.Type().In(i)).Elem())
			} else {
				fnArgs = append(fnArgs, reflect.ValueOf(val))
			}

			if val != nil && refFn.Type().In(i).Kind() != reflect.ValueOf(values[i]).Type().Kind() {
				return
			}
		}

		refFn.Call(fnArgs)
	}
}

// NewSubject returns a pointer
// to an empty instance of Subject
func NewSubject() *Subject {
	return &Subject{
		Subscriptions: make(map[Subscription]interface{}),
	}
}
