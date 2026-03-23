package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bruli/pinger/internal/app"
	"github.com/bruli/pinger/internal/config"
	"github.com/bruli/pinger/internal/domain"
	infrahttp "github.com/bruli/pinger/internal/infra/http"
	"github.com/bruli/pinger/internal/infra/icmp"
	"github.com/bruli/pinger/internal/infra/listeners"
	"github.com/bruli/pinger/internal/infra/nats"
	"github.com/bruli/pinger/internal/infra/yaml"
)

//go:embed checks.yml
var checks []byte

const serviceName = "pinger"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	log := runWorkers(ctx)

	srv := infrahttp.NewServer(log)

	go shutDown(ctx, log, srv)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.ErrorContext(ctx, "error while starting server",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func shutDown(ctx context.Context, log *slog.Logger, srv *http.Server) {
	<-ctx.Done()
	log.InfoContext(ctx, "received shutdown signal")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.ErrorContext(ctx, "error shutting down server",
			slog.String("error", err.Error()))
	}
}

func runWorkers(ctx context.Context) *slog.Logger {
	log := buildLog()
	conf, err := config.New()
	if err != nil {
		log.ErrorContext(ctx, "error while loading configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	repo := yaml.NewPingsRepository(checks)
	exec := icmp.NewPingExecutor()
	publish, err := nats.NewPublisher(conf.NatsServerURL)
	if err != nil {
		log.ErrorContext(ctx, "error while creating nats publisher", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer publish.Close()

	publishListener := listeners.NewPublishOnPingResult(publish)

	evBus := app.NewEventBus()
	evBus.Subscribe(domain.PingEvent{
		BasicEvent: &domain.BasicEvent{Name: domain.PingResultEventName},
	}, publishListener)

	chLogMdw := app.NewLogCommandHandlerMiddleware(log)
	chBusMdw := app.NewEventBusCommandHandlerMiddleware(evBus, log)
	chMultiMdw := app.NewCommandHandlerMultiMiddleware(chBusMdw, chLogMdw)
	qhLogMds := app.NewLogQueryHandlerMiddleware(log)

	findQh := qhLogMds(app.NewFindPings(repo))
	execCh := chMultiMdw(app.NewExecutePings(exec))

	result, err := findQh.Handle(ctx, app.FindPingsQuery{})
	if err != nil {
		log.ErrorContext(ctx, "error while finding pings", slog.String("error", err.Error()))
		os.Exit(1)
	}
	pings, ok := result.([]*domain.Ping)
	if !ok {
		log.ErrorContext(ctx, "error while casting result to pings", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.InfoContext(ctx, fmt.Sprintf("found %d pings", len(pings)))

	for _, pi := range pings {
		go execute(ctx, execCh, log, pi)
	}
	return log
}

func buildLog() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log := slog.New(handler)
	log.With("service", serviceName)
	return log
}

func execute(ctx context.Context, ch app.CommandHandler, log *slog.Logger, p *domain.Ping) {
	t := time.NewTicker(p.Interval())
	for {
		select {
		case <-ctx.Done():
			log.InfoContext(ctx, fmt.Sprintf("worker for resource %q stopped", p.ResourceName()))
			return
		case <-t.C:
			_, _ = ch.Handle(ctx, app.ExecutePingsCommand{Ping: p})

		}
	}
}
