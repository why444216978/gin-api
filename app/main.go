package main

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/loader"
	"github.com/why444216978/gin-api/app/module/test/job/grpc/cmux"
	"github.com/why444216978/gin-api/app/module/test/job/grpc/h2c"
	serviceGRPC "github.com/why444216978/gin-api/app/module/test/service/grpc"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/app/router"
	"github.com/why444216978/gin-api/bootstrap"
	jobLib "github.com/why444216978/gin-api/library/job"
	serverGRPC "github.com/why444216978/gin-api/server/grpc"
	serverH2C "github.com/why444216978/gin-api/server/grpc/h2c"
	httpServer "github.com/why444216978/gin-api/server/http"
	logMiddleware "github.com/why444216978/gin-api/server/http/middleware/log"
	panicMiddleware "github.com/why444216978/gin-api/server/http/middleware/panic"
	timeoutMiddleware "github.com/why444216978/gin-api/server/http/middleware/timeout"
)

var (
	job    = flag.String("job", "", "is job")
	server = flag.String("server", "http", "is server type")
)

func main() {
	log.Printf("Actual pid is %d", syscall.Getpid())

	flag.Parse()
	if *job != "" {
		jobLib.Handlers = map[string]jobLib.HandleFunc{
			"grpc-cmux": cmux.Start,
			"grpc-h2c":  h2c.Start,
		}
		jobLib.Handle(*job)
		return
	}

	if err := loader.Load(); err != nil {
		panic(err)
	}

	port := config.App.AppPort
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

func startGRPC(port int) {
	srv := serverH2C.NewH2C(fmt.Sprintf(":%d", port),
		[]serverGRPC.Register{serviceGRPC.NewService()},
	)

	if err := bootstrap.NewApp(srv, bootstrap.WithRegistry(resource.Registrar)).Start(); err != nil {
		log.Println(err)
	}
}
