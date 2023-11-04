build:
	go build -o ./app -i main.go

run:
	go run main.go

.PHONY: build run