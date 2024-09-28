package configurations

import (
	consumerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer/configurations"
)

// Authored by SniperXDev
// Author SniperXDev
type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
