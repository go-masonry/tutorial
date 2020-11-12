package mortar

import (
	"github.com/go-masonry/bjaeger"
	"github.com/go-masonry/bzerolog"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/providers"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"os"
)

func LoggerFxOption() fx.Option {
	return fx.Options(
		fx.Provide(zeroLogBuilder),
		providers.LoggerFxOption(),
		providers.LoggerGRPCIncomingContextExtractorFxOption(),
		bjaeger.TraceInfoContextExtractorFxOption(),
	)
}

func zeroLogBuilder(config cfg.Config) log.Builder {
	builder := bzerolog.Builder().IncludeCaller()
	if config.Get(mortar.LoggerWriterConsole).Bool() {
		builder = builder.SetWriter(bzerolog.ConsoleWriter(os.Stderr))
	}
	return builder
}


func CustomLoggerFxOption(serviceName string) fx.Option {
	return fx.Options(
		fx.Provide(jsonLogBuilder(serviceName)),
		providers.LoggerFxOption(),
		providers.LoggerGRPCIncomingContextExtractorFxOption(),
		bjaeger.TraceInfoContextExtractorFxOption(),
	)
}

func jsonLogBuilder(serviceName string) func() log.Builder {
	return func() log.Builder {
		builder := bzerolog.Builder().AddStaticFields(map[string]interface{}{"service":serviceName})
		zerolog.TimestampFieldName = "@timestamp"
		return builder
	}
}

