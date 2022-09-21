package external

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"time"
)

type ProductCreated struct {
	*types.Message
	ProductId   string    `json:"productId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}
