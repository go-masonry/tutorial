# Tutorial

This tutorial will show how to build a gRPC web service using Mortar   

## Prerequisites 

- You should be familiar with [Protobuf](https://developers.google.com/protocol-buffers), [gRPC](https://grpc.io) and [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) / [Swagger](https://en.wikipedia.org/wiki/OpenAPI_Specification)
    - Install by starting [here](https://developers.google.com/protocol-buffers/docs/gotutorial) 
- You should understand what [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) is
- Have access to Jaeger service

## Workshop

In this tutorial we are going to build a Workshop web service. Our Workshop specializes in painting cars.
In order to paint a car first it needs to be accepted in our Workshop. Once accepted we can paint our car.
Once painted the customer can collect it.

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
Since we own our proto files we can "enrich" them, hence we are going to use option #1. We are going to add custom options to our proto files.
  
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

Plot twist, our Workshop is not going to do the actual painting but will delegate this task to another Workshop->SubWorkshop.
SubWorkshop will expose only one Endpoint:
- Paint Car

We will use protobuf again to describe it
```protobuf
service SubWorkshop{
  rpc PaintCar(SubPaintCarRequest) returns (google.protobuf.Empty);
}
```

> Our internal SubWorkshop service will only expose gRPC API

In this tutorial our service will implement both APIs

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

You can find all the generated files [here](api) within `tutorial/api` directory. 