package mortar

import (
	"context"

	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	api "workshop-common/api"
	"workshop/app/controllers"
	"workshop/app/data"
	"workshop/app/services"
	"workshop/app/validations"
)

type tutorialServiceDeps struct {
	fx.In

	// API Implementations
	Workshop api.WorkshopServer
}

func WorkshopService() fx.Option {
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
		tutorialDependencies(),
	)
}

func gRPCHandlers(deps tutorialServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		api.RegisterWorkshopServer(srv, deps.Workshop)
		// Any additional gRPC Implementations should be called here
	}
}

func gRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register api REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return api.RegisterWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func tutorialDependencies() fx.Option {
	return fx.Provide(
		services.CreateWorkshopService,
		controllers.CreateWorkshopController,
		data.CreateCarDB,
		validations.CreateWorkshopValidations,
	)
}
