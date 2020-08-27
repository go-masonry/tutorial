# Tutorial - Part 6 Tests

Since this is a tutorial we are not going to write tests to cover all the code.
Instead we will focus on the business logic only.

## Fake REST API calls

Our Workshop logic remotely calls SubWorkshop using REST API. Although we can run a real SubWorkshop service that will answer our HTTP requests.
This has nothing to do with the code we really want to test. To remind you here is the dependencies that Workshop Controller needs to work.

```golang
type workshopControllerDeps struct {
  fx.In

  DB                data.CarDB
  Logger            log.Logger
  HTTPClientBuilder client.NewHTTPClientBuilder
}
```

You can manually create this struct and override only whats needed, but it's not interesting.
Instead we will use Fx to create everything needed for Workshop Controller.
Please look at [`workshop_test.go`](app/controllers/workshop_test.go) file to better understand this example.
> Please pay special attention to `TestPaintCar`. It also shows how you can use HTTP Client interceptor to avoid any Remote Calls.

## Fake gRPC API calls

While Workshop need to call SubWorkshop using REST API, SubWorkshop is calling back using gRPC.
This example also demonstrates the use of Mocked Interfaces.
> Every Mortar Interface have a Mock generated for it using [gomock](https://github.com/golang/mock).
> Mocked packages have a `mock_*` prefix
>
> - mock_client
> - mock_server
> - mock_trace
> - ...

Here are SubWorkshop dependencies

```golang
type subWorkshopControllerDeps struct {
  fx.In

  Logger            log.Logger
  GRPCClientBuilder client.GRPCClientConnectionBuilder
}
```

We mock `client.GRPCClientConnectionBuilder` and gRPC connection.

Please look at [`subworkshop_test.go`](app/controllers/subworkshop_test.go) file to better understand this example.

## Overriding Configuration values in Tests

It is sometimes convenient to override configuration values externally during tests.
We expect that the configuration library is capable of doing that.
In our case we use Viper and it allows to merge different configuration sources.
Hence you can provide 2 config files [`config.yml`, `config_test.yml`] to the builder, or even more.

Here is an example for tests

- `config_test.yml`

  ```yml
  mortar:
  name: "tutorial_test"
  logger:
    level: info
    console: true
  ```

- Mortar Fx Option configuration
  
  ```golang
    pwd, _ := os.Getwd()
    mortar.ViperFxOption(pwd+"/<relative-path>/config/config.yml", pwd+"/<relative-path>/config/config_test.yml"),
  ```
