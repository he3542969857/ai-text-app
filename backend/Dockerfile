# ---- build ----
FROM golang:1.26-alpine AS build
WORKDIR /src
# 使用国内代理加速依赖拉取(direct 作为兜底,海外环境同样可用)
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

# ---- runtime ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
