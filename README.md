# Tutorial

This tutorial will explain how to build a gRPC web service using [go-masonry/mortar](https://github.com/go-masonry/mortar)

- [Tutorial](#tutorial)
  - [Prerequisites](#prerequisites)
  - [Workshop](#workshop)
    - [gRPC](#grpc)
    - [Adding REST](#adding-rest)
  - [SubWorkshop](#subworkshop)
  - [Generating code from our proto file](#generating-code-from-our-proto-file)
  - [Project Structure](#project-structure)
  - [Building Workshop](#building-workshop)
    - [Services](#services)
    - [Validations](#validations)
    - [Controllers](#controllers)
    - [Fake DB](#fake-db)
  - [Building SubWorkshop](#building-subworkshop)
  - [Dependency Injection using Uber-FX](#dependency-injection-using-uber-fx)
    - [Introducing Mortar and Bricks](#introducing-mortar-and-bricks)
    - [Back to code](#back-to-code)
  - [Wiring](#wiring)
    - [main.go](#maingo)
    - [Configuration](#configuration)
    - [Logger](#logger)
    - [Wiring WebService](#wiring-webservice)

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
  
## Building Workshop

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

## Building SubWorkshop

SubWorkshop needs to do one thing, paint the car. Once it paints the car it needs to tell the Workshop that it finished. Now if you look at the SubWorkshop Request it has a callback field. We will use the callback value as an address to call the Workshop service. This time we will call Workshop gRPC API (one that wasn't exposed as REST).

```golang
func (s *subWorkshopController) PaintCar(ctx context.Context, request *workshop.SubPaintCarRequest) (*empty.Empty, error) {
  // Paint car
  if err := s.doActualPaint(ctx, request.GetCar()); err != nil {
    return nil, err
  }
  wrapper := s.deps.GRPCClientBuilder.Build() // we will explain this part later
  // Dial back to caller
  conn, err := wrapper.Dial(ctx, request.GetCallbackServiceAddress(), grpc.WithInsecure())
  if err != nil {
    return nil, fmt.Errorf("car painted but we can't callback to %s, %w", request.GetCallbackServiceAddress(), err)
  }
  // Make client and call method
  workshopClient := workshop.NewWorkshopClient(conn)
  return workshopClient.CarPainted(ctx, &workshop.PaintFinishedRequest{CarId: request.GetCar().GetId(), DesiredColor: request.GetDesiredColor()})
}
```

Here we only showing you Controller implementation, but there are also validations and service implementation. Feel free to browse.

## Dependency Injection using [Uber-FX](https://github.com/uber-go/fx)

If you are unfamiliar with it, it's best to read/watch all about it. No seriously, do it.

- <https://www.youtube.com/watch?v=LDGKQY8WJEM>
- <https://godoc.org/go.uber.org/fx>

Mortar is heavily based on uber-fx and because of it your application will also be based on it. Once you see the benefits it offers you'll probably never write services without it.

### Introducing Mortar and Bricks

After all you are reading how to build a service with Mortar and so far we haven't mentioned it at all.

Try to think of Mortar as an [Abstract Type](https://en.wikipedia.org/wiki/Abstract_type) while your project implementation (and this Tutorial) can be seen as *Concrete Type*.

Mortar defines several interfaces **without implementing** them, here are some.

- [Logger](https://github.com/go-masonry/mortar/blob/master/interfaces/log/interfaces.go)
- [Config](https://github.com/go-masonry/mortar/blob/master/interfaces/cfg/interfaces.go)

The reason Mortar doesn't implement them is because we don't want to reinvent the wheel. There are a lot of great libraries for logging and configuration. This way you can choose which implementation to use in your project. However, there are already some implementations for [Logger](https://github.com/go-masonry/bzerolog) and [Config](https://github.com/go-masonry/bviper) interfaces. Feel free to add yours.

So why not use original libraries directly ???

There are several reasons for that. First we wanted to reduce [boilerplate code](https://en.wikipedia.org/wiki/Boilerplate_code). Second is that for Mortar to be used by different projects we need to make sure everyone speaks the same language or Interfaces for that matter. This way every project can choose it's own implementation and it will not influence other projects. Third one is "middleware", many projects and Zerolog is no exclusion build their API without `context.Context` in mind. However if you ever built a gRPC Web Service or any others you become used to propagate `context.Context` interface as first parameter of every public Function. Since `context.Context` is also a map it is used as a storage to pass around functions/libraries. Now we want to capitalize on that. There are use-cases where we want to extract some of the information stored in the `context.Context` and use it else where. It can be real useful to **automatically** extract fields from the `context.Context` and add them to a log line. Hence Mortar Logger Interface is defined with that in mind. You can read more about Mortar "middleware concept" here. Actually Mortar provides a lot of different middleware, more on that later.

### Back to code

We will revisit our code with `uber-fx` and Mortar interfaces in mind.
For example let's look at `workshopControllerDeps` 

```golang
type workshopControllerDeps struct {
  fx.In

  DB                db.CarDB
  Logger            log.Logger
  HttpClientBuilder partial.HttpClientPartialBuilder
}

func (w *workshopController) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
  err := w.deps.DB.InsertCar(ctx, FromProtoCarToModelCar(car))
  w.deps.Logger.WithError(err).Debug(ctx, "car accepted")
  return &empty.Empty{}, err
}
```

As you can see we introduced `fx.In` marker for our struct, and used Mortar Logger to log that we accepted a car. If you haven't read about `uber-fx` yet, `fx.In` marker will tell `uber-fx` to Inject all the Types that are Publicly defined in this struct. Even if the struct itself is private.

You should browse the code to get a feel of what changed.

Actually at this point we have all our service logic completed. Now we need to "wire" everything.

## Wiring

Like you probably noticed Mortar is heavily based on `uber-fx` and it also introduces some Interfaces. Now let us create all the dependencies and wire everything together.

### main.go

Like any other program our tutorial must have a `main.go` file. Personally I prefer to keep it as simple and concise as possible. However given what we want to achieve here it will be somewhat hard, to help ourselves read/change this code later we will break it to functions.

```golang
func createApplication(configFilePath string, additionalFiles []string) *fx.App {
  return fx.New(
    mortar.ViperFxOption(configFilePath, additionalFiles...), // Configuration map
    mortar.LoggerFxOption(),      // Logger
  )
}
```

Now before we continue with the explanations I just want to remind you (yes you should get yourself familiar with `uber-fx`) that we need to explain `uber-fx` how to build our dependency graph. Since GO lacks meta programming we will need to do it explicitly. Do note that Mortar provides a lot of predefined `fx.Option`s, more about them later.

To create a Dependency we need to have a function where it's return type is the Dependency. For example, here is a function that creates our Workshop Controller.

```func CreateWorkshopController(deps workshopControllerDeps) WorkshopController```

You can think of them as **[Constructors](https://en.wikipedia.org/wiki/Constructor_(object-oriented_programming))**. Fx have two options that accept constructors, `fx.Provide` and `fx.Invoke` the former will call your constructor only if it's dependency needed by another constructor while the later will **Eagerly** try to call the provided constructor and will create every other **provided** dependency that your "invoked" constructor needed and their respected dependencies. By doing so it creates your dependency graph.

### Configuration

As mentioned before, you need to import Implementations for Mortar, in this Tutorial we are going to use [Viper](https://github.com/spf13/viper) for `Config`.

It is good practice to use constants in your code instead of [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming)) and it's even better to set them outside your code either by providing a config file or reading from Environment variable. Mortar have a `Config` interface that we mentioned earlier that is used everywhere to read them external configurations. But we need to provide an implementation for it. You can build one yourself or use <https://github.com/go-masonry/bviper>. To build a configuration you need to provide a configuration file or use one found under `config/config.yml` in our tutorial.
While Mortar can be configured explicitly and that gives you total control over it. It is much comfortable to use it's defaults. To read them Mortar expects a dedicated Configuration key called **mortar**

```yaml
mortar:
  name: "tutorial"
  ...
```

If you noticed we had an empty directory `app/mortar` that we are now going to fill with code.

Let's look at `app/mortar/config.go` file

```golang
package mortar

import (
  "github.com/go-masonry/bviper"
  "github.com/go-masonry/mortar/interfaces/cfg"
  "go.uber.org/fx"
)

func ViperFxOption(configFilePath string, additionalFilePaths ...string) fx.Option {
  return fx.Provide(func() (cfg.Config, error) {
    builder := bviper.Builder().SetConfigFile(configFilePath)
    for _, extraFile := range additionalFilePaths {
      builder = builder.AddExtraConfigFile(extraFile)
    }
    return builder.Build()
  })
}
```

Remember your code is not the one calling the Constructors, Fx does it. Hence we can't tell it to use our custom parameters. But we can wrap this with a [Closure](https://en.wikipedia.org/wiki/Closure_(computer_programming)). This way we have a Constructor `func() (cfg.Config, error)` that accepts no parameter and can be safely called by Fx.

### Logger

I'm not going to explain why logging is important. I'll just say that Mortar have a `Logger` interface and we will use [Zerolog](https://github.com/rs/zerolog) to implement it.

```golang
func LoggerFxOption() fx.Option {
  return fx.Options(
    fx.Provide(ZeroLogBuilder),
    providers.LoggerFxOption(),
  )
}

func ZeroLogBuilder(config cfg.Config) log.Builder {
  builder := bzerolog.Builder()

  if config.Get(mortar.LoggerWriterConsole).Bool() {
    builder = builder.SetWriter(bzerolog.ConsoleWriter(os.Stderr))
  }
  return builder
}
```

If you look at this Constructor function

`func ZeroLogBuilder(config cfg.Config) log.Builder`

You see that it depends on `Config` which we provided earlier. However this Constructor function doesn't produce `Logger` instead it produces something called `log.Builder` which will later be used by Mortar to enrich our `Logger`.

This is why if you look above the Constructor function there is this line `providers.LoggerFxOption()`. I don't want to explain here how it enriches the `Logger` let's just say that it's output is the Logger itself. Here is how it's defined in Mortar.

```golang
// LoggerFxOption adds Default Logger to the graph
func LoggerFxOption() fx.Option {
  return fx.Provide(constructors.DefaultLogger)
}
```

### Wiring WebService

As you might recall we are building a Web Service here with gRPC and REST. It is time to introduce how one should configure Mortar Web Service. Like with any other dependency, we are going to use [go-grpc](https://grpc.io/docs/languages/go/basics/) and [grpc-gateway](https://grpc-ecosystem.github.io/grpc-gateway) to implement Http Web Service.

You can look at [grpc-server-example](https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go) example. Especially `func main()` there you can see how one can create and start a simple gRPC service.

You can look at [grpc-gateway-example](https://grpc-ecosystem.github.io/grpc-gateway/docs/usage.html) example. Section 6 where there is also an example of how to create and start grpc-gateway service.

One of Mortar goals is to reduce boilerplate code. However, we also want to give you total control on how to configure Mortars web service. Meaning you can configure both `grpc-server` and `grpc-gateway` the way you need. But, we have some good defaults which are mostly good for most cases.

You can read all about Mortar Http Interfaces both for client and server [here](). This tutorial shows how to use it's defaults.

To create Mortar web service you need to provide 2 options

1. Web Server Builder `providers.HttpServerBuilderFxOption()`
2. Invoke everything related to Web server `providers.CreateEntireWebServiceDependencyGraph()`
  
First option creates a Web Server Builder using **implicitly** provided configuration. Second one uses this Builder to create a Web Service and all it's dependencies while also adding `fx.Lifecycle` OnStart/OnStop hooks. Once we run our application, OnStart `fx.Lifecycle` hooks will be run and start our service.

But creating Web Service is not enough, we need also to create our Workshop and SubWorkshop Implementations. Let's look at our `main.go` again.

```golang
func createApplication(configFilePath string, additionalFiles []string) *fx.App {
  return fx.New(
    mortar.ViperFxOption(configFilePath, additionalFiles...), // Configuration map
    mortar.LoggerFxOption(),                                  // Logger
    mortar.HttpClientFxOptions(),
    mortar.HttpServerFxOptions(),
    // Tutorial service dependencies
    mortar.TutorialAPIsAndOtherDependenciesFxOption(), // register tutorial APIs
    // This one invokes all the above
    providers.CreateEntireWebServiceDependencyGraph(), // http server invoker
  )
}
```

You can see that we added 4 new dependencies to our graph. Well actually that's not true, there are several dependencies hiding behind these options. All these options will finally satisfy this Tutorial Logic.