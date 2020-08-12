module github.com/go-masonry/tutorial/complete

go 1.14

require (
	github.com/alecthomas/kong v0.2.11
	github.com/go-masonry/bjaeger v0.0.0-00010101000000-000000000000
	github.com/go-masonry/bviper v0.0.0-00010101000000-000000000000
	github.com/go-masonry/bzerolog v0.0.0-00010101000000-000000000000
	github.com/go-masonry/mortar v0.0.0-20200811123033-8cb0c6045bc1
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/opentracing/opentracing-go v1.2.0
	go.uber.org/fx v1.13.0
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20200729194436-6467de6f59a7 // indirect
	google.golang.org/genproto v0.0.0-20200804151602-45615f50871c
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace (
	//github.com/go-masonry/bdatadog => /Users/talgendler/development/go/src/github.com/go-masonry/bdatadog
	github.com/go-masonry/bjaeger => /Users/talgendler/development/go/src/github.com/go-masonry/bjaeger
	github.com/go-masonry/bviper => /Users/talgendler/development/go/src/github.com/go-masonry/bviper
	github.com/go-masonry/bzerolog => /Users/talgendler/development/go/src/github.com/go-masonry/bzerolog
	github.com/go-masonry/mortar => /Users/talgendler/development/go/src/github.com/go-masonry/mortar
)
