# 指定基础镜像
FROM golang:1.20

ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 复制 Go 项目文件到容器中
COPY . .

# 构建 Go 项目
RUN go build -o main .

# 设置容器启动命令
CMD ["./main"]