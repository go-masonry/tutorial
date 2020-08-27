package services

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	workshop "github.com/go-masonry/tutorial/04-instrumentation/api"
	"github.com/go-masonry/tutorial/04-instrumentation/app/controllers"
	"github.com/go-masonry/tutorial/04-instrumentation/app/validations"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

type subWorkshopServiceDeps struct {
	fx.In

	Logger      log.Logger
	Controller  controllers.SubWorkshopController
	Validations validations.SubWorkshopValidations
}

type subWorkshopImpl struct {
	deps subWorkshopServiceDeps
	workshop.UnimplementedSubWorkshopServer
}

func CreateSubWorkshopService(deps subWorkshopServiceDeps) workshop.SubWorkshopServer {
	return &subWorkshopImpl{
		deps: deps,
	}
}

func (s *subWorkshopImpl) PaintCar(ctx context.Context, request *workshop.SubPaintCarRequest) (*empty.Empty, error) {
	if err := s.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	s.deps.Logger.Debug(ctx, "sub workshop - actually painting the car")
	return s.deps.Controller.PaintCar(ctx, request)
}
