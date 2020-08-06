package mortar

import (
	"os"

	"github.com/go-masonry/bzerolog"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

func LoggerFxOption() fx.Option {
	return fx.Options(
		fx.Provide(ZeroLogBuilder),
		providers.LoggerFxOption(),
	)
}

func ZeroLogBuilder(config cfg.Config) log.Builder {
	builder := bzerolog.Builder()

	if config.Get(mortar.LoggerWriterConsole).Bool() {
		builder = builder.SetWriter(bzerolog.ConsoleWriter(os.Stderr))
	}
	return builder
}
