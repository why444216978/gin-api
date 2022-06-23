package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/why444216978/gin-api/app/loader"
	jobGRPC "github.com/why444216978/gin-api/app/module/test/job/grpc"
	serviceGRPC "github.com/why444216978/gin-api/app/module/test/service/grpc"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/app/router"
	"github.com/why444216978/gin-api/bootstrap"
	"github.com/why444216978/gin-api/library/app"
	jobLib "github.com/why444216978/gin-api/library/job"
	serverGRPC "github.com/why444216978/gin-api/server/grpc"
	serverH2C "github.com/why444216978/gin-api/server/grpc/h2c"
	httpServer "github.com/why444216978/gin-api/server/http"
	logMiddleware "github.com/why444216978/gin-api/server/http/middleware/log"
	panicMiddleware "github.com/why444216978/gin-api/server/http/middleware/panic"
	timeoutMiddleware "github.com/why444216978/gin-api/server/http/middleware/timeout"
)

var (
	env    = flag.String("env", "dev", "config path")
	job    = flag.String("job", "", "is job")
	server = flag.String("server", "http", "is server type")
)

func main() {
	flag.Parse()

	var err error

	if err = bootstrap.Init(*env, loader.Load); err != nil {
		log.Printf("bootstrap.Init err %s", err.Error())
		return
	}

	if *job != "" {
		jobLib.Handlers = map[string]jobLib.HandleFunc{
			"grpc-test": jobGRPC.Start,
		}
		jobLib.Handle(*job, resource.ServiceLogger)
		return
	}

	port := app.Port()
	if *server == "http" {
		log.Printf("start http, port %d", port)
		startHTTP(port)
	} else {
		log.Printf("start grpc, port %d", port)
		startGRPC(port)
	}
}

func startHTTP(port int) {
	srv := httpServer.New(fmt.Sprintf(":%d", port),
		httpServer.WithReadTimeout(app.ReadTimeout()),
		httpServer.WithWriteTimeout(app.WriteTimeout()),
		httpServer.WithRegisterRouter(router.RegisterRouter),
		httpServer.WithMiddlewares(
			panicMiddleware.ThrowPanic(resource.ServiceLogger),
			timeoutMiddleware.TimeoutMiddleware(app.ContextTimeout()),
			logMiddleware.LoggerMiddleware(resource.ServiceLogger),
		),
		httpServer.WithPprof(app.Pprof()),
		httpServer.WithDebug(app.Debug()),
	)

	if err := bootstrap.NewApp(srv, resource.Registrar).Start(); err != nil {
		log.Println(err)
	}
}

func startGRPC(port int) {
	srv := serverH2C.NewH2C(fmt.Sprintf(":%d", port),
		[]serverGRPC.Register{serviceGRPC.NewService()},
	)

	if err := bootstrap.NewApp(srv, resource.Registrar).Start(); err != nil {
		log.Println(err)
	}
}
