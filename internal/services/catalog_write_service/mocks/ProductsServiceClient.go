// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	productsservice "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// ProductsServiceClient is an autogenerated mock type for the ProductsServiceClient type
type ProductsServiceClient struct {
	mock.Mock
}

type ProductsServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ProductsServiceClient) EXPECT() *ProductsServiceClient_Expecter {
	return &ProductsServiceClient_Expecter{mock: &_m.Mock}
}

// CreateProduct provides a mock function with given fields: ctx, in, opts
func (_m *ProductsServiceClient) CreateProduct(ctx context.Context, in *productsservice.CreateProductReq, opts ...grpc.CallOption) (*productsservice.CreateProductRes, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *productsservice.CreateProductRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.CreateProductReq, ...grpc.CallOption) (*productsservice.CreateProductRes, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.CreateProductReq, ...grpc.CallOption) *productsservice.CreateProductRes); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*productsservice.CreateProductRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *productsservice.CreateProductReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductsServiceClient_CreateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProduct'
type ProductsServiceClient_CreateProduct_Call struct {
	*mock.Call
}

// CreateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - in *products_service.CreateProductReq
//   - opts ...grpc.CallOption
func (_e *ProductsServiceClient_Expecter) CreateProduct(ctx interface{}, in interface{}, opts ...interface{}) *ProductsServiceClient_CreateProduct_Call {
	return &ProductsServiceClient_CreateProduct_Call{Call: _e.mock.On("CreateProduct",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProductsServiceClient_CreateProduct_Call) Run(run func(ctx context.Context, in *productsservice.CreateProductReq, opts ...grpc.CallOption)) *ProductsServiceClient_CreateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*productsservice.CreateProductReq), variadicArgs...)
	})
	return _c
}

func (_c *ProductsServiceClient_CreateProduct_Call) Return(_a0 *productsservice.CreateProductRes, _a1 error) *ProductsServiceClient_CreateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductsServiceClient_CreateProduct_Call) RunAndReturn(run func(context.Context, *productsservice.CreateProductReq, ...grpc.CallOption) (*productsservice.CreateProductRes, error)) *ProductsServiceClient_CreateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetProductById provides a mock function with given fields: ctx, in, opts
func (_m *ProductsServiceClient) GetProductById(ctx context.Context, in *productsservice.GetProductByIdReq, opts ...grpc.CallOption) (*productsservice.GetProductByIdRes, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *productsservice.GetProductByIdRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.GetProductByIdReq, ...grpc.CallOption) (*productsservice.GetProductByIdRes, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.GetProductByIdReq, ...grpc.CallOption) *productsservice.GetProductByIdRes); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*productsservice.GetProductByIdRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *productsservice.GetProductByIdReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductsServiceClient_GetProductById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductById'
type ProductsServiceClient_GetProductById_Call struct {
	*mock.Call
}

// GetProductById is a helper method to define mock.On call
//   - ctx context.Context
//   - in *products_service.GetProductByIdReq
//   - opts ...grpc.CallOption
func (_e *ProductsServiceClient_Expecter) GetProductById(ctx interface{}, in interface{}, opts ...interface{}) *ProductsServiceClient_GetProductById_Call {
	return &ProductsServiceClient_GetProductById_Call{Call: _e.mock.On("GetProductById",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProductsServiceClient_GetProductById_Call) Run(run func(ctx context.Context, in *productsservice.GetProductByIdReq, opts ...grpc.CallOption)) *ProductsServiceClient_GetProductById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*productsservice.GetProductByIdReq), variadicArgs...)
	})
	return _c
}

func (_c *ProductsServiceClient_GetProductById_Call) Return(_a0 *productsservice.GetProductByIdRes, _a1 error) *ProductsServiceClient_GetProductById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductsServiceClient_GetProductById_Call) RunAndReturn(run func(context.Context, *productsservice.GetProductByIdReq, ...grpc.CallOption) (*productsservice.GetProductByIdRes, error)) *ProductsServiceClient_GetProductById_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateProduct provides a mock function with given fields: ctx, in, opts
func (_m *ProductsServiceClient) UpdateProduct(ctx context.Context, in *productsservice.UpdateProductReq, opts ...grpc.CallOption) (*productsservice.UpdateProductRes, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *productsservice.UpdateProductRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.UpdateProductReq, ...grpc.CallOption) (*productsservice.UpdateProductRes, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *productsservice.UpdateProductReq, ...grpc.CallOption) *productsservice.UpdateProductRes); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*productsservice.UpdateProductRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *productsservice.UpdateProductReq, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProductsServiceClient_UpdateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateProduct'
type ProductsServiceClient_UpdateProduct_Call struct {
	*mock.Call
}

// UpdateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - in *products_service.UpdateProductReq
//   - opts ...grpc.CallOption
func (_e *ProductsServiceClient_Expecter) UpdateProduct(ctx interface{}, in interface{}, opts ...interface{}) *ProductsServiceClient_UpdateProduct_Call {
	return &ProductsServiceClient_UpdateProduct_Call{Call: _e.mock.On("UpdateProduct",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProductsServiceClient_UpdateProduct_Call) Run(run func(ctx context.Context, in *productsservice.UpdateProductReq, opts ...grpc.CallOption)) *ProductsServiceClient_UpdateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*productsservice.UpdateProductReq), variadicArgs...)
	})
	return _c
}

func (_c *ProductsServiceClient_UpdateProduct_Call) Return(_a0 *productsservice.UpdateProductRes, _a1 error) *ProductsServiceClient_UpdateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProductsServiceClient_UpdateProduct_Call) RunAndReturn(run func(context.Context, *productsservice.UpdateProductReq, ...grpc.CallOption) (*productsservice.UpdateProductRes, error)) *ProductsServiceClient_UpdateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// NewProductsServiceClient creates a new instance of ProductsServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductsServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductsServiceClient {
	mock := &ProductsServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
