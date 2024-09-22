package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type IConsumer interface {
	Consume(ctx context.Context, streamName string, groupName string, consumerName string, processFunc func(message map[string]interface{}))
}
type consumer struct {
	redisCli *redis.Client
}

func NewConsumer(redisCli *redis.Client) IConsumer {
	return &consumer{
		redisCli: redisCli,
	}
}

func (c *consumer) Consume(ctx context.Context, streamName string, groupName string, consumerName string, processFunc func(message map[string]interface{})) {
	// log.Println("Starting consumer...")

	// Create a consumer group (only once)
	err := c.redisCli.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Could not create group: %v", err)
	}

	// Read existing messages
	// TODO: if startID is "0", and there are error, retry it
	c.readMessages(ctx, streamName, groupName, consumerName, processFunc, "0")

	// Continue to read new messages
	for {
		c.readMessages(ctx, streamName, groupName, consumerName, processFunc, ">")

		time.Sleep(1 * time.Second)
	}
}

func (c *consumer) readMessages(ctx context.Context, streamName string, groupName string, consumerName string, processFunc func(message map[string]interface{}), startID string) {
	log.Printf("reading message with startID: %s\n", startID)

	messages, err := c.redisCli.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamName, startID},
		Count:    10,
		// Blocking call, if its 0, it will block indefinitely
		Block: 0,
	}).Result()
	// if startID == ">" && strings.Contains(err.Error(), "i/o timeout") {
	// 	log.Printf("no new message within block time\n")
	// 	return
	// }
	if err != nil {
		log.Fatalf("Could not read messages: %v", err)
	}

	log.Printf("Received %d messages\n", len(messages))

	for _, msg := range messages {
		for _, m := range msg.Messages {
			fmt.Printf("Received message: %+v\n", m.Values)

			// Process the message using the provided callback
			processFunc(m.Values)

			// Acknowledge the message
			err := c.redisCli.XAck(ctx, streamName, groupName, m.ID).Err()
			if err != nil {
				log.Fatalf("Could not acknowledge message: %v", err)
			}
		}
	}
}
