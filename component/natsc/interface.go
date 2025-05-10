package natsc

import (
	"context"

	pb "github.com/VanThen60hz/service-context/component/pubsub"
)

type Nats interface {
	// Publish publishes an event to a channel
	Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error
	// Subscribe subscribes to a channel and returns a channel to receive events
	Subscribe(ctx context.Context, channel pb.Channel, eventTitle string) (c <-chan *pb.Event, close func())
	// PublishAsync publishes an event to a channel asynchronously
	PublishAsync(ctx context.Context, channel pb.Channel, data *pb.Event) error
	// SubscribeWithHandler subscribes to a channel with a custom handler
	SubscribeWithHandler(ctx context.Context, channel pb.Channel, handler func(*pb.Event)) (close func(), err error)
}

var _ Nats = (*NatsComponent)(nil)
