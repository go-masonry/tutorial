package controllers

import (
	"context"
	"fmt"

	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	api "workshop-common/api"
)

// SubWorkshopController responsible for the business logic of our Sub Workshop
type SubWorkshopController interface {
	api.SubWorkshopServer
}

type subWorkshopControllerDeps struct {
	fx.In

	Logger            log.Logger
	GRPCClientBuilder client.GRPCClientConnectionBuilder
}

type subWorkshopController struct {
	deps subWorkshopControllerDeps
}

// CreateSubWorkshopController is a constructor to be used by Fx
func CreateSubWorkshopController(deps subWorkshopControllerDeps) SubWorkshopController {
	return &subWorkshopController{
		deps: deps,
	}
}

func (s *subWorkshopController) PaintCar(ctx context.Context, request *api.SubPaintCarRequest) (*empty.Empty, error) {
	// Paint car
	if err := s.doActualPaint(ctx, request.GetCar()); err != nil {
		return nil, err
	}
	wrapper := s.deps.GRPCClientBuilder.Build()
	// Dial back to caller
	conn, err := wrapper.Dial(ctx, request.GetCallbackServiceAddress(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("car painted but we can't callback to %s, %w", request.GetCallbackServiceAddress(), err)
	}
	// Make client and call method
	workshopClient := api.NewWorkshopClient(conn)
	return workshopClient.CarPainted(ctx, &api.PaintFinishedRequest{CarNumber: request.GetCar().GetNumber(), DesiredColor: request.GetDesiredColor()})
}

func (s *subWorkshopController) doActualPaint(ctx context.Context, car *api.Car) error {
	// here be paint logic...
	// ...
	// ...
	return nil
}
