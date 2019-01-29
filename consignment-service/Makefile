build:
    # 一定要注意 Makefile 中的缩进，否则 make build 可能报错 Nothing to be done for build
    # protoc 命令前边是一个 Tab，不是四个或八个空格
    # grpc
	# protoc --go_out=plugins=grpc:$(GOPATH)/src/github.com/porcorosso/shippy/consignment-service ./proto/consignment/consignment.proto

    # 不再使用 grpc 插件，而使用go-micro
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/porcorosso/shippy/consignment-service proto/consignment/consignment.proto
    # 告知 Go 编译器生成二进制文件的目标环境：amd64 CPU 的 Linux 系统
	GOOS=linux GOARCH=amd64 go build
    # 根据当前目录下的 Dockerfile 生成名为 consignment-service 的镜像
	docker build -t consignment-service . --no-cache

run:
	# 在 Docker alpine 容器的 50001 端口上运行 consignment-service 服务
	# 可添加 -d 参数将微服务放到后台运行
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns consignment-service