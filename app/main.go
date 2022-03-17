package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	jobLib "github.com/why444216978/gin-api/library/job"
	httpServer "github.com/why444216978/gin-api/server/http"
	logMiddleware "github.com/why444216978/gin-api/server/http/middleware/log"
	panicMiddleware "github.com/why444216978/gin-api/server/http/middleware/panic"
	timeoutMiddleware "github.com/why444216978/gin-api/server/http/middleware/timeout"

	"github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/loader"
	"github.com/why444216978/gin-api/app/module/test/job/grpc"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/app/router"
	"github.com/why444216978/gin-api/bootstrap"
)

var job = flag.String("job", "", "is job")

func main() {
	log.Printf("Actual pid is %d", syscall.Getpid())

	flag.Parse()
	if *job != "" {
		jobLib.Handlers = map[string]jobLib.HandleFunc{
			"grpc-cmux": grpc.GrpcCmux,
		}
		jobLib.Handle(*job)
		return
	}

	if err := loader.Load(); err != nil {
		panic(err)
	}

	srv := httpServer.New(fmt.Sprintf(":%d", config.App.AppPort),
		httpServer.WithReadTimeout(time.Duration(config.App.ReadTimeout)*time.Millisecond),
		httpServer.WithWriteTimeout(time.Duration(config.App.WriteTimeout)*time.Millisecond),
		httpServer.WithRegisterRouter(router.RegisterRouter),
		httpServer.WithMiddlewares(
			panicMiddleware.ThrowPanic(),
			timeoutMiddleware.TimeoutMiddleware(time.Duration(config.App.ContextTimeout)*time.Millisecond),
			logMiddleware.LoggerMiddleware(),
		),
		httpServer.WithPprof(config.App.Pprof),
		httpServer.WithDebug(config.App.IsDebug),
	)

	if err := bootstrap.NewApp(srv, bootstrap.WithRegistry(resource.Registrar)).Start(); err != nil {
		log.Println(err)
	}
}
