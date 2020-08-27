package services

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	workshop "github.com/go-masonry/tutorial/06-tests/api"
	"github.com/go-masonry/tutorial/06-tests/app/controllers"
	"github.com/go-masonry/tutorial/06-tests/app/validations"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

type workshopServiceDeps struct {
	fx.In

	Logger      log.Logger
	Controller  controllers.WorkshopController
	Validations validations.WorkshopValidations
}

type workshopImpl struct {
	deps                                 workshopServiceDeps
	workshop.UnimplementedWorkshopServer // if keep this one added even when you change your interface this code will compile
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
	w.deps.Logger.WithField("car", car).Debug(ctx, "accepting car")
	return w.deps.Controller.AcceptCar(ctx, car)
}

func (w *workshopImpl) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "sending car to be painted")
	return w.deps.Controller.PaintCar(ctx, request)
}

func (w *workshopImpl) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) (*workshop.Car, error) {
	if err := w.deps.Validations.RetrieveCar(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "retrieving car")
	return w.deps.Controller.RetrieveCar(ctx, request)
}

func (w *workshopImpl) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.CarPainted(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "car painted")
	return w.deps.Controller.CarPainted(ctx, request)
}
