.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${BINARY} cmd/main.go

.PHONY: run
run:
	docker-compose up --build -d

.PHONY: stop
stop:
	docker-compose down
