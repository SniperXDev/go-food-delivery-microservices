package v1

import (
	"context"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/integration"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Get_All_Product_Query_Handler(t *testing.T) {
	fixture := integration.NewIntegrationTestFixture()

	err := mediatr.RegisterRequestHandler[*GetProductsQuery, *dtos.GetProductsResponseDto](NewGetProductsQueryHandler(fixture.Log, fixture.Cfg, fixture.ProductRepository))
	if err != nil {
		return
	}

	defer fixture.Cleanup()

	query := NewGetProductsQuery(utils.NewListQuery(10, 1))
	queryResult, err := mediatr.Send[*GetProductsQuery, *dtos.GetProductsResponseDto](context.Background(), query)

	assert.NotNil(t, queryResult)
	assert.NotNil(t, queryResult.Products)
	assert.NotEmpty(t, queryResult.Products.Items)
}
