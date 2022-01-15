# 使用最新golang镜像
FROM golang:latest

# 作者
MAINTAINER AirGo

# 环境变量配置
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

# 指定工作目录，用于保存文件
WORKDIR /tmp

# 移动代码到容器中
COPY . .

# 指定工作目录，用于保存执行文件和配置文件
WORKDIR /app/

# 移动执行文件和配置文件
RUN cp /tmp/gin-api/app/app .
RUN cp -rf /tmp/gin-api/app/conf ./

# 向外暴露8888端口
EXPOSE 8888

# 启动容器时运行的命令
ENTRYPOINT ["./app -env online"]