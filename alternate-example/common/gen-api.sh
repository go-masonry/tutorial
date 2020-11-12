protoc -I.\
        -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.0/third_party/googleapis \
        --grpc-gateway_out=repeated_path_param_separator=ssv:. \
        --openapiv2_out=repeated_path_param_separator=ssv:. \
        --go_out=plugins=grpc:api \
        api/garage.proto
