FROM golang:1.13 as builder
ENV GOPROXY=http://goproxy.ymt360.com
ENV GOSUMDB=off

WORKDIR /purchase-server/
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build --ldflags '-extldflags "-static"' -o purchase-server

FROM debian:latest
WORKDIR /purchase-server/
COPY --from=builder /purchase-server/purchase-server /purchase-server/purchase-server
COPY --from=builder /purchase-server/configs/ /purchase-server/configs/
EXPOSE 777
CMD ["/purchase-server/purchase-server"]