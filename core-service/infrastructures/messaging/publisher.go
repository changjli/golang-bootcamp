package messaging

import (
	"context"
	"events"
)

type PaymentPublisher interface {
	PublishPaymentRequested(ctx context.Context, event *events.PaymentRequestedEvent) error
}
