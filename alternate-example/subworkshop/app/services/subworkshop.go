package services

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
	api "workshop-common/api"
	"subworkshop/app/controllers"
	"subworkshop/app/validations"
)

type subWorkshopServiceDeps struct {
	fx.In

	Logger      log.Logger
	Controller  controllers.SubWorkshopController
	Validations validations.SubWorkshopValidations
}

type subWorkshopImpl struct {
	deps subWorkshopServiceDeps
	api.UnimplementedSubWorkshopServer
}

func CreateSubWorkshopService(deps subWorkshopServiceDeps) api.SubWorkshopServer {
	return &subWorkshopImpl{
		deps: deps,
	}
}

func (s *subWorkshopImpl) PaintCar(ctx context.Context, request *api.SubPaintCarRequest) (*empty.Empty, error) {
	if err := s.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	s.deps.Logger.Debug(ctx, "sub api - actually painting the car")
	return s.deps.Controller.PaintCar(ctx, request)
}
