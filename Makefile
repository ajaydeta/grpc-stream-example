dependencies:
	@echo "Installing dependencies..."
	@go mod download && go mod tidy

proto:
	@echo "Generating proto files..."
	@protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative proto.proto

server:
	@echo "Running server..."
	@go run server.go

run:
	@echo "Running server and client..."
	@go run server.go & go run client.go