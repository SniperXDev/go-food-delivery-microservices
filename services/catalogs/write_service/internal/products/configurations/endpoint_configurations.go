package configurations

import (
	"context"
	"github.com/labstack/echo/v4"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	creatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/endpoints/v1"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product/endpoints/v1"
	gettingProductByIdV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/endpoints/v1"
	gettingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/endpoints/v1"
	searchingProductsV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product/endpoints/v1"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product/endpoints/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

func (c *productsModuleConfigurator) configEndpoints(ctx context.Context) {
	configV1Endpoints(c.EchoServer, c.InfrastructureConfiguration, ctx)
}

func configV1Endpoints(echoServer customEcho.EchoHttpServer, infra *infrastructure.InfrastructureConfiguration, ctx context.Context) {
	echoServer.ConfigGroup("/api/v1", func(v1 *echo.Group) {
		group := v1.Group("/products")
		productEndpointBase := &delivery.ProductEndpointBase{
			ProductsGroup:               group,
			InfrastructureConfiguration: infra,
		}
		// CreateNewProduct
		createProductEndpoint := creatingProductV1.NewCreteProductEndpoint(productEndpointBase)
		createProductEndpoint.MapRoute()

		// UpdateProductCommand
		updateProductEndpoint := updatingProductV1.NewUpdateProductEndpoint(productEndpointBase)
		updateProductEndpoint.MapRoute()

		// GetProductsQuery
		getProductsEndpoint := gettingProductsV1.NewGetProductsEndpoint(productEndpointBase)
		getProductsEndpoint.MapRoute()

		// SearchProductsQuery
		searchProducts := searchingProductsV1.NewSearchProductsEndpoint(productEndpointBase)
		searchProducts.MapRoute()

		// GetProductByIdQuery
		getProductByIdEndpoint := gettingProductByIdV1.NewGetProductByIdEndpoint(productEndpointBase)
		getProductByIdEndpoint.MapRoute()

		// DeleteProductCommand
		deleteProductEndpoint := deletingProductV1.NewDeleteProductEndpoint(productEndpointBase)
		deleteProductEndpoint.MapRoute()
	},
	)
}
