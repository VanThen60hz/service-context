package natsc

import (
	"context"
	"encoding/json"

	pb "github.com/VanThen60hz/service-context/component/pubsub"
	"github.com/nats-io/nats.go"
)

func (n *NatsComponent) Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error {
	dataByte, err := json.Marshal(data.Data)
	if err != nil {
		n.logger.Errorln(err)
		return err
	}

	if err := n.nc.Publish(string(channel), dataByte); err != nil {
		n.logger.Errorln(err)
		return err
	}

	return nil
}

func (n *NatsComponent) PublishAsync(ctx context.Context, channel pb.Channel, data *pb.Event) error {
	dataByte, err := json.Marshal(data.Data)
	if err != nil {
		n.logger.Errorln(err)
		return err
	}

	// Use PublishMsg for async publishing
	msg := &nats.Msg{
		Subject: string(channel),
		Data:    dataByte,
	}

	if err := n.nc.PublishMsg(msg); err != nil {
		n.logger.Errorln(err)
		return err
	}

	return nil
}
