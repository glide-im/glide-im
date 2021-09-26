protoc --proto_path=./proto --go_out=./ ./proto/*.proto
protoc --go_out=plugins=grpc:. -I=./proto -I=. common.proto