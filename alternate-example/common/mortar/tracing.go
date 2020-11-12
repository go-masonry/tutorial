package mortar

import (
	"context"
	"go.uber.org/multierr"
	"os"
	"strings"

	"github.com/go-masonry/bjaeger"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
)

func TracerFxOption() fx.Option {
	return fx.Provide(jaegerBuilder)
}

// This constructor assumes you have JAEGER environment variables set
//
// https://github.com/jaegertracing/jaeger-client-go#environment-variables
//
// Once built it will register Lifecycle hooks (connect on start, close on stop)
func jaegerBuilder(lc fx.Lifecycle, config cfg.Config, logger log.Logger) (opentracing.Tracer, error) {
	address := strings.Split(config.Get("jaeger.address").String(), ":")
	err := multierr.Combine(
		os.Setenv("JAEGER_AGENT_HOST", address[0]),
		os.Setenv("JAEGER_AGENT_PORT", address[1]),
		os.Setenv("JAEGER_SAMPLER_TYPE", config.Get("jaeger.sampler_type").String()),
		os.Setenv("JAEGER_SAMPLER_PARAM", config.Get("jaeger.sampler_param").String()),
	)
	if err != nil {
		return nil, err
	}

	openTracer, err := bjaeger.Builder().
		SetServiceName(config.Get(mortar.Name).String()).
		AddOptions(bjaeger.BricksLoggerOption(logger)). // verbose logging,
		Build()
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return openTracer.Connect(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return openTracer.Close(ctx)
		},
	})
	return openTracer.Tracer(), nil
}