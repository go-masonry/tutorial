package validations

import (
	"context"
	"fmt"

	workshop "github.com/go-masonry/tutorial/part2/api"
)

type SubWorkshopValidations interface {
	PaintCar(ctx context.Context, request *workshop.SubPaintCarRequest) error
}

type subWorkshopValidations struct{}

func CreateSubWorkshopValidations() SubWorkshopValidations {
	return new(subWorkshopValidations)
}
func (s subWorkshopValidations) PaintCar(ctx context.Context, request *workshop.SubPaintCarRequest) error {
	if len(request.GetCallbackServiceAddress()) == 0 {
		return fmt.Errorf("callback service address cannot be empty")
	}
	if request.GetCar() == nil {
		return fmt.Errorf("car can't be empty")
	}
	return nil
}
