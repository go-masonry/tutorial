package validations

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	api "workshop-common/api"
)

type WorkshopValidations interface {
	AcceptCar(ctx context.Context, car *api.Car) error
	PaintCar(ctx context.Context, request *api.PaintCarRequest) error
	RetrieveCar(ctx context.Context, request *api.RetrieveCarRequest) error
	CarPainted(ctx context.Context, request *api.PaintFinishedRequest) error
}

type workshopValidations struct {
}

func CreateWorkshopValidations() WorkshopValidations {
	return new(workshopValidations)
}

func (w *workshopValidations) AcceptCar(ctx context.Context, car *api.Car) error {
	return carIdValidation(car.GetNumber())
}

func (w *workshopValidations) PaintCar(ctx context.Context, request *api.PaintCarRequest) error {
	supportedColors := map[string]struct{}{"red": {}, "green": {}, "blue": {}}
	if _, supported := supportedColors[strings.ToLower(request.GetDesiredColor())]; supported {
		return nil
	}
	return status.Errorf(codes.InvalidArgument, "out of ink for %s", request.GetDesiredColor())
}

func (w *workshopValidations) RetrieveCar(ctx context.Context, request *api.RetrieveCarRequest) error {
	return carIdValidation(request.GetCarNumber())
}

func (w *workshopValidations) CarPainted(ctx context.Context, request *api.PaintFinishedRequest) error {
	return carIdValidation(request.GetCarNumber())
}

func carIdValidation(carID string) error {
	if len(carID) != 8 {
		return status.Errorf(codes.InvalidArgument, "%s should be 8 chars long", carID)
	}
	return nil
}
