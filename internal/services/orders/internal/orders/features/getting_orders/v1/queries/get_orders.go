package queries

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

// Ref: https://golangbot.com/inheritance/

type GetOrders struct {
	*utils.ListQuery
}

func NewGetOrders(query *utils.ListQuery) *GetOrders {
	return &GetOrders{ListQuery: query}
}
