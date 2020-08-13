module github.com/go-masonry/tutorial/complete

go 1.14

require (
	github.com/alecthomas/kong v0.2.11
	github.com/go-masonry/bjaeger v0.0.0-00010101000000-000000000000
	github.com/go-masonry/bviper v0.0.0-00010101000000-000000000000
	github.com/go-masonry/bzerolog v0.0.0-00010101000000-000000000000
	github.com/go-masonry/mortar v0.0.0-20200813075556-b02ba8f61e29
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.0
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	google.golang.org/genproto v0.0.0-20200808173500-a06252235341
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
)

replace (
	//github.com/go-masonry/bdatadog => /Users/talgendler/development/go/src/github.com/go-masonry/bdatadog
	github.com/go-masonry/bjaeger => /Users/talgendler/development/go/src/github.com/go-masonry/bjaeger
	github.com/go-masonry/bviper => /Users/talgendler/development/go/src/github.com/go-masonry/bviper
	github.com/go-masonry/bzerolog => /Users/talgendler/development/go/src/github.com/go-masonry/bzerolog
	github.com/go-masonry/mortar => /Users/talgendler/development/go/src/github.com/go-masonry/mortar
)
