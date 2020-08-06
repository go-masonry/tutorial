package services

import (
	"context"

	workshop "github.com/go-masonry/tutorial/part2/api"
	"github.com/go-masonry/tutorial/part2/app/controllers"
	"github.com/go-masonry/tutorial/part2/app/validations"
	"github.com/golang/protobuf/ptypes/empty"
)

type subWorkshopServiceDeps struct {
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
	return s.deps.Controller.PaintCar(ctx, request)
}
