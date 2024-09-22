package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:         RedisAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	// Start the consumer in a goroutine
	consumer := NewConsumer(rdb)
	go func() {
		consumer.Consume(
			ctx,
			RedisStreamName,
			"hercules",
			"consumer-1",
			func(message map[string]interface{}) {
				log.Printf("processing message:\n\tcontent:%s\n\tidentifier:%s\n", message["content"], message["identifier"])

				time.Sleep(1 * time.Second)
			},
		)
	}()

	// Give the consumer a moment to start
	time.Sleep(1 * time.Second)

	// publisher := NewPublisher(rdb)
	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		message := map[string]interface{}{
	// 			"content":    faker.Word(),
	// 			"identifier": faker.UUIDHyphenated(),
	// 		}

	// 		cmd := publisher.Publish(
	// 			ctx,
	// 			RedisStreamName,
	// 			message,
	// 		)

	// 		err := cmd.Err()
	// 		if err != nil {
	// 			log.Printf("Could not publish message: %v", err)
	// 		}

	// 		if err == nil {
	// 			log.Printf("Published message successfully identifier:%s\n", message["identifier"])
	// 		}

	// 		time.Sleep(50 * time.Millisecond) // Sleep for a bit to stagger messages
	// 	}
	// }()

	// Prevent main from exiting immediately
	select {}
}
