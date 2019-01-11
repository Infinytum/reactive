package reactive

// Skip will ignore a specified amount of updates
// and will pass through all following
func Skip(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		subject.Subscribe(func(args ...interface{}) {
			if count == 0 {
				newSubject.Next(args...)
			} else {
				count--
			}
		})
	}
}

// SkipEvery will skip every {count} update and will pass all others
func SkipEvery(count int) func(Observable, Subjectable) {
	return func(subject Observable, newSubject Subjectable) {
		var current = 0
		subject.Subscribe(func(args ...interface{}) {
			current++
			if count != current {
				newSubject.Next(args...)
			} else {
				current = 0
			}
		})
	}
}
