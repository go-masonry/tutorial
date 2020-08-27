# Tutorial - Part 2 Business logic

## Project Structure

This tutorial proposes a project structure.
Having a consistent project structure helps a lot, however you can create your own, if needed.
What's more important is that your services project structure should look as similar as possible.
This helps a lot when there are different projects/groups and handful of developers.

Let's make a brief overview

```s
.
├── api
├── app
│   ├── controllers
│   ├── data
│   ├── mortar
│   ├── services
│   └── validations
├── build
├── config
└── tests
```

- `api` directory stores this service API definitions.
- `app` directory holds our service code.
  - `services` directory implements our gRPC API interfaces, also an entry point.
  - `validations` directory treats everything related to input validations.
  - `data` directory is all about storing/persisting data.
  - `controllers` directory is responsible for business logic.
  - `mortar` later on this one.
- `build` directory stores everything you need to build your service (CI/CD).
- `config` directory stores everything related to your service configuration.
- `tests` directory stores functional/integration tests.
  
## Building Workshop

In this part we will build the following

- [Tutorial - Part 2 Business logic](#tutorial---part-2-business-logic)
  - [Project Structure](#project-structure)
  - [Building Workshop](#building-workshop)
    - [Services](#services)
    - [Validations](#validations)
    - [Controllers](#controllers)
    - [Fake DB](#fake-db)
    - [Building SubWorkshop](#building-subworkshop)

### Services

Every input starts within the `service`, once we get an input we want to `validate` it. Once validated we will call a `controller` to do the actual business logic. For the sake of brevity we will show only one Function. But feel free to browse the code.

```golang
func (w *workshopImpl) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
  if err := w.deps.Validations.AcceptCar(ctx, car); err != nil {
    return nil, err
  }
  return w.deps.Controller.AcceptCar(ctx, car)
}
```

As you can see I wasn't lying, we first `validate` then use the logic in the `controller` to treat this request.

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

Please remember that we also have a SubWorkshop service. Our business logic should call this SubWorkshop to actually paint the car, given the car was previously accepted by the Workshop.

>In this example we will use `*http.Client` a.k.a REST Client to call SubWorkshop API.

```golang
func (w *workshopController) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
  ...
  httpReq, err := w.makePaintRestRequest(ctx, car, request)
  ...
  response, err := w.client.Do(httpReq)
  ...
}
```

### Fake DB

Once our Workshop accepts a car it needs to store it somewhere. We will fake a DB by using a simple MAP `map[string]*CarEntity`

```golang
type CarEntity struct {
  CarNumber     string
  Owner         string
  BodyStyle     string
  OriginalColor string
  CurrentColor  string
  Painted       bool
}

type CarDB interface {
  InsertCar(ctx context.Context, car *CarEntity) error
  PaintCar(ctx context.Context, CarNumber string, newColor string) error
  GetCar(ctx context.Context, CarNumber string) (*CarEntity, error)
  RemoveCar(ctx context.Context, CarNumber string) (*CarEntity, error)
}
```

> As a practice don't use your external DTOs as your DB models/Entities.

### Building SubWorkshop

SubWorkshop needs to do one thing, paint the car. Once it paints the car it needs to tell the Workshop that it finished. Now if you look at the SubWorkshop Request it has a callback field. We will use the callback value as an address to callback the Workshop service. This time we will call Workshop gRPC API (one that wasn't exposed as REST).

```golang
func (s *subWorkshopController) PaintCar(ctx context.Context, request *workshop.SubPaintCarRequest) (*empty.Empty, error) {
  // Paint car
  ...
  // Dial back to caller
  conn, err := grpc.DialContext(ctx, request.GetCallbackServiceAddress(), grpc.WithInsecure())
  if err != nil {
    return nil, fmt.Errorf("car painted but we can't callback to %s, %w", request.GetCallbackServiceAddress(), err)
  }
  // Make client and call remote method
  workshopClient := workshop.NewWorkshopClient(conn)
  return workshopClient.CarPainted(ctx, &workshop.PaintFinishedRequest{CarNumber: request.GetCar().GetNumber(), DesiredColor: request.GetDesiredColor()}
}
```
