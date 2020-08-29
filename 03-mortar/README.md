# Tutorial - Part 3 Mortar and Uber-Fx

## Dependency Injection

If you are unfamiliar with it, it's best to read/watch all about it. No seriously, do it.
Mortar is heavily based on uber-fx and because of it your application will also be based on it.

- <https://www.youtube.com/watch?v=LDGKQY8WJEM>
- <https://godoc.org/go.uber.org/fx>

## Introducing Mortar and Bricks

After all you are reading how to build a service with Mortar and so far we haven't mentioned it at all. In fact the Business Logic in [02-logic](../02-logic/README.md) is not so special and can exist and work without Fx or Mortar.

Let's introduce Mortar

Try to think of Mortar as an [Abstract Type](https://en.wikipedia.org/wiki/Abstract_type) while your project implementation (and this Tutorial) can be seen as *Concrete Type*.

Mortar defines different interfaces **without implementing** all of them.

The reason Mortar doesn't implement them is because we don't want to reinvent the wheel. There are a lot of great libraries that solve different problems. One just need to make sure they are wrapped to implement Mortar Interfaces. Feel free to add yours.

By the way we call these external Implementations [Bricks](https://github.com/go-masonry/mortar/wiki/bricks.md)

So why not use original libraries directly ???

There are several reasons for that. First we wanted to reduce [boilerplate code](https://en.wikipedia.org/wiki/Boilerplate_code).
Second is that for Mortar to be used by different projects we need to make sure everyone speaks the same language or [Interfaces](https://github.com/go-masonry/mortar/tree/master/interfaces) for that matter.
This way every project can choose it's own implementation and it will not influence other projects.
Third one is [Middleware](https://github.com/go-masonry/mortar/blob/master/wiki/middleware.md), many projects/libraries build their API without `context.Context` in mind.
However if you ever built a gRPC Web Service or any other you become used to propagate `context.Context` interface as first parameter of at least every public Function.
Since `context.Context` is also a map it is used as a storage to pass around functions/libraries.
We want to capitalize on that. There are use-cases where we want to extract some of the information stored in the `context.Context` and use it else where.
It can be real useful to **automatically** extract fields from the `context.Context` and add them to a log line. Hence Mortar Logger Interface is defined with that in mind. Actually Mortar provides a lot of different middleware, more on that later.

## Back to code

We will revisit our code with Fx and Mortar in mind.
For example let's look at `workshopControllerDeps`

```golang
type workshopControllerDeps struct {
 fx.In

 DB                data.CarDB
 Logger            log.Logger
 HTTPClientBuilder client.NewHTTPClientBuilder
}

func (w *workshopController) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
  err := w.deps.DB.InsertCar(ctx, FromProtoCarToModelCar(car))
  w.deps.Logger.WithError(err).Debug(ctx, "car accepted")
  return &empty.Empty{}, err
}
```

As you can see we introduced `fx.In` marker to our struct, and used Mortar Logger to log that we accepted a car.
`fx.In` marker will tell Fx to Inject all the **Types** that are **Publicly** defined in this struct.
Even if the struct itself is private.

Actually at this point we have all our service logic revised with Mortar and Fx.
