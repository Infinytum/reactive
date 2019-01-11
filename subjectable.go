package reactive

// Subjectable defines required methods
// for an object to be considered a subject
type Subjectable interface {
	Observable

	// Next takes an undefined amount of parameters
	// which will be passed to subscribed functions
	Next(values ...interface{})
}
