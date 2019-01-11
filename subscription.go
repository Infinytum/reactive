package reactive

import "github.com/google/uuid"

// Subscription represents a string to identify
// a subscription in an obserable so it can be removed
type Subscription string

// NewSubscription generates a new subscription
func NewSubscription() Subscription {
	return Subscription(uuid.New().String())
}
