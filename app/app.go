package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/why444216978/gin-api/config"
	job_service "github.com/why444216978/gin-api/library/job"
	"github.com/why444216978/gin-api/library/registry"
	"github.com/why444216978/gin-api/library/registry/etcd"
	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/router"
	"github.com/why444216978/gin-api/services/test/job/grpc"

	"github.com/why444216978/go-util/sys"
	"golang.org/x/sync/errgroup"
)

var (
	job = flag.String("job", "", "is job")
)

type App struct {
	ctx    context.Context
	server *http.Server
	cancel func()
}

func Start() {
	if *job != "" {
		job_service.Handlers = map[string]job_service.HandleFunc{
			"grpc-cmux": grpc.GrpcCmux,
		}
		job_service.Handle(*job)
		return
	}

	app := newApp()

	g, _ := errgroup.WithContext(app.ctx)
	//start serever
	g.Go(func() (err error) {
		err = app.start()
		return
	})
	g.Go(func() (err error) {
		app.registerSignal()
		return
	})
	g.Go(func() (err error) {
		err = app.registerService()
		if err != nil {
			panic(err)
		}
		return
	})
	g.Go(func() (err error) {
		app.shutdown()
		return
	})
	log.Printf("%s: errgroup exit %v\n", time.Now().Format("2006-01-02 15:04:05"), g.Wait())
}

func newApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	initResource(ctx)

	return &App{
		ctx:    ctx,
		cancel: cancel,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.App.AppPort),
			Handler:      router.InitRouter(),
			ReadTimeout:  time.Duration(config.App.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(config.App.WriteTimeout) * time.Millisecond,
		},
	}
}

func (a *App) start() error {
	err := a.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (a *App) registerSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Printf("%s: exit by signal %v\n", time.Now().Format("2006-01-02 15:04:05"), <-ch)

	//trigger shutdown
	a.cancel()
}

func (a *App) registerService() (err error) {
	var (
		localIP   string
		cfg       = &registry.RegistryConfig{}
		registrar *etcd.EtcdRegistrar
	)

	if err = resource.Config.ReadConfig("registry", "toml", cfg); err != nil {
		return err
	}

	if localIP, err = sys.LocalIP(); err != nil {
		return err
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
		return err
	}

	err = RegisterCloseFunc(registrar.DeRegister)

	return nil
}

func (a *App) shutdown() {
	<-a.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	//资源清理
	for _, f := range closeFunc {
		f(ctx)
	}

	//server shutdown
	err := a.server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
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
