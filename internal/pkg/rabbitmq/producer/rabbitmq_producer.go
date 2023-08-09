//go:build go1.18

package producer

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/rabbitmq/amqp091-go"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	messageHeader "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/message_header"
	producer2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/otel/tracing/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/producer"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

type rabbitMQProducer struct {
	logger                  logger.Logger
	connection              types.IConnection
	eventSerializer         serializer.EventSerializer
	producersConfigurations map[string]*configurations.RabbitMQProducerConfiguration
	isProducedNotifications []func(message types2.IMessage)
}

func NewRabbitMQProducer(
	connection types.IConnection,
	rabbitmqProducersConfiguration map[string]*configurations.RabbitMQProducerConfiguration,
	logger logger.Logger,
	eventSerializer serializer.EventSerializer,
	isProducedNotifications ...func(message types2.IMessage),
) (producer.Producer, error) {
	p := &rabbitMQProducer{
		logger:                  logger,
		connection:              connection,
		eventSerializer:         eventSerializer,
		producersConfigurations: rabbitmqProducersConfiguration,
	}

	p.isProducedNotifications = isProducedNotifications

	return p, nil
}

func (r *rabbitMQProducer) IsProduced(h func(message types2.IMessage)) {
	r.isProducedNotifications = append(r.isProducedNotifications, h)
}

func (r *rabbitMQProducer) PublishMessage(
	ctx context.Context,
	message types2.IMessage,
	meta metadata.Metadata,
) error {
	return r.PublishMessageWithTopicName(ctx, message, meta, "")
}

func (r *rabbitMQProducer) getProducerConfigurationByMessage(
	message types2.IMessage,
) *configurations.RabbitMQProducerConfiguration {
	messageType := utils.GetMessageBaseReflectType(message)
	return r.producersConfigurations[messageType.String()]
}

func (r *rabbitMQProducer) PublishMessageWithTopicName(
	ctx context.Context,
	message types2.IMessage,
	meta metadata.Metadata,
	topicOrExchangeName string,
) error {
	if message.IsMessage() == false {
		return errors.New(
			fmt.Sprintf(
				"message %s is not a message type or message property is nil",
				utils.GetMessageBaseReflectType(message),
			),
		)
	}

	producerConfiguration := r.getProducerConfigurationByMessage(message)

	if producerConfiguration == nil {
		producerConfiguration = configurations.NewDefaultRabbitMQProducerConfiguration(message)
	}

	var exchange string
	var routingKey string

	if topicOrExchangeName != "" {
		exchange = topicOrExchangeName
	} else if producerConfiguration != nil && producerConfiguration.ExchangeOptions.Name != "" {
		exchange = producerConfiguration.ExchangeOptions.Name
	} else {
		exchange = utils.GetTopicOrExchangeName(message)
	}

	if producerConfiguration != nil && producerConfiguration.RoutingKey != "" {
		routingKey = producerConfiguration.RoutingKey
	} else {
		routingKey = utils.GetRoutingKey(message)
	}

	meta = r.getMetadata(message, meta)

	producerOptions := &producer2.ProducerTracingOptions{
		MessagingSystem: "rabbitmq",
		DestinationKind: "exchange",
		Destination:     exchange,
		OtherAttributes: []attribute.KeyValue{
			semconv.MessagingRabbitmqRoutingKeyKey.String(routingKey),
		},
	}

	serializedObj, err := r.eventSerializer.Serialize(message)
	if err != nil {
		return err
	}

	ctx, beforeProduceSpan := producer2.StartProducerSpan(
		ctx,
		message,
		&meta,
		string(serializedObj.Data),
		producerOptions,
	)

	// https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/publisher_confirms.go
	if r.connection == nil {
		return producer2.FinishProducerSpan(beforeProduceSpan, errors.New("connection is nil"))
	}

	if r.connection.IsClosed() {
		return producer2.FinishProducerSpan(
			beforeProduceSpan,
			errors.New("connection is closed, wait for connection alive"),
		)
	}

	// create a unique channel on the connection and in the end close the channel
	channel, err := r.connection.Channel()
	if err != nil {
		return producer2.FinishProducerSpan(beforeProduceSpan, err)
	}
	defer channel.Close()

	err = r.ensureExchange(producerConfiguration, channel, exchange)
	if err != nil {
		return producer2.FinishProducerSpan(beforeProduceSpan, err)
	}

	if err := channel.Confirm(false); err != nil {
		return producer2.FinishProducerSpan(beforeProduceSpan, err)
	}

	confirms := make(chan amqp091.Confirmation)
	channel.NotifyPublish(confirms)

	props := amqp091.Publishing{
		CorrelationId:   messageHeader.GetCorrelationId(meta),
		MessageId:       message.GeMessageId(),
		Timestamp:       time.Now(),
		Headers:         metadata.MetadataToMap(meta),
		Type:            message.GetEventTypeName(), // typeMapper.GetTypeName(message) - just message type name not full type name because in other side package name for type could be different
		ContentType:     serializedObj.ContentType,
		Body:            serializedObj.Data,
		DeliveryMode:    producerConfiguration.DeliveryMode,
		Expiration:      producerConfiguration.Expiration,
		AppId:           producerConfiguration.AppId,
		Priority:        producerConfiguration.Priority,
		ReplyTo:         producerConfiguration.ReplyTo,
		ContentEncoding: producerConfiguration.ContentEncoding,
	}

	err = channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		true,
		false,
		props,
	)
	if err != nil {
		return producer2.FinishProducerSpan(beforeProduceSpan, err)
	}

	if confirmed := <-confirms; !confirmed.Ack {
		return producer2.FinishProducerSpan(beforeProduceSpan, errors.New("ack not confirmed"))
	}

	if len(r.isProducedNotifications) > 0 {
		for _, notification := range r.isProducedNotifications {
			if notification != nil {
				notification(message)
			}
		}
	}

	return producer2.FinishProducerSpan(beforeProduceSpan, err)
}

func (r *rabbitMQProducer) getMetadata(
	message types2.IMessage,
	meta metadata.Metadata,
) metadata.Metadata {
	meta = metadata.FromMetadata(meta)

	if message.GetEventTypeName() == "" {
		message.SetEventTypeName(
			typeMapper.GetTypeName(message),
		) // just message type name not full type name because in other side package name for type could be different)
	}
	messageHeader.SetMessageType(meta, message.GetEventTypeName())
	messageHeader.SetMessageContentType(meta, r.eventSerializer.ContentType())

	if messageHeader.GetMessageId(meta) == "" {
		messageHeader.SetMessageId(meta, message.GeMessageId())
	}

	if messageHeader.GetMessageCreated(meta) == *new(time.Time) {
		messageHeader.SetMessageCreated(meta, message.GetCreated())
	}

	if messageHeader.GetCorrelationId(meta) == "" {
		cid := uuid.NewV4().String()
		messageHeader.SetCorrelationId(meta, cid)
	}
	messageHeader.SetMessageName(meta, utils.GetMessageName(message))

	return meta
}

func (r *rabbitMQProducer) ensureExchange(
	producersConfigurations *configurations.RabbitMQProducerConfiguration,
	channel *amqp091.Channel,
	exchangeName string,
) error {
	err := channel.ExchangeDeclare(
		exchangeName,
		string(producersConfigurations.ExchangeOptions.Type),
		producersConfigurations.ExchangeOptions.Durable,
		producersConfigurations.ExchangeOptions.AutoDelete,
		false,
		false,
		producersConfigurations.ExchangeOptions.Args,
	)
	if err != nil {
		return err
	}

	return nil
}
