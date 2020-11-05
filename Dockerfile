FROM golang:1.13 as builder
ENV GOPROXY=http://goproxy.ymt360.com
ENV GOSUMDB=off

WORKDIR /gin-api/
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --ldflags '-extldflags "-static"' -o gin-api

FROM debian:latest
WORKDIR /gin-api/
COPY --from=builder /gin-api/gin-api /gin-api/gin-api
COPY --from=builder /gin-api/configs/ /gin-api/configs/
EXPOSE 777
CMD ["/gin-api/gin-api"]