// Package openapi GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package openapi

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Mehdi Hadeli",
            "url": "https://github.com/mehdihadeli"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/orders": {
            "get": {
                "description": "Get all orders",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Get all orders",
                "parameters": [
                    {
                        "type": "string",
                        "name": "orderBy",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GetOrdersResponseDto"
                        }
                    }
                }
            },
            "post": {
                "description": "Create new order",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Create order",
                "parameters": [
                    {
                        "description": "Order data",
                        "name": "CreateOrderRequestDto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateOrderRequestDto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/dtos.CreateOrderResponseDto"
                        }
                    }
                }
            }
        },
        "/api/v1/orders/{id}": {
            "get": {
                "description": "Get order by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Orders"
                ],
                "summary": "Get order by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.GetOrderByIdResponseDto"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.CreateOrderRequestDto": {
            "type": "object",
            "properties": {
                "accountEmail": {
                    "type": "string"
                },
                "deliveryAddress": {
                    "type": "string"
                },
                "deliveryTime": {
                    "type": "string"
                },
                "shopItems": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ShopItemDto"
                    }
                }
            }
        },
        "dtos.CreateOrderResponseDto": {
            "type": "object",
            "properties": {
                "orderId": {
                    "type": "string"
                }
            }
        },
        "dtos.GetOrderByIdResponseDto": {
            "type": "object",
            "properties": {
                "order": {
                    "$ref": "#/definitions/dtos.OrderReadDto"
                }
            }
        },
        "dtos.GetOrdersResponseDto": {
            "type": "object",
            "properties": {
                "orders": {
                    "type": "object"
                }
            }
        },
        "dtos.OrderReadDto": {
            "type": "object",
            "properties": {
                "accountEmail": {
                    "type": "string"
                },
                "cancelReason": {
                    "type": "string"
                },
                "canceled": {
                    "type": "boolean"
                },
                "completed": {
                    "type": "boolean"
                },
                "createdAt": {
                    "type": "string"
                },
                "deliveredTime": {
                    "type": "string"
                },
                "deliveryAddress": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "orderId": {
                    "type": "string"
                },
                "paid": {
                    "type": "boolean"
                },
                "paymentId": {
                    "type": "string"
                },
                "shopItems": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ShopItemReadDto"
                    }
                },
                "submitted": {
                    "type": "boolean"
                },
                "totalPrice": {
                    "type": "number"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "dtos.ShopItemDto": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "dtos.ShopItemReadDto": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "quantity": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "utils.FilterModel": {
            "type": "object",
            "properties": {
                "comparison": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Orders Service Api",
	Description:      "Orders Service Api",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
