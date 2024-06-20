# 设置基础镜像为阿里云的 Go 1.21 镜像
FROM golang:1.21 AS builder
WORKDIR /src
COPY . .
# Accept the build argument
ARG HOST
ARG PORT
ARG USER
ARG DBNAME
ARG PASSWORD
# Make sure to use the ARG in ENV
ENV HOST=$HOST
ENV PORT=$PORT
ENV USER=$USER
ENV DBNAME=$DBNAME
ENV PASSWORD=$PASSWORD
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o go_web main.go

# 使用相同版本的镜像作为部署环境
FROM golang:1.21 AS deployer
WORKDIR /app
COPY --from=builder /src/go_web /app/
# 确保 ENTRYPOINT 指向正确的二进制文件
ENTRYPOINT ["/app/go_web"]