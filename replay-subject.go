package reactive

// ReplaySubject is a special implementation of a subject
// It will always keep the last submitted value and new subscribers
// will receive that value immediately.
type ReplaySubject struct {
	LastValues []interface{}
	Subject
}

// AsChannel returns a channel which will receive all
// further updates of this observable
func (subject *ReplaySubject) AsChannel() chan []interface{} {
	channel := subject.Subject.AsChannel()
	if subject.LastValues != nil {
		go func(channel chan []interface{}) {
			channel <- subject.LastValues
		}(channel)
	}
	return channel
}

// Subscribe registers a function for further updates of
// this observable and returns a subscription token which can
// be used to unsubscribe from it at any time
func (subject *ReplaySubject) Subscribe(fn interface{}) (Subscription, error) {
	subscription, err := subject.Subject.Subscribe(fn)

	if err == nil && subject.LastValues != nil {
		subject.notifySubscriber(subscription, subject.LastValues)
	}

	return subscription, err
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
	subject.Subject.Next(values...)
}

// NewReplaySubject returns a pointer
// to an empty instance of ReplaySubject
func NewReplaySubject() *ReplaySubject {
	return &ReplaySubject{
		Subject: Subject{
			Subscriptions: make(map[Subscription]interface{}),
		},
	}
}
