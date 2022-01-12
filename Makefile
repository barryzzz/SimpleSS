BINARY_NAME=ss
VERSION=1.0.0
build:
	go clean
	go build -o ./bin/$(BINARY_NAME)-client.exe cmd/local/main.go 
	go build -o ./bin/${BINARY_NAME}-server.exe cmd/server/main.go

clean:
	rm ./bin/$(BINARY_NAME)-client.exe
	rm ./bin/${BINARY_NAME}-server.exe