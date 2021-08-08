module gin-api

go 1.13

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/why444216978/go-util v1.0.11
	go.uber.org/zap v1.17.0
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/mysql v1.1.0
	gorm.io/gorm v1.21.9
	gorm.io/plugin/dbresolver v1.1.0
)

replace github.com/gomodule/redigo v2.0.0+incompatible => github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
