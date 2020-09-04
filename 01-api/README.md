# Tutorial - Part 1 API

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
Since we own our `proto` files we can "enrich" them, hence we are going to use option 1. Let's go ahead and add custom options to our `proto` files.
  
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
      put: "/v1/workshop/cars/{car_number}/paint"
      body: "*"
    };
  }

  rpc RetrieveCar(RetrieveCarRequest) returns (Car) {
    option (google.api.http) = {
      get: "/v1/workshop/cars/{car_number}"
    };
  }
}
```

grpc-gateway protoc plugin will use the provided `(google.api.http)` options to generate our reverse-proxy REST layer.

### SubWorkshop

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

### Generating code from our garage.proto file

If you haven't installed grpc-gateway plugin by now, please do.

To generate our code we will run the following command from within the `api` directory

```shell
protoc  -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=repeated_path_param_separator=ssv:. \
  --openapiv2_out=repeated_path_param_separator=ssv:. \
  --go_out=plugins=grpc:. \
  garage.proto
```

You can find all the files [here](api/).
