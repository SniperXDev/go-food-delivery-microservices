package dtos

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

// https://echo.labstack.com/guide/binding/
// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

// GetProductsRequestDto validation will handle in command level
type GetProductsRequestDto struct {
	*utils.ListQuery
}
