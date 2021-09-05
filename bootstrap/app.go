package bootstrap

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-api/config"
	"gin-api/jobs"
	"gin-api/libraries/registry"
	"gin-api/libraries/registry/etcd"
	"gin-api/resource"
	"gin-api/routers"

	"github.com/why444216978/go-util/sys"
	"golang.org/x/sync/errgroup"
)

var (
	job = flag.String("job", "", "is job")
)

type App struct {
	ctx       context.Context
	server    *http.Server
	registrar registry.Registrar
	cancel    func()
}

func StartApp() {
	if *job != "" {
		jobs.Handle(*job)
		return
	}

	app := newApp()
	go app.start()

	app.registerService()

	app.registerSignal()

	<-app.ctx.Done()
}

func newApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.App.AppPort),
			Handler:      routers.InitRouter(),
			ReadTimeout:  time.Duration(config.App.ReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(config.App.WriteTimeout) * time.Millisecond,
		},
	}
}

func (a *App) registerSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	timeout := time.Second * 3

	sig := <-ch
	var cancel context.CancelFunc
	a.ctx, cancel = context.WithTimeout(context.Background(), timeout)
	a.shutdown()
	cancel()
	log.Println(fmt.Sprintf("%s exit by signal %v\n", time.Now(), sig))
}

func (a *App) start() {
	g, _ := errgroup.WithContext(a.ctx)
	g.Go(func() (err error) {
		log.Println("start by server")
		log.Println("Start with " + a.server.Addr)
		err = a.server.ListenAndServe()
		return
	})
	err := g.Wait()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (a *App) registerService() {
	var (
		err     error
		localIP string
		cfg     = &registry.RegistryConfig{}
	)

	if err = resource.Config.ReadConfig("registry", "toml", cfg); err != nil {
		panic(err)
	}

	if localIP, err = sys.LocalIP(); err != nil {
		panic(err)
	}

	if a.registrar, err = etcd.NewRegistry(
		etcd.WithRegistrarClient(resource.Etcd.Client),
		etcd.WithRegistrarServiceName(config.App.AppName),
		etcd.WithRegistarHost(localIP),
		etcd.WithRegistarPort(config.App.AppPort),
		etcd.WithRegistrarLease(cfg.Lease)); err != nil {
		panic(err)
	}
	if err = a.registrar.Register(a.ctx); err != nil {
		panic(err)
	}
}

func (a *App) shutdown() {
	defer a.cancel()

	err := a.registrar.DeRegister(a.ctx)
	if err != nil {
		log.Println(err)
	}

	err = a.server.Shutdown(a.ctx)
	if err != nil {
		log.Println(err)
	}
}
