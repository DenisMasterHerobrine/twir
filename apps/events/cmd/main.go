package main

import (
	"log/slog"

	eventsActivity "github.com/satont/twir/apps/events/internal/activities/events"
	"github.com/satont/twir/apps/events/internal/chat_alerts"
	"github.com/satont/twir/apps/events/internal/gorm"
	"github.com/satont/twir/apps/events/internal/grpc_impl"
	"github.com/satont/twir/apps/events/internal/hydrator"
	"github.com/satont/twir/apps/events/internal/redis"
	"github.com/satont/twir/apps/events/internal/workers"
	"github.com/satont/twir/apps/events/internal/workflows"
	"github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			cfg.NewFx,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: "events"}),
			logger.NewFx(
				logger.Opts{
					Service: "events",
					Level:   slog.LevelDebug,
				},
			),
			func(config cfg.Config) bots.BotsClient {
				return clients.NewBots(config.AppEnv)
			},
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			func(config cfg.Config) websockets.WebsocketClient {
				return clients.NewWebsocket(config.AppEnv)
			},
			gorm.New,
			redis.New,
			hydrator.New,
			eventsActivity.New,
			workflows.NewEventsWorkflow,
			chat_alerts.New,
		),
		fx.Invoke(
			workers.NewEventsWorker,
			grpc_impl.New,
			func(l logger.Logger) {
				l.Info("Events service started")
			},
		),
	).Run()
}
