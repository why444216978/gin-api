FROM golang:1.16 as builder
ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off

WORKDIR /gin-api/
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --ldflags '-extldflags "-static"' -o gin-api

FROM debian:latest
WORKDIR /gin-api/
COPY --from=builder /gin-api/gin-api /gin-api/gin-api
COPY --from=builder /gin-api/conf/ /gin-api/conf/
EXPOSE 777
CMD ["/gin-api/gin-api -env online"]
