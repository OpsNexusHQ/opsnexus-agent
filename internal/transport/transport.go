package transport

import "context"

// Transport defines the interface for transmitting data/metrics to a destination.
type Transport interface {
	Send(ctx context.Context, payload any) error
}
