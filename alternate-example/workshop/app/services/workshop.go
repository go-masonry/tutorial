package services

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
	api "workshop-common/api"
	"workshop/app/controllers"
	"workshop/app/validations"
)

type workshopServiceDeps struct {
	fx.In

	Logger      log.Logger
	Controller  controllers.WorkshopController
	Validations validations.WorkshopValidations
}

type workshopImpl struct {
	deps                            workshopServiceDeps
	api.UnimplementedWorkshopServer // if keep this one added even when you change your interface this code will compile
}

func CreateWorkshopService(deps workshopServiceDeps) api.WorkshopServer {
	return &workshopImpl{
		deps: deps,
	}
}

func (w *workshopImpl) AcceptCar(ctx context.Context, car *api.Car) (*empty.Empty, error) {
	if err := w.deps.Validations.AcceptCar(ctx, car); err != nil {
		return nil, err
	}
	w.deps.Logger.WithField("car", car).Debug(ctx, "accepting car")
	return w.deps.Controller.AcceptCar(ctx, car)
}

func (w *workshopImpl) PaintCar(ctx context.Context, request *api.PaintCarRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "sending car to be painted")
	return w.deps.Controller.PaintCar(ctx, request)
}

func (w *workshopImpl) RetrieveCar(ctx context.Context, request *api.RetrieveCarRequest) (*api.Car, error) {
	if err := w.deps.Validations.RetrieveCar(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "retrieving car")
	return w.deps.Controller.RetrieveCar(ctx, request)
}

func (w *workshopImpl) CarPainted(ctx context.Context, request *api.PaintFinishedRequest) (*empty.Empty, error) {
	if err := w.deps.Validations.CarPainted(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.Debug(ctx, "car painted")
	return w.deps.Controller.CarPainted(ctx, request)
}
