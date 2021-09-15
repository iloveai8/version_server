BINARY_NAME=GSVsnServer
MAIN_PATH=main.go

.PHONY:run build-linux deploy

run:
#	@echo start order
	go run -race main.go --config=./conf/dev.yaml

build-linux:
	@echo build linux
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_PATH)
#	upx $(BINARY_NAME)
    @echo "可执行文件大小为:`du -sh $(BINARY_NAME)`"

