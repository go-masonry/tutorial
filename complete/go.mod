module github.com/go-masonry/tutorial/complete

go 1.14

require (
	github.com/alecthomas/kong v0.2.11
	github.com/go-masonry/bjaeger v0.0.0-20200826110714-a69431b81d36
	github.com/go-masonry/bprometheus v0.0.0-20200826141841-db9728cb76a0
	github.com/go-masonry/bviper v0.0.0-20200806151318-2a173e4784f5
	github.com/go-masonry/bzerolog v0.0.0-20200819164318-647ac9dc481f
	github.com/go-masonry/mortar v0.0.0-20200826155723-3cc4cbd94e04
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/common v0.13.0 // indirect
	github.com/spf13/afero v1.3.4 // indirect
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.1
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/ini.v1 v1.60.1 // indirect
)

replace (
	github.com/go-masonry/bjaeger => /Users/talgendler/development/go/src/github.com/go-masonry/bjaeger
	github.com/go-masonry/bprometheus => /Users/talgendler/development/go/src/github.com/go-masonry/bprometheus
	github.com/go-masonry/bviper => /Users/talgendler/development/go/src/github.com/go-masonry/bviper
	github.com/go-masonry/bzerolog => /Users/talgendler/development/go/src/github.com/go-masonry/bzerolog
	github.com/go-masonry/mortar => /Users/talgendler/development/go/src/github.com/go-masonry/mortar
)
