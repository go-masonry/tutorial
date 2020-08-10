package services

import (
	"context"

	workshop "github.com/go-masonry/tutorial/part2/api"
	"github.com/go-masonry/tutorial/part2/app/controllers"
	"github.com/go-masonry/tutorial/part2/app/validations"
	"github.com/golang/protobuf/ptypes/empty"
)

type workshopServiceDeps struct {
	Controller  controllers.WorkshopController
	Validations validations.WorkshopValidations
}

type workshopImpl struct {
	deps                                 workshopServiceDeps
	workshop.UnimplementedWorkshopServer // if you keep this one embedded even when you change your interface this code will compile
}

func CreateWorkshopService(deps workshopServiceDeps) workshop.WorkshopServer {
	return &workshopImpl{
		deps: deps,
	}
}

func (w *workshopImpl) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
	if err := w.deps.Validations.AcceptCar(ctx, car); err != nil {
		return nil, err
	}
	return w.deps.Controller.AcceptCar(ctx, car)
}

func (w *workshopImpl) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.PaintCar(ctx, request)
}

func (w *workshopImpl) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) (*workshop.Car, error) {
	if err := w.deps.Validations.RetrieveCar(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.RetrieveCar(ctx, request)
}

func (w *workshopImpl) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.CarPainted(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.CarPainted(ctx, request)
}
