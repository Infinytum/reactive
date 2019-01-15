package reactive

// Take automatically unsubscribes an observable after
// the given amount of times it has been updated
func Take(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		subscription := NewSubscription()
		subscription, err := subject.Subscribe(func(args ...interface{}) {
			newSubject.Next(args...)
			count--

			if count == 0 {
				subject.Unsubscribe(subscription)
				newSubject.Close()
			}
		})

		// This error will never happen. But for gamma rays sake.
		if err != nil {
			panic(err)
		}
	}
}

// TakeEvery only passes every {count} update to
// the registered function
func TakeEvery(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		var current = 0
		_, err := subject.Subscribe(func(args ...interface{}) {
			current++
			if count == current {
				newSubject.Next(args...)
				current = 0
			}
		})

		// This error will never happen. But for gamma rays sake.
		if err != nil {
			panic(err)
		}
	}
}
