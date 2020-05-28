build:
	mkdir -p functions
	GOOS=linux GOARCH=amd64 go build -o functions/first ./first.go

