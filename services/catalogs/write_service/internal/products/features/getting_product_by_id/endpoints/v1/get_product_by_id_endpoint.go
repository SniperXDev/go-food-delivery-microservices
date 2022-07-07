package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id/dtos"
	"net/http"
)

type getProductByIdEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewGetProductByIdEndpoint(productEndpointBase *delivery.ProductEndpointBase) *getProductByIdEndpoint {
	return &getProductByIdEndpoint{productEndpointBase}
}

func (ep *getProductByIdEndpoint) MapRoute() {
	ep.ProductsGroup.GET("/:id", ep.getProductByID())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) getProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.GetProductByIdHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.getProductByID")
		defer span.Finish()

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := getting_product_by_id.NewGetProductById(request.ProductId)

		if err := ep.Validator.StructCtx(ctx, query); err != nil {
			ep.Log.WarnMsg("validate", err)
			tracing.TraceErr(span, err)
			return err
		}

		queryResult, err := ep.ProductMediator.Send(ctx, query)

		if err != nil {
			ep.Log.WarnMsg("GetProductById", err)
			tracing.TraceErr(span, err)
			return err
		}

		response, ok := queryResult.(*dtos.GetProductByIdResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, response)
	}
}
