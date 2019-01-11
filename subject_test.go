package reactive

import (
	"testing"
)

func TestNewSubject(t *testing.T) {
	subject := NewSubject()

	if subject == nil {
		t.Error("NewSubject returned nil")
	}

	if subject.Subscriptions == nil {
		t.Error("NewSubject did not create empty subscriptions map")
	}

	if len(subject.Subscriptions) != 0 {
		t.Error("NewSubject did not create empty subscriptions map")
	}
}

func TestSubject_Subscribe(t *testing.T) {
	subject := NewSubject()
	subscription := subject.Subscribe(func() {})

	if _, exists := subject.Subscriptions[subscription]; !exists {
		t.Error("Subscription is not in subscription map")
	}
}

func TestSubject_Next(t *testing.T) {
	subject := NewSubject()
	didRun := false
	subject.Subscribe(func(run bool) {
		didRun = run
	})

	subject.Next(true)

	if !didRun {
		t.Error("Subscription handler wasnt called!")
	}
}
