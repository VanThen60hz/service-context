package natsc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/VanThen60hz/service-context/component/pubsub"
	"github.com/nats-io/nats.go"
)

func (n *NatsComponent) Subscribe(ctx context.Context, channel pb.Channel, eventTitle string) (c <-chan *pb.Event, cl func()) {
	ch := make(chan *pb.Event)

	sub, err := n.nc.Subscribe(string(channel), func(msg *nats.Msg) {
		var data interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			n.logger.Errorln("Error unmarshaling message:", err)
			return
		}

		evt := &pb.Event{
			Id:         fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:      eventTitle,
			Channel:    channel,
			RemoteData: msg.Data,
			Data:       data,
			CreatedAt:  time.Now().UTC(),
		}
		ch <- evt
	})
	if err != nil {
		n.logger.Errorln(err)
	}

	return ch, func() {
		_ = sub.Unsubscribe()
		close(ch)
	}
}

func (n *NatsComponent) SubscribeWithHandler(ctx context.Context, channel pb.Channel, handler func(*pb.Event)) (close func(), err error) {
	sub, err := n.nc.Subscribe(string(channel), func(msg *nats.Msg) {
		var data interface{}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			n.logger.Errorln("Error unmarshaling message:", err)
			return
		}

		evt := &pb.Event{
			Id:         fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:      string(channel),
			Channel:    channel,
			RemoteData: msg.Data,
			Data:       data,
			CreatedAt:  time.Now().UTC(),
		}
		handler(evt)
	})
	if err != nil {
		n.logger.Errorln(err)
		return nil, err
	}

	return func() {
		_ = sub.Unsubscribe()
	}, nil
}
