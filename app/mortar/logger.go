package mortar

import (
	"github.com/go-masonry/bjaeger"
	"github.com/go-masonry/bzerolog"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
	"os"
)

func LoggerFxOption() fx.Option {
	return fx.Options(
		fx.Provide(ZeroLogBuilder),
		providers.LoggerFxOption(),
		providers.LoggerGRPCIncomingContextExtractorFxOption(),
	)
}

func ZeroLogBuilder(config cfg.Config) log.Builder {
	builder := bzerolog.
		Builder().
		// You can add explicit context extractors here or use the implicit fx.Group used by `go-masonry/mortar/constructors/logger.go`
		AddContextExtractors(bjaeger.TraceInfoExtractorFromContext)

	if config.Get(mortar.LoggerWriterConsole).Bool() {
		builder = builder.SetWriter(bzerolog.ConsoleWriter(os.Stderr))
	}
	return builder
}
