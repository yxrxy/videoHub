# 第一阶段：依赖阶段
FROM golang:1.21-alpine AS deps
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

# 第二阶段：构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
ARG SERVICE_PATH
RUN go build -o service_binary ./cmd/${SERVICE_PATH}

# 第三阶段：运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
ARG SERVICE_PORT
COPY --from=builder /app/config/config.yaml /app/config/
COPY --from=builder /app/service_binary /app/
EXPOSE ${SERVICE_PORT}
CMD ["./service_binary"]