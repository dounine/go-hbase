FROM golang:1.23-alpine3.19 as builder
# 设置工作目录
WORKDIR /app
# Go 代码复制到容器中
COPY . .
# 构建 Go 二进制文件
RUN go build -ldflags "-w -s" -o myapp .
# 使用轻量的 alpine 镜像作为最终的容器镜像
FROM alpine:3.14
# 设置工作目录
WORKDIR /app
# 从构建阶段复制 Go 二进制文件
COPY --from=builder /app/myapp .
# 暴露应用程序的端口
EXPOSE 8000
# 运行应用程序
CMD ["./myapp"]