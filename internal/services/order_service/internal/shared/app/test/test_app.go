package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/store"
	config4 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/eventstroredb/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	config3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/eventstoredb"
	mongo2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/mongo"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/redis"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/configurations/orders"
	ordersService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/grpc/genproto"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                  *config.Config
	Bus                  bus.RabbitmqBus
	Container            contracts.Container
	Logger               logger.Logger
	RabbitmqOptions      *config2.RabbitmqOptions
	EchoHttpOptions      *config3.EchoHttpOptions
	EventStoreDbOptions  *config4.EventStoreDbOptions
	OrderMongoRepository repositories.OrderMongoRepository
	OrderAggregateStore  store.AggregateStore[*aggregate.Order]
	OrdersServiceClient  ordersService.OrdersServiceClient
	MongoClient          *mongo.Client
	EsdbClient           *esdb.Client
	MongoDbOptions       *mongodb.MongoDbOptions
	GrpcClient           grpc.GrpcClient
}

func NewTestApp() *TestApp {
	return &TestApp{}
}

func (a *TestApp) Run(t *testing.T) (result *TestAppResult) {
	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewOrdersTestApplicationBuilder(t)
	appBuilder.ProvideModule(orders.OrderServiceModule)

	appBuilder.Decorate(rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(eventstoredb.EventstoreDBContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(mongo2.MongoContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(redis.RedisContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	testApp.ConfigureOrders()

	testApp.MapOrdersEndpoints()

	testApp.ResolveFunc(
		func(
			cfg *config.Config,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			echoOptions *config3.EchoHttpOptions,
			grpcClient grpc.GrpcClient,
			eventStoreDbOptions *config4.EventStoreDbOptions,
			orderMongoRepository repositories.OrderMongoRepository,
			orderAggregateStore store.AggregateStore[*aggregate.Order],
			mongoClient *mongo.Client,
			esdbClient *esdb.Client,
			mongoDbOptions *mongodb.MongoDbOptions,
		) {
			result = &TestAppResult{
				Bus:                  bus,
				Cfg:                  cfg,
				Container:            testApp,
				Logger:               logger,
				RabbitmqOptions:      rabbitmqOptions,
				MongoClient:          mongoClient,
				MongoDbOptions:       mongoDbOptions,
				EchoHttpOptions:      echoOptions,
				EsdbClient:           esdbClient,
				EventStoreDbOptions:  eventStoreDbOptions,
				OrderMongoRepository: orderMongoRepository,
				OrderAggregateStore:  orderAggregateStore,
				OrdersServiceClient: ordersService.NewOrdersServiceClient(
					grpcClient.GetGrpcConnection(),
				),
				GrpcClient: grpcClient,
			}
		},
	)

	// we need a longer timout for up and running our testcontainers
	duration := time.Second * 300

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := testApp.Start(startCtx)
	if err != nil {
		t.Errorf("Error starting, err: %v", err)
		os.Exit(1)
	}

	//// waiting for grpc endpoint becomes ready in the given timeout
	//err = result.GrpcClient.WaitForAvailableConnection()
	//require.NoError(t, err)

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		err = testApp.Stop(stopCtx)
		require.NoError(t, err)
	})

	return
}
