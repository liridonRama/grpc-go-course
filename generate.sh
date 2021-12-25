protoc --proto_path=greet --go_out=greet  greet/greetpb/greet.proto

protoc --proto_path=greet --go_out=greet --go_opt=paths=source_relative  --go-grpc_out=greet --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto