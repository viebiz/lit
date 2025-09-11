# Kafka Messaging Helpers

The `broker/kafka` package provides helpers for producing and consuming Kafka messages with consistent configuration, retry and monitoring support.

## Producer Helpers

Two producers are available:

### SyncProducer

`SyncProducer` publishes messages and waits for the broker to acknowledge delivery. It is created with `NewSyncProducer` and exposes `SendMessage` for sending payloads【F:broker/kafka/producer_sync.go†L21-L62】【F:broker/kafka/producer_sync.go†L65-L94】

### AsyncProducer

`AsyncProducer` enqueues messages without blocking. Messages are placed on an internal queue by `SendMessage` and published by calling `Listen`, which handles successes and errors asynchronously【F:broker/kafka/producer_async.go†L13-L68】【F:broker/kafka/producer_async.go†L71-L155】

### Configuration Options

Producer behaviour can be tuned with options:

- `ProducerWithAckMode` sets acknowledgement level for writes【F:broker/kafka/producer_option.go†L13-L18】
- `ProducerWithAutoCreateTopics` allows automatic topic creation【F:broker/kafka/producer_option.go†L21-L26】
- `ProducerWithRoundRobinPartitioner` forces round-robin partitioning【F:broker/kafka/producer_option.go†L28-L33】
- `ProducerWithFlushFrequency` configures batching frequency【F:broker/kafka/producer_option.go†L35-L41】
- `ProducerWithCompression` chooses a compression codec【F:broker/kafka/producer_option.go†L43-L48】
- `ProducerWithTLS` enables TLS connections【F:broker/kafka/producer_option.go†L50-L55】

## Consumer Helpers

`NewConsumerGroup` builds a consumer group and wires the handler, retry settings and monitoring【F:broker/kafka/consumer.go†L24-L80】. Consumption is started with `Consume`, which blocks until the provided context is cancelled【F:broker/kafka/consumer.go†L83-L112】.

### Configuration Options

Consumers support several options:

- `ConsumerWithTLS` uses TLS for broker connections【F:broker/kafka/consumer_option.go†L12-L17】
- `ConsumerWithOffsetNewest` starts from the latest offset【F:broker/kafka/consumer_option.go†L19-L23】
- `ConsumerMaxRetryPerMessage` caps retry attempts per message【F:broker/kafka/consumer_option.go†L26-L31】
- `ConsumerWithAutoCreateTopics` enables automatic topic creation【F:broker/kafka/consumer_option.go†L33-L37】
- `ConsumerWithCustomConsumerGroupID` sets a custom group id【F:broker/kafka/consumer_option.go†L40-L45】
- `ConsumerDisablePayloadLogging` suppresses payload logging【F:broker/kafka/consumer_option.go†L48-L52】

## Retry Logic

Each consumed message is processed with exponential backoff. The handler is retried up to `maxRetriesPerMsg` and the offset is committed regardless of success, ensuring the consumer continues to the next message【F:broker/kafka/consumer_handler.go†L98-L115】【F:broker/kafka/consumer_handler.go†L118-L134】

## Monitoring Hooks

Producers and consumers emit telemetry through the `monitoring` and `instrumentkafka` packages. Handlers create spans for publish, consume and commit events, capturing partition and offset information for tracing【F:broker/kafka/producer_sync.go†L76-L86】【F:broker/kafka/producer_async.go†L77-L118】【F:broker/kafka/consumer_handler.go†L57-L74】【F:broker/kafka/consumer_handler.go†L124-L133】

## Sample Implementations

### Sync Producer

```go
ctx := context.Background()
producer, err := kafka.NewSyncProducer(ctx, kafka.Config{AppName: "my-app", Server: "local"}, []string{"localhost:9092"},
    kafka.ProducerWithAutoCreateTopics())
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

err = producer.SendMessage(ctx, "my-topic", []byte("hello"), kafka.ProducerMessageOption{
    Key: "msg-1",
})
if err != nil {
    log.Fatal(err)
}
```

### Consumer Group

```go
ctx := context.Background()
handler := func(ctx context.Context, msg kafka.ConsumerMessage) error {
    fmt.Printf("received: %s\n", msg.Value)
    return nil
}

consumer, err := kafka.NewConsumerGroup(ctx, kafka.Config{AppName: "my-app", Server: "local"}, []string{"my-topic"},
    []string{"localhost:9092"}, handler, kafka.ConsumerMaxRetryPerMessage(3))
if err != nil {
    log.Fatal(err)
}

if err := consumer.Consume(ctx); err != nil {
    log.Fatal(err)
}
```

