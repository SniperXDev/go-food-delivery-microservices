package dtos

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
)

type GetOrdersResponseDto struct {
	Orders *utils.ListResult[*dtosV1.OrderReadDto]
}
