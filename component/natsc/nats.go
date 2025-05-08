package natsc

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/VanThen60hz/service-context/component/logger"
	pb "github.com/VanThen60hz/service-context/component/pubsub"
	"github.com/nats-io/nats.go"
)

type NatsOpt struct {
	prefix   string
	server   string
	username string
	password string
	token    string
}

type NatsComponent interface {
	InitFlags()
	Configure() error
	Run() error
	Stop() <-chan bool
	Get() interface{}
	Name() string
	GetPrefix() string
	Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error
	Subscribe(ctx context.Context, channel pb.Channel, eventTitle string) (c <-chan *pb.Event, close func())
}

type natsComponent struct {
	name      string
	logger    logger.Logger
	nc        *nats.Conn
	isRunning bool
	*NatsOpt
}

func NewNatsComponent(name string, prefix string) NatsComponent {
	return &natsComponent{
		name: name,
		NatsOpt: &NatsOpt{
			prefix: prefix,
		},
		isRunning: false,
	}
}

func (n *natsComponent) GetPrefix() string {
	if n.prefix == "" {
		return n.name
	}
	return n.prefix
}

func (n *natsComponent) Get() interface{} {
	return n
}

func (n *natsComponent) Name() string {
	return n.name
}

func (n *natsComponent) InitFlags() {
	prefix := n.prefix
	if n.prefix != "" {
		prefix += "-"
	}

	flag.StringVar(&n.server, prefix+"nats-server", "", "Nats connect server. Ex: \"nats://..., nats://\"")
	flag.StringVar(&n.username, prefix+"nats-username", "", "Nats username")
	flag.StringVar(&n.password, prefix+"nats-password", "", "Nats password")
	flag.StringVar(&n.token, prefix+"nats-token", "", "Nats token")
}

func (n *natsComponent) Configure() error {
	if n.isRunning {
		return nil
	}
	n.logger = logger.GetCurrent().GetLogger(n.name)
	n.logger.Info("Connecting to Nats at ", n.server, " ...")

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

	if n.username != "" {
		options = append(options, nats.UserInfo(n.username, n.password))
	}
	if n.token != "" {
		options = append(options, nats.Token(n.token))
	}

	nc, err := nats.Connect(n.server, options...)
	if err != nil {
		return err
	}

	n.nc = nc
	n.isRunning = true
	return nil
}

func (n *natsComponent) Run() error {
	return n.Configure()
}

func (n *natsComponent) Stop() <-chan bool {
	if n.nc != nil {
		err := n.nc.Drain()
		if err != nil {
			n.logger.Errorf("Error when drain nats connection: %q\n", err)
		}
	}
	n.isRunning = false

	c := make(chan bool)
	go func() { c <- true }()
	return c
}

func (n *natsComponent) Publish(ctx context.Context, channel pb.Channel, data *pb.Event) error {
	dataByte, err := json.Marshal(data.Data)
	if err != nil {
		n.logger.Errorln(err)
		return err
	}

	if err := n.nc.Publish(string(channel), dataByte); err != nil {
		n.logger.Errorln(err)
		return err
	}

	//if err := n.nc.Flush(); err != nil {
	//	n.logger.Errorln(err)
	//	return err
	//}

	return nil
}

func (n *natsComponent) Subscribe(ctx context.Context, channel pb.Channel, eventTitle string) (c <-chan *pb.Event, cl func()) {
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
