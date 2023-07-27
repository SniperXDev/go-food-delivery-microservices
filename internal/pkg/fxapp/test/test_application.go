package test

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type testApplication struct {
	provides  []interface{}
	decorates []interface{}
	options   []fx.Option
	invokes   []interface{}
	logger    logger.Logger
	env       environemnt.Environment
	tb        fxtest.TB
	fxtestApp *fxtest.App
}

func (a *testApplication) Logger() logger.Logger {
	return a.logger
}

func (a *testApplication) Environment() environemnt.Environment {
	return a.env
}

func NewTestApplication(
	tb fxtest.TB,
	provides []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	env environemnt.Environment,
) contracts.Application {
	return &testApplication{
		tb:        tb,
		env:       env,
		logger:    logger,
		options:   options,
		provides:  provides,
		decorates: decorates,
	}
}

func (a *testApplication) ResolveFunc(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *testApplication) ResolveFuncWithParamTag(function interface{}, paramTagName string) {
	a.invokes = append(a.invokes, fx.Annotate(function, fx.ParamTags(paramTagName)))
}

func (a *testApplication) RegisterHook(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *testApplication) Run() {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxTestApp := CreateFxTestApp(
		a.tb,
		a.provides,
		a.decorates,
		a.invokes,
		a.options,
		a.logger,
		a.env,
	)
	// running phase will do in this stage and all register event hooks like OnStart and OnStop
	// instead of run for handling start and stop and create a ctx and cancel we can handle them manually with appconfigfx.start and appconfigfx.stop
	// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
	fxTestApp.Run()
}

func (a *testApplication) Start(ctx context.Context) error {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := CreateFxTestApp(a.tb, a.provides, a.decorates, a.invokes, a.options, a.logger, a.env)
	a.fxtestApp = fxApp

	return fxApp.Start(ctx)
}

func (a *testApplication) Stop(ctx context.Context) error {
	if a.fxtestApp == nil {
		a.logger.Fatal("Failed to stop because application not started.")
	}
	return a.fxtestApp.Stop(ctx)
}

func (a *testApplication) Wait() <-chan fx.ShutdownSignal {
	if a.fxtestApp == nil {
		a.logger.Fatal("Failed to wait because application not started.")
	}
	return a.fxtestApp.Wait()
}
