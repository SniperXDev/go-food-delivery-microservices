package integrationEvents

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"

	uuid "github.com/satori/go.uuid"
)

type ProductDeletedV1 struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}

func NewProductDeletedV1(productId string) *ProductDeletedV1 {
	return &ProductDeletedV1{ProductId: productId, Message: types.NewMessage(uuid.NewV4().String())}
}
