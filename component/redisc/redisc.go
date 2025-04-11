package redisc

import (
	"flag"

	sctx "github.com/VanThen60hz/service-context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisComponent struct {
	id     string
	logger sctx.Logger
	cfg    redisConfig
	client *redis.Client
}

type redisConfig struct {
	address  string
	password string
	db       int
}

func NewRedisComponent(id string) *RedisComponent {
	return &RedisComponent{id: id}
}

func (r *RedisComponent) ID() string {
	return r.id
}

func (r *RedisComponent) InitFlags() {
	flag.StringVar(&r.cfg.address, "redis-address", "localhost:6379", "Redis server address")
	flag.StringVar(&r.cfg.password, "redis-password", "", "Redis password")
	flag.IntVar(&r.cfg.db, "redis-db", 0, "Redis database number")
}

func (r *RedisComponent) Activate(ctx sctx.ServiceContext) error {
	r.logger = ctx.Logger(r.id)

	if err := r.cfg.validate(); err != nil {
		return err
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     r.cfg.address,
		Password: r.cfg.password,
		DB:       r.cfg.db,
	})

	return nil
}

func (r *RedisComponent) Stop() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

func (cfg *redisConfig) validate() error {
	if cfg.address == "" {
		return errors.New("Redis address is missing")
	}
	return nil
}
