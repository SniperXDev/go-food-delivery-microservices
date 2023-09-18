package createProductCommand

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type CreateProduct struct {
	ProductID   uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

func NewCreateProduct(name string, description string, price float64) (*CreateProduct, error) {
	command := &CreateProduct{
		ProductID:   uuid.NewV4(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}
	err := command.Validate()
	if err != nil {
		return nil, err
	}
	return command, nil
}

func (c *CreateProduct) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ProductID, validation.Required),
		validation.Field(&c.Name, validation.Required, validation.Length(0, 255)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 5000)),
		validation.Field(&c.Price, validation.Required, validation.Min(0.0).Exclusive()),
		validation.Field(&c.CreatedAt, validation.Required),
	)
}
