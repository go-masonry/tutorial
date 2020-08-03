package db

import (
	"context"
	"fmt"
	"go.uber.org/fx"
)

// This interface will represent our car db
type CarDB interface {
	InsertCar(ctx context.Context, car *CarEntity) error
	PaintCar(ctx context.Context, carID string, newColor string) error
	GetCar(ctx context.Context, carID string) (*CarEntity, error)
	RemoveCar(ctx context.Context, carID string) (*CarEntity, error)
}

type carDBDeps struct {
	fx.In
}

func CreateCarDB(deps carDBDeps) CarDB {
	return &carDB{
		deps: deps,
		cars: make(map[string]*CarEntity),
	}
}

type carDB struct {
	deps carDBDeps
	cars map[string]*CarEntity
}

func (c *carDB) InsertCar(ctx context.Context, car *CarEntity) error {
	if _, exists := c.cars[car.CarID]; exists {
		return fmt.Errorf("car %s already exists", car.CarID)
	}
	c.cars[car.CarID] = car
	return nil
}

func (c *carDB) PaintCar(ctx context.Context, carID string, newColor string) error {
	if car, exists := c.cars[carID]; exists {
		car.CurrentColor = newColor
		car.Painted = true
		return nil
	}
	return fmt.Errorf("unknown car ID %s", carID)
}

func (c *carDB) GetCar(ctx context.Context, carID string) (*CarEntity, error) {
	if car, exists := c.cars[carID]; exists {
		return car, nil
	}
	return nil, fmt.Errorf("unknown car ID %s", carID)
}
func (c *carDB) RemoveCar(ctx context.Context, carID string) (*CarEntity, error) {
	car, err := c.GetCar(ctx, carID)
	if err == nil {
		delete(c.cars, carID)
		return car, nil
	}
	return nil, err
}
