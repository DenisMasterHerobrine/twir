package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/queue"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"google.golang.org/grpc"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/grpc_impl"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	"go.uber.org/zap"
)

func main() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())

	config, err := cfg.New()
	if err != nil || config == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	if config.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	if config.SentryDsn != "" {
		sentry.Init(
			sentry.ClientOptions{
				Dsn:              config.SentryDsn,
				Environment:      config.AppEnv,
				Debug:            false,
				TracesSampleRate: 1.0,
			},
		)
	}

	var logger *zap.Logger

	if config.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

	// gorm
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(10)
	d.SetConnMaxIdleTime(1 * time.Minute)
	defer d.Close()

	// sqlx
	dbConnOpts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("cannot parse postgres url connection: %w", err))
	}
	pgConn, err := sqlx.ConnectContext(appCtx, "postgres", dbConnOpts)
	defer pgConn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// redis
	url, err := redis.ParseURL(config.RedisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     url.Addr,
			Password: url.Password,
			DB:       url.DB,
			Username: url.Username,
		},
	)
	defer redisClient.Close()

	redisClient.Conn()

	tokensGrpc := clients.NewTokens(config.AppEnv)

	queueDistributor := queue.NewRedisTaskDistributor(config, logger)
	queueProcessor := queue.NewRedisTaskProcessor(
		queue.RedisTaskProcessorOpts{
			Cfg:        *config,
			Logger:     logger,
			Gorm:       db,
			TokensGrpc: tokensGrpc,
		},
	)
	defer queueProcessor.Stop()

	go func() {
		err := queueProcessor.Start()
		if err != nil {
			logger.Fatal("Error starting queue processor", zap.Error(err))
		}
	}()

	s := &services.Services{
		Config: config,
		Logger: logger,
		Gorm:   db,
		Sqlx:   pgConn,
		Redis:  redisClient,
		GrpcClients: &services.Grpc{
			WebSockets: clients.NewWebsocket(config.AppEnv),
			Bots:       clients.NewBots(config.AppEnv),
			Dota:       clients.NewDota(config.AppEnv),
			Eval:       clients.NewEval(config.AppEnv),
			Tokens:     tokensGrpc,
			Events:     clients.NewEvents(config.AppEnv),
			Ytsr:       clients.NewYtsr(config.AppEnv),
		},
		TaskDistributor: queueDistributor,
	}

	variablesService := variables.New(
		&variables.Opts{
			Services: s,
		},
	)
	commandsService := commands.New(
		&commands.Opts{
			Services:         s,
			VariablesService: variablesService,
		},
	)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.PARSER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()
	parser.RegisterParserServer(
		grpcServer,
		grpc_impl.NewServer(s, commandsService, variablesService),
	)
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	logger.Info("Parser microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting")
	appCtxCancel()
}
