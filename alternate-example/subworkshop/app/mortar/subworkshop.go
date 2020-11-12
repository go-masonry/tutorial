package mortar

import (
	"context"

	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"subworkshop/app/controllers"
	"subworkshop/app/services"
	"subworkshop/app/validations"
	api "workshop-common/api"
)

type tutorialServiceDeps struct {
	fx.In

	// API Implementations
	SubWorkshop api.SubWorkshopServer
}

func SubWorkshopService() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: gRPCHandlers,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: gRPCGatewayHandlers,
		}),
		// All other dependencies
		serviceDependencies(),
	)
}

func gRPCHandlers(deps tutorialServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		api.RegisterSubWorkshopServer(srv, deps.SubWorkshop)
		// Any additional gRPC Implementations should be called here
	}
}

func gRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register sub api REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return api.RegisterSubWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func serviceDependencies() fx.Option {
	return fx.Provide(
		services.CreateSubWorkshopService,
		controllers.CreateSubWorkshopController,
		validations.CreateSubWorkshopValidations,
	)
}
