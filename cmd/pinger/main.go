package main

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bruli/pinger/internal/app"
	"github.com/bruli/pinger/internal/domain"
	infrahttp "github.com/bruli/pinger/internal/infra/http"
	"github.com/bruli/pinger/internal/infra/icmp"
	"github.com/bruli/pinger/internal/infra/yaml"
	"github.com/rs/zerolog"
)

//go:embed checks.yml
var checks []byte

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	log := runWorkers(ctx)

	srv := infrahttp.NewServer(log)

	go shutDown(ctx, log, srv)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("error while starting server")
	}
}

func shutDown(ctx context.Context, log *zerolog.Logger, srv *http.Server) {
	<-ctx.Done()
	log.Info().Msg("shutdown signal received, stopping server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("error shutting down HTTP server")
	}
}

func runWorkers(ctx context.Context) *zerolog.Logger {
	log := buildLogger()
	repo := yaml.NewPingsRepository(checks)
	exec := icmp.NewPingExecutor()

	chLogMdw := app.NewLogCommandHandlerMiddleware(log)
	qhLogMds := app.NewLogQueryHandlerMiddleware(log)

	findQh := qhLogMds(app.NewFindPings(repo))
	execCh := chLogMdw(app.NewExecutePings(exec))

	result, err := findQh.Handle(ctx, app.FindPingsQuery{})
	if err != nil {
		log.Err(err).Msg("error while finding pings")
		os.Exit(1)
	}
	pings, ok := result.([]*domain.Ping)
	if !ok {
		log.Err(err).Msg("error while casting result to pings")
		os.Exit(1)
	}
	log.Info().Msgf("found %d pings", len(pings))

	for _, pi := range pings {
		go execute(ctx, execCh, log, pi)
	}
	return log
}

func buildLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}

func execute(ctx context.Context, ch app.CommandHandler, log *zerolog.Logger, p *domain.Ping) {
	t := time.NewTicker(p.Interval())
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("worker for resource %q stopped", p.ResourceName())
			return
		case <-t.C:
			_, _ = ch.Handle(ctx, app.ExecutePingsCommand{Ping: p})

		}
	}
}
