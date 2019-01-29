# 若运行环境是 Linux 则需把 alpine 换成 debian
# 使用最新版 alpine 作为基础镜像
FROM alpine:latest

# 在容器的根目录下创建 app 目录
RUN mkdir -p /app

# 将工作目录切换到 /app 下
WORKDIR /app

# 将当前目录下的货物信息文件 consignment.json 拷贝到 /app 目录下
ADD consignment.json /app/consignment.json
ADD vessel.json /app/vessel.json
ADD consignment-cli /app/consignment-cli

ENTRYPOINT ["./consignment-cli"]
# 运行服务端
CMD ["./consignment-cli"]
