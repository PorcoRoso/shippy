build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/porcorosso/shippy/vessel-service proto/vessel/vessel.proto
	# dep 工具暂不可用，直接手动编译
	GOOS=linux GOARCH=amd64 go build
	docker build -t vessel-service . --no-cache

run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service