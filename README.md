# Tutorial

This tutorial will explain how to build a gRPC web service using [go-masonry/mortar](https://github.com/go-masonry/mortar)

## Prerequisites

- You should be familiar with [Protobuf](https://developers.google.com/protocol-buffers), [gRPC](https://grpc.io) and [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) / [Swagger](https://en.wikipedia.org/wiki/OpenAPI_Specification)
  - Install everything related to `gRPC` starting [here](https://developers.google.com/protocol-buffers/docs/gotutorial)
- You should understand what [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) is
- Have access to [Jaeger](https://www.jaegertracing.io/docs/1.18/getting-started) service

## Workshop

In this tutorial we are going to build a Workshop (Garage) web service. Our Workshop specializes in painting cars.
In order to paint a car first it needs to be *accepted* in our Workshop. Once accepted we can *paint* it.
Once painted the customer can *collect* it.

Given all the above we should expose these Endpoints:

- Accept Car
- Paint Car
- Retrieve Car

### gRPC

We will use protobuf to describe this service

```protobuf
service Workshop {
  rpc AcceptCar(Car) returns (google.protobuf.Empty);
  rpc PaintCar(PaintCarRequest) returns (google.protobuf.Empty);
  rpc RetrieveCar(RetrieveCarRequest) returns (Car);
}
```

Our service will expose a gRPC API based on the above protobuf definition.

### Adding REST

While gRPC is great, REST API however is still heavily used. If you look above, our service will have a gRPC implementation of the above API, why not reuse it ?
Well why not? To help us achieve this we will use [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) which will create a REST API reverse-proxy layer on top of our gRPC API.
If you are not familiar with it please read about it, it's an **amazing project**.

REST API still needs to be defined and grpc-gateway supports [1](https://grpc-ecosystem.github.io/grpc-gateway/docs/usage.html), [2](https://grpc-ecosystem.github.io/grpc-gateway/docs/grpcapiconfiguration.html) ways of doing it.
Since we own our proto files we can "enrich" them, hence we are going to use option 1. Let's go ahead and add custom options to our proto files.
  
```protobuf
import "google/api/annotations.proto";
service Workshop {
  rpc AcceptCar(Car) returns (google.protobuf.Empty){
    option (google.api.http) = {
      post: "/v1/workshop/cars"
      body: "*"
    };
  }

  rpc PaintCar(PaintCarRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      put: "/v1/workshop/cars/{car_id}/paint"
      body: "*"
    };
  }

  rpc RetrieveCar(RetrieveCarRequest) returns (Car) {
    option (google.api.http) = {
      get: "/v1/workshop/cars/{car_id}"
    };
  }
}
```

grpc-gateway protoc plugin will use provided `(google.api.http)` options to generate our reverse-proxy REST layer.

## SubWorkshop

**Plot twist**, our Workshop is not going to do the actual painting but will delegate this task to another Workshop a.k.a SubWorkshop.
SubWorkshop will expose only one Endpoint:

- Paint Car

Once our SubWorkshop paints the car it returns it to main Workshop. Since we want to show different features here (later on that), we will add an additional RPC to our Workshop. This new endpoint will only going to be exposed via gRPC.

```protobuf
service Workshop {
  ....
  rpc CarPainted(PaintFinishedRequest) returns (google.protobuf.Empty);
}
```

**In this tutorial our service will implement both APIs.**

To define SubWorkshop API we will use protobuf again to describe it

```protobuf
service SubWorkshop{
  rpc PaintCar(SubPaintCarRequest) returns (google.protobuf.Empty) {
    option (google.api.http) = {
      post: "/v1/subworkshop/paint"
      body: "*"
    };
  }
}
```

## Generating code from our proto file

If you haven't installed grpc-gateway plugin by now, please do.

To generate our code we will run the following command from within the `tutorial/api` directory

```shell script
protoc  -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=:. \
  --swagger_out=:. \
  --go_out=plugins=grpc:. \
  garage.proto
```

You can find all the generated files [here](api).

## Project Structure 

This tutorial proposes a project structure, it doesn't really matter which one you choose to create. What's more important is that your services project structure should look as similar as possible. This helps a lot when there are different projects/groups and handful of developers.

Let's make a brief overview

```shell
.
├── api
├── app
│   ├── controllers
│   ├── db
│   ├── mortar
│   ├── services
│   └── validations
├── build
├── config
└── tests
```

- `api` directory will store this service API definitions.
- `app` holds our service code.
  - `services` implements our gRPC API interfaces, also an entry point. Code here accepts all the input first.
  - `validations` treat everything related to validating input.
  - `db` take a guess.
  - `controllers` business logic lays here.
  - `mortar` later on this one.
- `build` you will need to build your service (CI/CD).
- `config` all configuration files should be here.
- `tests` functional/integration tests.
  
## Building our Workshop

In this part we will build the following

- [Services](#services)
- [Validations](#validations)
- [Controllers](#controllers)
- [DB](#fake-db)

### Services

Every input starts with `service`, once we get an input we want to `validate` it. Once validated we will call a `controller` to do the actual business logic. For the sake of brevity we will show only one Function

```golang
type workshopImpl struct {
  deps workshopServiceDeps
}

func CreateWorkshopService(deps workshopServiceDeps) workshop.WorkshopServer {
  return &workshopImpl{deps: deps}
}

func (w *workshopImpl) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
  if err := w.deps.Validations.AcceptCar(ctx, car); err != nil {
    return nil, err
  }
  return w.deps.Controller.AcceptCar(ctx, car)
}
```

As you can see I wasn't lying, we first `validate` then use the logic in the `controller` to treat this request.

So far nothing special, but you probably (or not) noticed an undefined struct called `workshopServiceDeps`. Here it is

```golang
type workshopServiceDeps struct {
  Controller  controllers.WorkshopController
  Validations validations.WorkshopValidations
}
```

### Validations

To simplify how we call our dependencies (validations, controllers) we defined Validation Interface to match our Service API. The only difference is that Validations functions return just an `error`.

```golang
type WorkshopValidations interface {
  AcceptCar(ctx context.Context, car *workshop.Car) error
  PaintCar(ctx context.Context, request *workshop.PaintCarRequest) error
  RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) error
  CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) error
}
```

### Controllers

Since Controller will implement the same exact gRPC API we can simply embed it.

```golang
type WorkshopController interface {
  workshop.WorkshopServer
}
```

If you were following this tutorial you will remember that we also have a SubWorkshop service. Our business logic should call this SubWorkshop to actually paint the car, given the car was previously accepted by the Workshop. In this example we will use `*http.Client` a.k.a REST Client to call SubWorkshop API.

```golang
type workshopController struct {
  deps    workshopControllerDeps // try to guess what we have here
  client  *http.Client
}

func (w *workshopController) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
  car, err := w.deps.DB.GetCar(ctx, request.GetCarId())
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
```

### Fake DB

Once our Workshop accepts a car it needs to store it somewhere. We will fake a DB by using a simple MAP `map[string]*CarEntity`

```golang
type CarEntity struct {
  CarID         string
  Owner         string
  BodyStyle     string
  OriginalColor string
  CurrentColor  string
  Painted       bool
}

type CarDB interface {
  InsertCar(ctx context.Context, car *CarEntity) error
  PaintCar(ctx context.Context, carID string, newColor string) error
  GetCar(ctx context.Context, carID string) (*CarEntity, error)
  RemoveCar(ctx context.Context, carID string) (*CarEntity, error)
}
```

> As a practice don't use your external DTOs as your DB models/Entities.

If you want to understand how everything should work, please take a look at the code within this directory.



## Dependency Injection using [Uber-FX](https://github.com/uber-go/fx)

If you are unfamiliar with it, it's best to read/watch all about it. No seriously, do it.

- https://www.youtube.com/watch?v=LDGKQY8WJEM
- https://godoc.org/go.uber.org/fx
