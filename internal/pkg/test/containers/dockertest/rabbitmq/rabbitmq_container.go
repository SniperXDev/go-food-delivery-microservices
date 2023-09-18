package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	bus2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rabbitmq/amqp091-go"
)

type rabbitmqDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.RabbitMQContainerOptions
	pool           *dockertest.Pool
	logger         logger.Logger
}

func NewRabbitMQDockerTest(logger logger.Logger) contracts.RabbitMQContainer {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return &rabbitmqDockerTest{
		defaultOptions: &contracts.RabbitMQContainerOptions{
			Ports:       []string{"5672", "15672"},
			Host:        "localhost",
			VirtualHost: "",
			UserName:    "dockertest",
			Password:    "dockertest",
			Tag:         "management",
			ImageName:   "rabbitmq",
			Name:        "rabbitmq-dockertest",
		},
		logger: logger,
		pool:   pool,
	}
}

func (g *rabbitmqDockerTest) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.RabbitMQContainerOptions,
) (*config.RabbitmqHostOptions, error) {
	runOptions := g.getRunOptions(options...)

	// pull mongodb docker image for version 5.0
	resource, err := g.pool.RunWithOptions(runOptions, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource (RabbitMQ Container): %s", err)
	}

	resource.Expire(
		120,
	) // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application_exceptions in the container might not be ready to accept connections yet

	g.resource = resource
	hostPort, err := strconv.Atoi(
		resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Ports[0])),
	) // 5672
	httpPort, err := strconv.Atoi(
		resource.GetPort(fmt.Sprintf("%s/tcp", g.defaultOptions.Ports[1])),
	) // 15672

	g.defaultOptions.HostPort = hostPort
	g.defaultOptions.HttpPort = httpPort

	t.Cleanup(func() { _ = resource.Close() })

	isConnectable := isConnectable(g.logger, g.defaultOptions)
	if !isConnectable {
		return g.CreatingContainerOptions(context.Background(), t, options...)
	}

	opt := &config.RabbitmqHostOptions{
		UserName:    g.defaultOptions.UserName,
		Password:    g.defaultOptions.Password,
		HostName:    g.defaultOptions.Host,
		VirtualHost: g.defaultOptions.VirtualHost,
		Port:        g.defaultOptions.HostPort,
		HttpPort:    g.defaultOptions.HttpPort,
	}

	return opt, nil
}

func (g *rabbitmqDockerTest) Start(
	ctx context.Context,
	t *testing.T,
	serializer serializer.EventSerializer,
	rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
	options ...*contracts.RabbitMQContainerOptions,
) (bus.Bus, error) {
	rabbitmqHostOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	var mqBus bus.Bus
	if err := g.pool.Retry(func() error {
		config := &config.RabbitmqOptions{
			RabbitmqHostOptions: rabbitmqHostOptions,
		}
		conn, err := types.NewRabbitMQConnection(config)
		if err != nil {
			return err
		}

		mqBus, err = bus2.NewRabbitmqBus(
			config,
			serializer,
			g.logger,
			conn,
			rabbitmqBuilderFunc)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		g.logger.Errorf("Could not connect to docker: %s", err)
		return nil, err
	}

	return mqBus, nil
}

func (g *rabbitmqDockerTest) Cleanup(ctx context.Context) error {
	return nil
}

func (g *rabbitmqDockerTest) getRunOptions(
	opts ...*contracts.RabbitMQContainerOptions,
) *dockertest.RunOptions {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if len(option.Ports) > 0 {
			g.defaultOptions.Ports = option.Ports
		}
		if option.UserName != "" {
			g.defaultOptions.UserName = option.UserName
		}
		if option.Password != "" {
			g.defaultOptions.Password = option.Password
		}
		if option.Tag != "" {
			g.defaultOptions.Tag = option.Tag
		}
	}

	runOptions := &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"RABBITMQ_DEFAULT_USER=" + g.defaultOptions.UserName,
			"RABBITMQ_DEFAULT_PASS=" + g.defaultOptions.Password,
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: g.defaultOptions.Ports,
	}

	return runOptions
}

func isConnectable(logger logger.Logger, options *contracts.RabbitMQContainerOptions) bool {
	conn, err := amqp091.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d", options.UserName, options.Password, options.Host, options.HostPort),
	)
	if err != nil {
		logError(logger, options.UserName, options.Password, options.Host, options.HostPort)

		return false
	}

	defer conn.Close()

	if err != nil || (conn != nil && conn.IsClosed()) {
		logError(logger, options.UserName, options.Password, options.Host, options.HostPort)

		return false
	}
	logger.Infof(
		"Opened rabbitmq connection on host: %s",
		fmt.Sprintf("amqp://%s:%s@%s:%d", options.UserName, options.Password, options.Host, options.HostPort),
	)

	return true
}

func logError(logger logger.Logger, userName string, password string, host string, hostPort int) {
	// we should not use `t.Error` or `t.Errorf` for logging errors because it will `fail` our test at the end and, we just should use logs without error like log.Error (not log.Fatal)
	logger.Errorf(
		fmt.Sprintf(
			"Error in creating rabbitmq connection with %s",
			fmt.Sprintf("amqp://%s:%s@%s:%d", userName, password, host, hostPort),
		),
	)
}
