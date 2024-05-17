package internal

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/s-larionov/process-manager"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/goverland-labs/goverland-ipfs-fetcher/protocol/ipfsfetcherpb"

	"github.com/goverland-labs/goverland-ipfs-fetcher/internal/config"
	"github.com/goverland-labs/goverland-ipfs-fetcher/internal/message"
	"github.com/goverland-labs/goverland-ipfs-fetcher/pkg/grpcsrv"
	"github.com/goverland-labs/goverland-ipfs-fetcher/pkg/health"
	"github.com/goverland-labs/goverland-ipfs-fetcher/pkg/prometheus"
)

type Application struct {
	sigChan <-chan os.Signal
	manager *process.Manager
	cfg     config.App
	db      *gorm.DB
}

func NewApplication(cfg config.App) (*Application, error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	a := &Application{
		sigChan: sigChan,
		cfg:     cfg,
		manager: process.NewManager(),
	}

	err := a.bootstrap()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Application) Run() {
	a.manager.StartAll()
	a.registerShutdown()
}

func (a *Application) bootstrap() error {
	initializers := []func() error{
		//a.initDB,

		// Init Dependencies
		a.initServices,

		// Init Workers: System
		a.initPrometheusWorker,
		a.initHealthWorker,
	}

	for _, initializer := range initializers {
		if err := initializer(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Application) initDB() error {
	db, err := gorm.Open(postgres.Open(a.cfg.DB.DSN), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		}),
	})
	if err != nil {
		return err
	}

	ps, err := db.DB()
	if err != nil {
		return err
	}
	ps.SetMaxOpenConns(a.cfg.DB.MaxOpenConnections)

	a.db = db
	if a.cfg.DB.Debug {
		a.db = db.Debug()
	}

	return err
}

func (a *Application) initServices() error {
	nc, err := nats.Connect(
		a.cfg.Nats.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(a.cfg.Nats.MaxReconnects),
		nats.ReconnectWait(a.cfg.Nats.ReconnectTimeout),
	)
	if err != nil {
		return err
	}

	pb, err := natsclient.NewPublisher(nc)
	if err != nil {
		return err
	}

	_ = pb

	err = a.initAPI()
	if err != nil {
		return fmt.Errorf("init api: %w", err)
	}

	return nil
}

func (a *Application) initAPI() error {
	authInterceptor := grpcsrv.NewAuthInterceptor()
	srv := grpcsrv.NewGrpcServer(
		[]string{
			"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
		},
		authInterceptor.AuthAndIdentifyTickerFunc,
	)

	ipfsfetcherpb.RegisterMessageServer(srv, message.NewServer())

	a.manager.AddWorker(grpcsrv.NewGrpcServerWorker("API", srv, a.cfg.InternalAPI.Bind))

	return nil
}

func (a *Application) initPrometheusWorker() error {
	srv := prometheus.NewServer(a.cfg.Prometheus.Listen, "/metrics")
	a.manager.AddWorker(process.NewServerWorker("prometheus", srv))

	return nil
}

func (a *Application) initHealthWorker() error {
	srv := health.NewHealthCheckServer(a.cfg.Health.Listen, "/status", health.DefaultHandler(a.manager))
	a.manager.AddWorker(process.NewServerWorker("health", srv))

	return nil
}

func (a *Application) registerShutdown() {
	go func(manager *process.Manager) {
		<-a.sigChan

		manager.StopAll()
	}(a.manager)

	a.manager.AwaitAll()
}
