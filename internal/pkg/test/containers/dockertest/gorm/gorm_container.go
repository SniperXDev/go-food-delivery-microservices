package gorm

import (
	"context"
	"log"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/phayes/freeport"
	"gorm.io/gorm"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/contracts"
)

type gormDockerTest struct {
	resource       *dockertest.Resource
	defaultOptions *contracts.PostgresContainerOptions
}

func NewGormDockerTest() contracts.GormContainer {
	return &gormDockerTest{
		defaultOptions: &contracts.PostgresContainerOptions{
			Database:  "test_db",
			Port:      "5432",
			Host:      "localhost",
			UserName:  "dockertest",
			Password:  "dockertest",
			Tag:       "latest",
			ImageName: "postgres",
			Name:      "postgresql-dockertest",
		},
	}
}

func (g *gormDockerTest) CreatingContainerOptions(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gormPostgres.GormOptions, error) {
	// https://github.com/ory/dockertest/blob/v3/examples/PostgreSQL.md
	// https://github.com/bozd4g/fb.testcontainers
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	runOption := g.getRunOptions(options...)
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(
		runOption,
		func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
	if err != nil {
		log.Fatalf("Could not start resource (Postgresql Test Container): %s", err)
	}

	resource.Expire(
		120,
	) // Tell docker to hard kill the container in 120 seconds exponential backoff-retry, because the application_exceptions in the container might not be ready to accept connections yet

	g.resource = resource
	port, _ := strconv.Atoi(resource.GetPort("5432/tcp"))
	g.defaultOptions.HostPort = port

	t.Cleanup(func() { _ = resource.Close() })

	var postgresoptions *gormPostgres.GormOptions

	if err = pool.Retry(func() error {
		postgresoptions = &gormPostgres.GormOptions{
			Port:     g.defaultOptions.HostPort,
			Host:     g.defaultOptions.Host,
			Password: g.defaultOptions.Password,
			DBName:   g.defaultOptions.Database,
			SSLMode:  false,
			User:     g.defaultOptions.UserName,
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
		return nil, err
	}

	return postgresoptions, nil
}

func (g *gormDockerTest) Start(
	ctx context.Context,
	t *testing.T,
	options ...*contracts.PostgresContainerOptions,
) (*gorm.DB, error) {
	gormOptions, err := g.CreatingContainerOptions(ctx, t, options...)
	if err != nil {
		return nil, err
	}

	db, err := gormPostgres.NewGorm(gormOptions)

	return db, nil
}

func (g *gormDockerTest) Cleanup(ctx context.Context) error {
	return g.resource.Close()
}

func (g *gormDockerTest) getRunOptions(
	opts ...*contracts.PostgresContainerOptions,
) *dockertest.RunOptions {
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

	hostFreePort, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	g.defaultOptions.HostPort = hostFreePort

	runOptions := &dockertest.RunOptions{
		Repository: g.defaultOptions.ImageName,
		Tag:        g.defaultOptions.Tag,
		Env: []string{
			"POSTGRES_USER=" + g.defaultOptions.UserName,
			"POSTGRES_PASSWORD=" + g.defaultOptions.Password,
			"POSTGRES_DB=" + g.defaultOptions.Database,
			"listen_addresses = '*'",
		},
		Hostname:     g.defaultOptions.Host,
		ExposedPorts: []string{g.defaultOptions.Port},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(g.defaultOptions.Port): {
				{HostIP: "0.0.0.0", HostPort: strconv.Itoa(g.defaultOptions.HostPort)},
			},
		},
	}

	return runOptions
}
