package natsc

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	sctx "github.com/VanThen60hz/service-context"
	pb "github.com/VanThen60hz/service-context/component/pubsub"
	"github.com/nats-io/nats.go"
)

type NatsComponent struct {
	id        string
	logger    sctx.Logger
	nc        *nats.Conn
	isRunning bool
	cfg       natsConfig
}

type natsConfig struct {
	server   string
	username string
	password string
	token    string
}

func NewNatsComponent(id string) *NatsComponent {
	return &NatsComponent{
		id:        id,
		isRunning: false,
	}
}

func (n *NatsComponent) ID() string {
	return n.id
}

func (n *NatsComponent) InitFlags() {
	flag.StringVar(&n.cfg.server, "nats-server", "", "Nats connect server. Ex: \"nats://..., nats://\"")
	flag.StringVar(&n.cfg.username, "nats-username", "", "Nats username")
	flag.StringVar(&n.cfg.password, "nats-password", "", "Nats password")
	flag.StringVar(&n.cfg.token, "nats-token", "", "Nats token")
}

func (n *NatsComponent) Activate(ctx sctx.ServiceContext) error {
	if n.isRunning {
		return nil
	}
	n.logger = ctx.Logger(n.id)
	n.logger.Info("Connecting to Nats at ", n.cfg.server, " ...")

	var options []nats.Option

	options = append(options,
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			n.logger.Errorf("Got disconnected! Reason: %q\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			n.logger.Errorf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			n.logger.Errorf("Connection closed. Reason: %q\n", nc.LastError())
		}))

	if n.cfg.username != "" {
		options = append(options, nats.UserInfo(n.cfg.username, n.cfg.password))
	}
	if n.cfg.token != "" {
		options = append(options, nats.Token(n.cfg.token))
	}

	nc, err := nats.Connect(n.cfg.server, options...)
	if err != nil {
		return err
	}

	n.nc = nc
	n.isRunning = true
	return nil
}

func (n *NatsComponent) Stop() error {
	if n.nc != nil {
		err := n.nc.Drain()
		if err != nil {
			n.logger.Errorf("Error when drain nats connection: %q\n", err)
			return err
		}
	}
	n.isRunning = false
	return nil
}

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
