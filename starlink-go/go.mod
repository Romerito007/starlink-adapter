module github.com/Romerito007/starlink-adapter/starlink-go

go 1.22.0

require (
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.4
)

replace google.golang.org/grpc => ./internal/stub/grpc
