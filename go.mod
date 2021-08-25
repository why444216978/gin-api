module gin-api

go 1.13

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/alicebob/miniredis/v2 v2.15.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/go-redis/redismock/v8 v8.0.6
	github.com/golang/mock v1.6.0
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/gopherjs/gopherjs v0.0.0-20210803090616-8f023c250c89 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/why444216978/go-util v1.0.12
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.32.0
	gopkg.in/ini.v1 v1.57.0 // indirect
	gorm.io/driver/mysql v1.1.0
	gorm.io/gorm v1.21.9
	gorm.io/plugin/dbresolver v1.1.0
	icode.baidu.com/baidu/health/kylin v1.1.3 // indirect
)

replace github.com/gomodule/redigo v2.0.0+incompatible => github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
