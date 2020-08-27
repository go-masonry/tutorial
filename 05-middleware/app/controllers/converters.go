package controllers

import (
	workshop "github.com/go-masonry/tutorial/05-middleware/api"
	"github.com/go-masonry/tutorial/05-middleware/app/data"
)

// FromProtoCarToModelCar converts workshop proto model to our data Entity
func FromProtoCarToModelCar(car *workshop.Car) *data.CarEntity {
	if car == nil {
		return nil
	}
	return &data.CarEntity{
		CarNumber:     car.GetNumber(),
		Owner:         car.GetOwner(),
		BodyStyle:     workshop.CarBody_name[int32(car.GetBodyStyle())],
		OriginalColor: car.GetColor(),
		CurrentColor:  car.GetColor(),
	}
}

// FromModelCarToProtoCar converts our data Entity to workshop proto model
func FromModelCarToProtoCar(car *data.CarEntity) *workshop.Car {
	if car == nil {
		return nil
	}
	return &workshop.Car{
		Number:    car.CarNumber,
		Owner:     car.Owner,
		BodyStyle: workshop.CarBody(workshop.CarBody_value[car.BodyStyle]),
		Color:     car.CurrentColor,
	}
}
