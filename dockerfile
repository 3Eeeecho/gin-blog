# 第一阶段：构建应用
FROM golang:1.23-alpine AS builder

WORKDIR /go/src/ginblog
ENV GOPROXY=https://goproxy.cn,direct
# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译 Go 应用并生成二进制文件，直接生成在项目根目录
RUN CGO_ENABLED=0 GOOS=linux go build -o ginblog .

# 第二阶段：创建运行时镜像
FROM alpine:latest

# 将构建阶段生成的二进制文件复制到根目录
COPY --from=builder /go/src/ginblog/ginblog /ginblog

# 复制配置文件
COPY conf/app.ini /conf/app.ini

EXPOSE 8000
# 设置容器启动时的命令
CMD ["/ginblog"]
