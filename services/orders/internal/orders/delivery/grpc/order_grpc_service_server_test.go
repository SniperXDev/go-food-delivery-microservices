package grpc

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	ordersService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/e2e"
	"github.com/stretchr/testify/assert"
	"testing"
)

type OrderGrpcServiceTests struct {
	*testing.T
	*e2e.E2ETestFixture
	*OrderGrpcServiceServer
}

func TestRunner(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("GRPC", func(t *testing.T) {
		// Before running the tests
		orderGrpcService := NewOrderGrpcService(fixture.InfrastructureConfiguration)
		ordersService.RegisterOrdersServiceServer(fixture.GrpcServer.GetCurrentGrpcServer(), orderGrpcService)

		go func() {
			if err := fixture.GrpcServer.RunGrpcServer(nil); err != nil {
				fixture.Log.Errorf("(s.RunGrpcServer) err: %v", err)
			}
		}()

		orderGrpcServiceTests := OrderGrpcServiceTests{
			T:                      t,
			E2ETestFixture:         fixture,
			OrderGrpcServiceServer: orderGrpcService,
		}

		// Run Tests
		orderGrpcServiceTests.Test_GetOrder_By_Id()
		orderGrpcServiceTests.Test_Create_Order()

		// After running the tests
		fixture.GrpcServer.GracefulShutdown()
		fixture.Cleanup()
	})
}

func (p *OrderGrpcServiceTests) Test_Create_Order() {
	//request := &productService.CreateProductReq{
	//	Price:       gofakeit.Price(100, 1000),
	//	Name:        gofakeit.Name(),
	//	Description: gofakeit.AdjectiveDescriptive(),
	//}
	//
	//res, err := p.CreateProduct(context.Background(), request)
	//assert.NoError(p.T, err)
	//assert.NotZero(p.T, res.ProductID)
}

func (p *OrderGrpcServiceTests) Test_GetOrder_By_Id() {
	res, err := p.GetOrderByID(context.Background(), &ordersService.GetOrderByIDReq{OrderId: "97e2d953-ed25-4afb-8578-782cc5d365ba"})
	assert.NoError(p.T, err)
	assert.NotNil(p.T, res.Order)
	assert.Equal(p.T, res.Order.OrderId, "97e2d953-ed25-4afb-8578-782cc5d365ba")
}
