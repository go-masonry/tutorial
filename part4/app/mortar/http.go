package mortar

import (
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

func HttpClientFxOptions() fx.Option {
	return fx.Options(
		providers.HttpClientBuildersFxOption(), // client builders
	)
}

func HttpServerFxOptions() fx.Option {
	return fx.Options(
		providers.HttpServerBuilderFxOption(), // Web Server Builder
	)
}
