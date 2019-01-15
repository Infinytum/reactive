package reactive

import (
	"testing"
	"time"
)

func TestNewSubject(t *testing.T) {
	subject := NewSubject().(*subject)

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
	subject := NewSubject().(*subject)
	subscription, _ := subject.Subscribe(func() {})

	if _, exists := subject.Subscriptions[subscription]; !exists {
		t.Error("Subscription is not in subscription map")
	}

	if sub, err := subject.Subscribe(3); err == nil || sub != EmptySubscription() {
		t.Error("Subscribe accepted non-function parameters")
	}
}

func TestSubject_Next(t *testing.T) {
	subject := NewSubject().(*subject)
	didRun := false
	subject.Subscribe(func(run bool) {
		didRun = run
	})

	subject.Next(true)

	if !didRun {
		t.Error("Subscription handler wasnt called!")
	}
}

func TestSubject_Pipe(t *testing.T) {
	subject := NewSubject().(*subject)

	if subject != subject.Pipe() {
		t.Error("Empty pipe is different from original")
	}

	if subject == subject.Pipe(Take(1)) {
		t.Error("Take pipe resulted in original subject")
	}
}

func TestSubject_Unsubscribe(t *testing.T) {
	subject := NewSubject().(*subject)
	subscription, _ := subject.Subscribe(func() {})

	if _, exists := subject.Subscriptions[subscription]; !exists {
		t.Error("Subscription is not in subscription map")
	}

	subject.Unsubscribe(subscription)

	if _, exists := subject.Subscriptions[subscription]; exists {
		t.Error("Subscription is still in subscription map")
	}
}

func TestSubject_AsChannel(t *testing.T) {
	subject := NewSubject().(*subject)
	channel := make(chan interface{})

	go func() {
		data := <-subject.AsChannel()
		channel <- data
	}()

	// Wait for handler to register
	<-time.After(time.Duration(400) * time.Millisecond)
	subject.Next(true)

	select {
	case <-channel:
		return
	case <-time.After(time.Duration(500) * time.Millisecond):
		t.Error("Subscription handler wasnt called!")
		break
	}
}
