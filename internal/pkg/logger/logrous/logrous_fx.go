package logrous

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("logrousfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		fx.Annotate(
			NewLogrusLogger,
			fx.As(new(logger.Logger)),
		),
		config.ProvideLogConfig,
	))

var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module("logrousfx",

		fx.Provide(config.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
	)
}
