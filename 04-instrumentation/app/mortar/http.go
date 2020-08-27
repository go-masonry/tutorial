package mortar

import (
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

func HttpClientFxOptions() fx.Option {
	return fx.Options(
		providers.HTTPClientBuildersFxOption(), // client builders
	)
}

func HttpServerFxOptions() fx.Option {
	return fx.Options(
		providers.HTTPServerBuilderFxOption(), // Web Server Builder
	)
}
