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
