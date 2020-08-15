FROM golang
MAINTAINER  weihaoyu
WORKDIR /go/gin-frame
COPY . .
EXPOSE 777
CMD ["/bin/bash", "/go/gin-frame/start.sh"]
