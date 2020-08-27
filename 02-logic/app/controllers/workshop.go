package controllers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	workshop "github.com/go-masonry/tutorial/02-logic/api"
	"github.com/go-masonry/tutorial/02-logic/app/data"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
)

const (
	grpcServerPort   = "5380"
	externalRestPort = "5381"
)

// WorkshopController responsible for the business logic of our Workshop
type WorkshopController interface {
	workshop.WorkshopServer
}

type workshopControllerDeps struct {
	DB data.CarDB
}

type workshopController struct {
	deps    workshopControllerDeps
	client  *http.Client
	encoder *jsonpb.Marshaler
}

// CreateWorkshopController is a constructor for Fx
func CreateWorkshopController(deps workshopControllerDeps) WorkshopController {
	encoder := &jsonpb.Marshaler{OrigName: true}
	return &workshopController{
		deps:    deps,
		client:  http.DefaultClient,
		encoder: encoder,
	}
}

func (w *workshopController) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
	err := w.deps.DB.InsertCar(ctx, FromProtoCarToModelCar(car))
	return &empty.Empty{}, err
}

func (w *workshopController) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
	car, err := w.deps.DB.GetCar(ctx, request.GetCarNumber())
	if err != nil {
		return nil, err
	}
	httpReq, err := w.makePaintRestRequest(ctx, car, request)
	if err != nil {
		return nil, err
	}
	response, err := w.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("painting failed with status %d", response.StatusCode)
	}
	return &empty.Empty{}, nil
}

func (w *workshopController) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) (*workshop.Car, error) {
	car, err := w.deps.DB.GetCar(ctx, request.GetCarNumber())
	if err != nil {
		return nil, err
	}
	if car.Painted {
		car, err = w.deps.DB.RemoveCar(ctx, request.GetCarNumber())
		if err != nil {
			return nil, err
		}
		return FromModelCarToProtoCar(car), nil
	}
	return nil, fmt.Errorf("car %s is not painted", request.GetCarNumber())
}

func (w *workshopController) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) (*empty.Empty, error) {
	err := w.deps.DB.PaintCar(ctx, request.GetCarNumber(), request.GetDesiredColor())
	return &empty.Empty{}, err
}

func (w *workshopController) makePaintRestRequest(ctx context.Context, car *data.CarEntity, request *workshop.PaintCarRequest) (httpReq *http.Request, err error) {
	pbReq := &workshop.SubPaintCarRequest{
		Car:                    FromModelCarToProtoCar(car),
		DesiredColor:           request.GetDesiredColor(),
		CallbackServiceAddress: fmt.Sprintf(":%s", grpcServerPort),
	}
	body := new(bytes.Buffer)
	if err = w.encoder.Marshal(body, pbReq); err != nil {
		return nil, err
	}
	if httpReq, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/v1/subworkshop/paint", externalRestPort), body); err == nil {
		httpReq = httpReq.WithContext(ctx)
	}
	return
}
