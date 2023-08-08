package otelMetrics

// ref:https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor add request status metrics to the otel
func UnaryServerInterceptor(meter api.Meter, serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		attrs := api.WithAttributes(
			attribute.Key("MetricsType").String("Grpc"),
		)

		if err != nil {
			counter, err := meter.Float64Counter(
				fmt.Sprintf("%s_error_grpc_requests_total", serviceName),
				api.WithDescription("The total number of error grpc requests"),
			)
			if err != nil {
				return nil, err
			}
			counter.Add(ctx, 1, attrs)
		} else {
			counter, err := meter.Float64Counter(fmt.Sprintf("%s_success_grpc_requests_total", serviceName), api.WithDescription("The total number of success grpc requests"))
			if err != nil {
				return nil, err
			}
			counter.Add(ctx, 1, attrs)
		}

		return resp, err
	}
}

// StreamServerInterceptor add request status metrics to the otel
func StreamServerInterceptor(meter api.Meter, serviceName string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)

		attrs := api.WithAttributes(
			attribute.Key("MetricsType").String("Grpc"),
		)

		ctx := ss.Context()

		if err != nil {
			counter, err := meter.Float64Counter(
				fmt.Sprintf("%s_error_grpc_requests_total", serviceName),
				api.WithDescription("The total number of error grpc requests"),
			)
			if err != nil {
				return err
			}
			counter.Add(ctx, 1, attrs)
		} else {
			counter, err := meter.Float64Counter(fmt.Sprintf("%s_success_grpc_requests_total", serviceName), api.WithDescription("The total number of success grpc requests"))
			if err != nil {
				return err
			}
			counter.Add(ctx, 1, attrs)
		}

		return err
	}
}
