gen-cal:
	protoc calculator\calculatorpb\calculator.proto --go-grpc_out=.
	protoc calculator\calculatorpb\calculator.proto --go_out=. 
run-server:
	go run server/server.go
run-client:
	go run client/client.go