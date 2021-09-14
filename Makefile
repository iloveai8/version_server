BINARY_NAME=GameSlotsVsn
MAIN_PATH=main.go

.PHONY:run build-linux

run:
	go run -race main.go run --config=conf/dev.yaml

build-linux:
	@echo build linux
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_PATH)
