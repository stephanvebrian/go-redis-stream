package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type IPublisher interface {
	Publish(ctx context.Context, streamName string, data map[string]interface{}) *redis.StringCmd
}

type publisher struct {
	redisCli *redis.Client
}

func NewPublisher(rdsCli *redis.Client) IPublisher {
	return &publisher{
		redisCli: rdsCli,
	}
}

func (p *publisher) Publish(ctx context.Context, streamName string, data map[string]interface{}) *redis.StringCmd {
	cmd := p.redisCli.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: data,
	})

	return cmd
}
