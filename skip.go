package reactive

// Skip will ignore a specified amount of updates
// and will pass through all following
func Skip(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		_, err := subject.Subscribe(func(args ...interface{}) {
			if count == 0 {
				newSubject.Next(args...)
			} else {
				count--
			}
		})

		// This error will never happen. But for gamma rays sake.
		if err != nil {
			panic(err)
		}
	}
}

// SkipEvery will skip every {count} update and will pass all others
func SkipEvery(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		var current = 0
		_, err := subject.Subscribe(func(args ...interface{}) {
			current++
			if count != current {
				newSubject.Next(args...)
			} else {
				current = 0
			}
		})

		// This error will never happen. But for gamma rays sake.
		if err != nil {
			panic(err)
		}
	}
}
