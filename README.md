# Sample Stream Application

This is a sample application that demonstrates the use of Redis streams for message publishing and consumption with consumer group.

## Overview

The application consists of two main components:

- A publisher that sends messages to a Redis stream
- A consumer that reads messages from the Redis stream and processes them

## Prerequisites

spin up container through docker-compose

## Configuration

The application uses the following variables:

- `REDIS_ADDRESS`: the address of the Redis instance (default: `localhost:6379`)
- `REDIS_STREAM_NAME`: the name of the Redis stream (default: `sample-stream`)

## Running the Application

To run the application, use the following command:

```bash
Make run-local
```

## TODO

- Add graceful handling for timeouts in the consumer when using XREADGROUP with the ">" ID. Implement retry logic and exponential backoff to avoid continuous retries when no new messages are available.
- Add configuration options to adjust the block timeout dynamically based on system load or message frequency

### References

- [set redis xreadgroup block time for a year ?](https://stackoverflow.com/questions/63673564/redis-xreadgroup-stream-with-block-for-a-year-stupid)
