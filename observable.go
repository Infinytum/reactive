package reactive

// Observable defines the requirements for a class
// to be considered a valid observable
type Observable interface {
	// AsChannel returns a channel which will receive all
	// further updates of this observable
	AsChannel() chan []interface{}

	// Pipe decorates an observable with one or multiple middlewares
	// and returns a new observable with the decoration applied
	Pipe(fns ...func(Observable, Subjectable)) Observable

	// Subscribe registers a function for further updates of
	// this observable and returns a subscription token which can
	// be used to unsubscribe from it at any time
	Subscribe(fn interface{}) (Subscription, error)

	// Unsubscribe unregisters a previously registered function for all
	// further updates of this observable or until re-registering.
	Unsubscribe(subscription Subscription)
}
