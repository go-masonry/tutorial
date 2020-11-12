package controllers

import (
	api "workshop-common/api"
	"workshop/app/data"
)

// FromProtoCarToModelCar converts workshop proto model to our data Entity
func FromProtoCarToModelCar(car *api.Car) *data.CarEntity {
	if car == nil {
		return nil
	}
	return &data.CarEntity{
		CarNumber:     car.GetNumber(),
		Owner:         car.GetOwner(),
		BodyStyle:     api.CarBody_name[int32(car.GetBodyStyle())],
		OriginalColor: car.GetColor(),
		CurrentColor:  car.GetColor(),
	}
}

// FromModelCarToProtoCar converts our data Entity to workshop proto model
func FromModelCarToProtoCar(car *data.CarEntity) *api.Car {
	if car == nil {
		return nil
	}
	return &api.Car{
		Number:    car.CarNumber,
		Owner:     car.Owner,
		BodyStyle: api.CarBody(api.CarBody_value[car.BodyStyle]),
		Color:     car.CurrentColor,
	}
}
