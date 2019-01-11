package reactive

import "testing"

func TestNewSubscription(t *testing.T) {
	if NewSubscription() == "" {
		t.Error("NewSubscription didnt return a new subscription id")
	}
}
