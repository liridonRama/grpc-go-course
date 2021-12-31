protoc --proto_path=greet --go_out=greet --go_opt=paths=source_relative  --go-grpc_out=greet --go-grpc_opt=paths=source_relative greet/greetpb/greet.proto
protoc --proto_path=calculator --go_out=calculator --go_opt=paths=source_relative  --go-grpc_out=calculator --go-grpc_opt=paths=source_relative calculator/calculatorpb/calculator.proto
protoc --proto_path=blog --go_out=blog --go_opt=paths=source_relative  --go-grpc_out=blog --go-grpc_opt=paths=source_relative blog/blogpb/blog.proto



# To run mongodb
mongod --config /usr/local/etc/mongod.conf