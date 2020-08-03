package mortar

import (
	"context"
	"github.com/go-masonry/mortar/constructors/partial"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	workshop "github.com/go-masonry/tutorial/api"
	"github.com/go-masonry/tutorial/app/controllers"
	"github.com/go-masonry/tutorial/app/db"
	"github.com/go-masonry/tutorial/app/services"
	"github.com/go-masonry/tutorial/app/validations"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type draftServiceDeps struct {
	fx.In

	// API Implementations
	Workshop    workshop.WorkshopServer
	SubWorkshop workshop.SubWorkshopServer
}

func DraftAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  partial.FxGroupGRPCServerAPIs,
			Target: draftGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  partial.FxGroupGRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: draftGRPCGatewayHandlers,
		}),
		// All other draft dependencies
		draftDependencies(),
	)
}

func draftGRPCServiceAPIs(deps draftServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		workshop.RegisterWorkshopServer(srv, deps.Workshop)
		workshop.RegisterSubWorkshopServer(srv, deps.SubWorkshop)
		// Any additional gRPC Implementations should be called here
	}
}

func draftGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register workshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return workshop.RegisterWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Register sub workshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return workshop.RegisterSubWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func draftDependencies() fx.Option {
	return fx.Provide(
		services.CreateWorkshopService,
		services.CreateSubWorkshopService,
		controllers.CreateWorkshopController,
		controllers.CreateSubWorkshopController,
		db.CreateCarDB,
		validations.CreateWorkshopValidations,
		validations.CreateSubWorkshopValidations,
	)
}
