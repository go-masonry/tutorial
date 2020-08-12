package controllers

import (
	workshop "github.com/go-masonry/tutorial/complete/api"
	"github.com/go-masonry/tutorial/complete/app/data"
)

func FromProtoCarToModelCar(car *workshop.Car) *data.CarEntity {
	if car == nil {
		return nil
	}
	return &data.CarEntity{
		CarID:         car.GetNumber(),
		Owner:         car.GetOwner(),
		BodyStyle:     workshop.CarBody_name[int32(car.GetBodyStyle())],
		OriginalColor: car.GetColor(),
		CurrentColor:  car.GetColor(),
	}
}

func FromModelCarToProtoCar(car *data.CarEntity) *workshop.Car {
	if car == nil {
		return nil
	}
	return &workshop.Car{
		Number:    car.CarID,
		Owner:     car.Owner,
		BodyStyle: workshop.CarBody(workshop.CarBody_value[car.BodyStyle]),
		Color:     car.CurrentColor,
	}
}
