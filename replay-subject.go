package reactive

import (
	"reflect"
)

// ReplaySubject is a special implementation of a subject
// It will always keep the last submitted value and new subscribers
// will receive that value immediately.
type ReplaySubject struct {
	LastValues    []interface{}
	Subscriptions map[Subscription]interface{}
}

// AsChannel returns a channel which will receive all
// further updates of this observable
func (subject *ReplaySubject) AsChannel() <-chan []interface{} {
	channel := make(chan []interface{})
	if subject.LastValues != nil {
		go func(channel chan []interface{}) {
			channel <- subject.LastValues
		}(channel)
	}
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
func (subject *ReplaySubject) Subscribe(fn interface{}) Subscription {
	subscription := NewSubscription()
	subject.Subscriptions[subscription] = fn

	if subject.LastValues != nil {
		subject.notifySubscriber(subscription, subject.LastValues)
	}

	return subscription
}

// Unsubscribe unregisters a previously registered function for all
// further updates of this observable or until re-registering.
func (subject *ReplaySubject) Unsubscribe(subscription Subscription) {
	if _, ok := subject.Subscriptions[subscription]; ok {
		delete(subject.Subscriptions, subscription)
	}
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

// Next takes an undefined amount of parameters
// which will be passed to subscribed functions
func (subject *ReplaySubject) Next(values ...interface{}) {
	subject.LastValues = values
	for subscription := range subject.Subscriptions {
		subject.notifySubscriber(subscription, values)
	}
}

func (subject ReplaySubject) notifySubscriber(subscription Subscription, values []interface{}) {
	passedArguments := make([]reflect.Value, 0)
	for _, arg := range values {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}

	if fn, ok := subject.Subscriptions[subscription]; ok {
		reflect.ValueOf(fn).Call(passedArguments)
	}
}

// NewReplaySubject returns a pointer
// to an empty instance of ReplaySubject
func NewReplaySubject() *ReplaySubject {
	return &ReplaySubject{
		Subscriptions: make(map[Subscription]interface{}),
	}
}
