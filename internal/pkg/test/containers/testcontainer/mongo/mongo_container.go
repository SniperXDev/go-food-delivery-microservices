package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"
)

type mongoTestContainers struct {
	container      testcontainers.Container
	defaultOptions *contracts.MongoContainerOptions
}

func NewMongoTestContainers() contracts.MongoContainer {
	return &mongoTestContainers{
		defaultOptions: &contracts.MongoContainerOptions{
			Database:  "test_db",
			Port:      "27017/tcp",
			Host:      "localhost",
			UserName:  "testcontainers",
			Password:  "testcontainers",
			Tag:       "latest",
			ImageName: "mongo",
			Name:      "mongo-testcontainer",
		},
	}
}

func (g *mongoTestContainers) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongodb.MongoDbOptions, error) {
	// https://github.com/testcontainers/testcontainers-go
	// https://dev.to/remast/go-integration-tests-using-testcontainers-9o5
	containerReq := g.getRunOptions(options...)

	// TODO: Using Parallel Container
	dbContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	// get a free random host hostPort
	hostPort, err := dbContainer.MappedPort(ctx, nat.Port(g.defaultOptions.Port))
	if err != nil {
		return nil, err
	}
	g.defaultOptions.HostPort = hostPort.Int()

	host, err := dbContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	g.container = dbContainer

	// Clean up the container after the test is complete
	t.Cleanup(func() { _ = dbContainer.Terminate(ctx) })
	option := &mongodb.MongoDbOptions{
		User:     g.defaultOptions.UserName,
		Password: g.defaultOptions.Password,
		UseAuth:  false,
		Host:     host,
		Port:     g.defaultOptions.HostPort,
		Database: g.defaultOptions.Database,
	}

	return option, nil
}

func (g *mongoTestContainers) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.MongoContainerOptions,
) (*mongo.Client, error) {
	mongoOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := mongodb.NewMongoDB(mongoOptions)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (g *mongoTestContainers) Cleanup(ctx context.Context) error {
	return g.container.Terminate(ctx)
}

func (g *mongoTestContainers) getRunOptions(
	opts ...*contracts.MongoContainerOptions,
) testcontainers.ContainerRequest {
	if len(opts) > 0 && opts[0] != nil {
		option := opts[0]
		if option.ImageName != "" {
			g.defaultOptions.ImageName = option.ImageName
		}
		if option.Host != "" {
			g.defaultOptions.Host = option.Host
		}
		if option.Port != "" {
			g.defaultOptions.Port = option.Port
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

	containerReq := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("%s:%s", g.defaultOptions.ImageName, g.defaultOptions.Tag),
		ExposedPorts: []string{g.defaultOptions.Port},
		WaitingFor:   wait.ForListeningPort(nat.Port(g.defaultOptions.Port)),
		Hostname:     g.defaultOptions.Host,
		SkipReaper:   true,
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": g.defaultOptions.UserName,
			"MONGO_INITDB_ROOT_PASSWORD": g.defaultOptions.Password,
		},
	}

	return containerReq
}
