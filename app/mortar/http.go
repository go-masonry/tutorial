package mortar

import (
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

func HttpClientFxOptions() fx.Option {
	return fx.Options(
		providers.HttpClientBuildersFxOption(), // client builders
		providers.TracerGRPCClientInterceptorFxOption(),
		providers.TracerRESTClientInterceptorFxOption(),
		providers.CopyGRPCHeadersClientInterceptorFxOption(),
	)
}

func HttpServerFxOptions() fx.Option {
	return fx.Options(
		providers.HttpServerBuilderFxOption(), // Web Server Builder
		providers.GRPCTracingUnaryServerInterceptorFxOption(),
		providers.GRPCGatewayMetadataTraceCarrierFxOption(), // read it's documentation to understand better
		providers.LoggerGRPCInterceptorFxOption(),
		providers.MonitorGRPCInterceptorFxOption(),
	)
}

// These will have you to debug/profile or understand the internals of your service
func InternalHttpHandlersFxOptions() fx.Option {
	return fx.Options(
		providers.InternalDebugHandlersFxOption(),
		providers.InternalProfileHandlerFunctionsFxOption(),
		providers.InternalSelfHandlersFxOption(),
	)
}
