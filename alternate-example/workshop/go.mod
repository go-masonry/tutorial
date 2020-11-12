module workshop

go 1.14

require (
	github.com/alecthomas/kong v0.2.11
	github.com/go-masonry/mortar v0.1.2
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/spf13/afero v1.3.4 // indirect
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.1
	google.golang.org/grpc v1.33.0
	gopkg.in/ini.v1 v1.60.2 // indirect
	workshop-common v1.0.0
)

replace workshop-common => ../common
