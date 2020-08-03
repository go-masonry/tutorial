package controllers

import (
	workshop "github.com/go-masonry/tutorial/api"
	"github.com/go-masonry/tutorial/app/db"
)

func FromProtoCarToModelCar(car *workshop.Car) *db.CarEntity {
	if car == nil {
		return nil
	}
	return &db.CarEntity{
		CarID:         car.GetId(),
		Owner:         car.GetOwner(),
		BodyStyle:     workshop.CarBody_name[int32(car.GetBodyStyle())],
		OriginalColor: car.GetColor(),
		CurrentColor:  car.GetColor(),
	}
}

func FromModelCarToProtoCar(car *db.CarEntity) *workshop.Car {
	if car == nil {
		return nil
	}
	return &workshop.Car{
		Id:        car.CarID,
		Owner:     car.Owner,
		BodyStyle: workshop.CarBody(workshop.CarBody_value[car.BodyStyle]),
		Color:     car.CurrentColor,
	}
}
