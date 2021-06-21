module github.com/plyovchev/sorting-robot-go/sort/fulfillment-service

go 1.16

replace github.com/plyovchev/sorting-robot-go/sort/gen => ../gen

require (
	github.com/plyovchev/sorting-robot-go/sort/gen v0.0.0-00010101000000-000000000000
	github.com/preslavmihaylov/ordertocubby v0.0.0-20210617074346-1704d311e402
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0 // indirect
)
