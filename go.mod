module gin-api

go 1.13

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/google/uuid v1.1.2
	github.com/jinzhu/gorm v1.9.14
	github.com/larspensjo/config v0.0.0-20160228172812-b6db95dc6321
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/streadway/amqp v1.0.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/why444216978/go-util v0.0.0-20201219111016-08e42297ee4d
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	google.golang.org/grpc v1.32.0
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/ini.v1 v1.57.0
)

replace github.com/gomodule/redigo v2.0.0+incompatible => github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
