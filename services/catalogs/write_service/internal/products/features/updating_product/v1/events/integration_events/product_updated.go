package integration_events

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	dto "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/dto/v1"
	uuid "github.com/satori/go.uuid"
)

type ProductUpdatedV1 struct {
	*types.Message
	*dto.ProductDto
}

func NewProductUpdatedV1(productDto *dto.ProductDto) *ProductUpdatedV1 {
	return &ProductUpdatedV1{Message: types.NewMessage(uuid.NewV4().String()), ProductDto: productDto}
}
