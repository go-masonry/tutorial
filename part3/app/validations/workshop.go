package validations

import (
	"context"
	"fmt"
	"strings"

	workshop "github.com/go-masonry/tutorial/part3/api"
)

type WorkshopValidations interface {
	AcceptCar(ctx context.Context, car *workshop.Car) error
	PaintCar(ctx context.Context, request *workshop.PaintCarRequest) error
	RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) error
	CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) error
}

type workshopValidations struct {
}

func CreateWorkshopValidations() WorkshopValidations {
	return new(workshopValidations)
}

func (w *workshopValidations) AcceptCar(ctx context.Context, car *workshop.Car) error {
	return carIdValidation(car.GetId())
}

func (w *workshopValidations) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) error {
	supportedColors := map[string]struct{}{"red": {}, "green": {}, "blue": {}}
	if _, supported := supportedColors[strings.ToLower(request.GetDesiredColor())]; supported {
		return nil
	}
	return fmt.Errorf("out of ink for %s", request.GetDesiredColor())
}

func (w *workshopValidations) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) error {
	return carIdValidation(request.GetCarId())
}

func (w *workshopValidations) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) error {
	return carIdValidation(request.GetCarId())
}

func carIdValidation(carID string) error {
	if len(carID) != 8 {
		return fmt.Errorf("%s should be 8 chars long", carID)
	}
	return nil
}
