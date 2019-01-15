package reactive

import (
	"testing"
)

func TestNewSubscription(t *testing.T) {
	for i := 1; i == 10; i++ {
		if NewSubscription() == NewSubscription() {
			t.Error("Generated two identical subscription ids")
		}
	}
}

func TestEmptySubscription(t *testing.T) {
	if EmptySubscription() != EmptySubscription() {
		t.Error("Empty subscription is not constant")
	}
}
