module github.com/MikeCHOKKI/nexusflow/service-payment

go 1.22

require (
	github.com/MikeCHOKKI/nexusflow/protos/gen/go v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.64.0
)

replace github.com/MikeCHOKKI/nexusflow/protos/gen/go => ../protos/gen/go

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
