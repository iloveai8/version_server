## 设置编译环境
FROM harbor.nuclearport.com/golang/golang:1.19 as builder
MAINTAINER chenhu@fotoable.com

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=gitlab.ftsview.com \
    GOPROXY=https://proxy.golang.org,direct

WORKDIR /app
COPY . .

RUN go build -o gsv .

## 设置运行环境
FROM alpine:latest
MAINTAINER chenhu@fotoable.com

WORKDIR /service

# 拷贝编译好的二进制文件
COPY --from=builder /app/gsv .
COPY --from=builder /app/conf/ conf/
COPY --from=builder /app/static/ static/
COPY --from=builder /app/data/ data/
EXPOSE 9091
ENTRYPOINT ./gsv run --config=./conf/${env}.yaml
