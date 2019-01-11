package reactive

import (
	"testing"
)

func TestNewReplaySubject(t *testing.T) {
	subject := NewReplaySubject()

	if subject == nil {
		t.Error("NewReplaySubject returned nil")
	}

	if subject.Subscriptions == nil {
		t.Error("NewReplaySubject did not create empty subscriptions map")
	}

	if len(subject.Subscriptions) != 0 {
		t.Error("NewReplaySubject did not create empty subscriptions map")
	}
}

func TestReplaySubject_Next(t *testing.T) {
	subject := NewReplaySubject()
	subject.Next(1)

	if subject.LastValues[0] != 1 {
		t.Error("LastValues were not set")
	}
}

func TestReplaySubject_Pipe(t *testing.T) {
	subject := NewReplaySubject()

	if subject != subject.Pipe() {
		t.Error("Empty pipe is different from original")
	}

	if subject == subject.Pipe(Take(1)) {
		t.Error("Take pipe resulted in original subject")
	}
}

func TestReplaySubject_Subscribe(t *testing.T) {
	subject := NewReplaySubject()
	subscription := subject.Subscribe(func(nr int) {})

	if _, exists := subject.Subscriptions[subscription]; !exists {
		t.Error("Subscription is not in subscription map")
	}

	subject.Next(1)

	didRun := false
	subject.Subscribe(func(nr int) {
		didRun = true
	})

	if !didRun {
		t.Error("LastValue was not given to new subscriber")
	}
}

func TestReplaySubject_AsChannel(t *testing.T) {
	subject := NewReplaySubject()
	subject.Next(1)

	data := <-subject.AsChannel()
	if data[0] != 1 {
		t.Error("AsChannel did not return value of Next call")
	}
}
