package bootstrap

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/why444216978/go-util/sys"
	"golang.org/x/sync/errgroup"

	"github.com/why444216978/gin-api/app/config"
	"github.com/why444216978/gin-api/app/module/test/job/grpc"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/app/router"
	jobLib "github.com/why444216978/gin-api/library/job"
	"github.com/why444216978/gin-api/library/registry"
	"github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/server"
	httpServer "github.com/why444216978/gin-api/server/http"
	logMiddleware "github.com/why444216978/gin-api/server/http/middleware/log"
	panicMiddleware "github.com/why444216978/gin-api/server/http/middleware/panic"
	timeoutMiddleware "github.com/why444216978/gin-api/server/http/middleware/timeout"
)

var (
	job = flag.String("job", "", "is job")
)

type App struct {
	ctx    context.Context
	server server.RPCServer
	cancel func()
}

func Start() {
	flag.Parse()
	if *job != "" {
		jobLib.Handlers = map[string]jobLib.HandleFunc{
			"grpc-cmux": grpc.GrpcCmux,
		}
		jobLib.Handle(*job)
		return
	}

	app := newApp()

	g, _ := errgroup.WithContext(app.ctx)
	//start serever
	g.Go(func() (err error) {
		return app.start()
	})
	g.Go(func() (err error) {
		return app.registerSignal()
	})
	g.Go(func() (err error) {
		return app.registerService()
	})
	g.Go(func() (err error) {
		return app.shutdown()
	})
	log.Printf("%s: errgroup exit %v\n", time.Now().Format("2006-01-02 15:04:05"), g.Wait())
}

func newApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	initResource(ctx)

	app := &App{
		ctx:    ctx,
		cancel: cancel,
	}
	app.server = httpServer.New(app.ctx, fmt.Sprintf(":%d", config.App.AppPort),
		httpServer.WithReadTimeout(time.Duration(config.App.ReadTimeout)*time.Millisecond),
		httpServer.WithWriteTimeout(time.Duration(config.App.WriteTimeout)*time.Millisecond),
		httpServer.WithRegisterRouter(router.RegisterRouter),
		httpServer.WithMiddlewares(
			panicMiddleware.ThrowPanic(),
			timeoutMiddleware.TimeoutMiddleware(time.Duration(config.App.ContextTimeout)*time.Millisecond),
			logMiddleware.InitContext(),
			logMiddleware.LoggerMiddleware(),
		),
		httpServer.WithPprof(config.App.Pprof),
	)

	return app
}

func (a *App) start() error {
	return a.server.Start()
}

func (a *App) registerSignal() (err error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Printf("%s: exit by signal %v\n", time.Now().Format("2006-01-02 15:04:05"), <-ch)

	//trigger shutdown
	a.cancel()

	return
}

func (a *App) registerService() (err error) {
	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	var (
		localIP   string
		cfg       = &registry.RegistryConfig{}
		registrar *etcd.EtcdRegistrar
	)

	if err = resource.Config.ReadConfig("registry", "toml", cfg); err != nil {
		return
	}

	if localIP, err = sys.LocalIP(); err != nil {
		return
	}

	if resource.Etcd == nil || resource.Etcd.Client == nil {
		return
	}

	if registrar, err = etcd.NewRegistry(
		etcd.WithRegistrarClient(resource.Etcd.Client),
		etcd.WithRegistrarServiceName(config.App.AppName),
		etcd.WithRegistarHost(localIP),
		etcd.WithRegistarPort(config.App.AppPort),
		etcd.WithRegistrarLease(cfg.Lease)); err != nil {
		return
	}
	if err = registrar.Register(a.ctx); err != nil {
		return
	}

	if err = RegisterCloseFunc(registrar.DeRegister); err != nil {
		return
	}

	return nil
}

func (a *App) shutdown() (err error) {
	<-a.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	//server shutdown
	err = a.server.Close()

	//clean resource
	for _, f := range closeFunc {
		f(ctx)
	}

	return
}

// closeFunc 资源回收方法列表
var closeFunc = make([]func(ctx context.Context) error, 0)

// RegisterCloseFunc 注册资源回收方法
func RegisterCloseFunc(cf interface{}) error {
	f, ok := cf.(func(ctx context.Context) error)
	if !ok {
		return errors.New("func type error")
	}

	closeFunc = append(closeFunc, f)
	return nil
}
