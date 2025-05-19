# GoKeeper

docker run -p 8082:8082 keeper:v00

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/server.proto