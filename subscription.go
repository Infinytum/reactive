package reactive

import "github.com/google/uuid"

// Subscription represents a string to identify
// a subscription in an obserable so it can be removed
type Subscription uuid.UUID

// NewSubscription generates a new subscription
func NewSubscription() Subscription {
	return Subscription(uuid.New())
}

func EmptySubscription() Subscription {
	emptyBytes := make([]byte, 16)
	emptyUuid, _ := uuid.ParseBytes(emptyBytes)
	return Subscription(emptyUuid)
}
