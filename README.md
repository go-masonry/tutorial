# Tutorial

This tutorial will explain how to build a gRPC web service using [go-masonry/mortar](https://github.com/go-masonry/mortar) step by step.

## Prerequisites

- You should be familiar with [Protobuf](https://developers.google.com/protocol-buffers), [gRPC](https://grpc.io) and [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) / [Swagger](https://en.wikipedia.org/wiki/OpenAPI_Specification)
  - Install everything related to `gRPC` starting [here](https://developers.google.com/protocol-buffers/docs/gotutorial)
- You should understand what [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) is
- Have access to [Jaeger](https://www.jaegertracing.io/docs/1.18/getting-started) service - **Optional**
- Have access to [Prometheus](https://prometheus.io) - **Optional**

## How to read this tutorial

### There are 7 parts in this tutorial, each part adds on top the previous one

1. [Tutorial service API](01-api)
2. [Project structure and business logic](02-logic)
3. [Dependency Injection and Mortar](03-mortar)
4. [Wiring/Instrumentation of everything together](04-instrumentation)
5. [Mortar Middleware](05-middleware)
6. [Tests example](06-tests)
7. [Makefile](07-makefile)

>You can create a local git repository and copy [01-api](01-api/) to it. Once you seen/understand the code, commit. Then copy [02-api](02-logic/) contents to your directory. This will overwrite some code.
>
>Repeat for each subsequent part. This way you will have git to show you what actually changed.
>
>*** Make sure to adjust imports accordingly.
