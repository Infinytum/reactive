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

// Close will remove all subscribers and render the subjectable useless
func (subject *Subject) Close() {
	subject.Subscriptions = make(map[Subscription]interface{})
}

// Next takes an undefined amount of parameters which will be passed to
// subscribed functions
func (subject *Subject) Next(values ...interface{}) {
	for subscription := range subject.Subscriptions {
		subject.notifySubscriber(subscription, values)
	}
}

func (subject Subject) notifySubscriber(subscription Subscription, values []interface{}) {
	if fn, ok := subject.Subscriptions[subscription]; ok {
		refFn := reflect.TypeOf(fn)
		fnArgs := make([]reflect.Value, 0, refFn.NumIn())

		for argIndex := 0; argIndex < refFn.NumIn(); argIndex++ {
			providedVal := values[argIndex]

			// Variadic arguments need special treatment
			if refFn.IsVariadic() {
				sliceType := refFn.In(argIndex).Elem()

				for _, innerVal := range values[argIndex:len(values)] {
					if providedVal == nil {
						fnArgs = append(fnArgs, reflect.New(sliceType).Elem())
						continue
					}

					if !reflect.TypeOf(innerVal).AssignableTo(sliceType) {
						// Slice does not match received data, skipping this subscriber
						return
					}
					fnArgs = append(fnArgs, reflect.ValueOf(innerVal))
				}
				// Finish loop as we have filled in all data to the slice
				break
			} else {
				argType := refFn.In(argIndex)
				if providedVal == nil {
					values[argIndex] = reflect.New(argType).Elem()
					providedVal = values[argIndex]
				}

				if !reflect.TypeOf(providedVal).AssignableTo(argType) {
					// Method signature not compatible with this input. Skipping subscriber
					return
				}

				fnArgs = append(fnArgs, reflect.ValueOf(values[argIndex]))

				if argIndex == refFn.NumIn()-1 {
					if refFn.NumIn() != len(fnArgs) {
						// Skipping non-slice overflow
						return
					}
				}
			}

		}

		reflect.ValueOf(fn).Call(fnArgs)
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

// NewSubject returns a pointer
// to an empty instance of Subject
func NewSubject() *Subject {
	return &Subject{
		Subscriptions: make(map[Subscription]interface{}),
	}
}
