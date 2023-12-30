package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tetsuya28/ouraring-exporter/config"
	"github.com/tetsuya28/ouraring-exporter/ouraring"
	"github.com/tetsuya28/ouraring-exporter/repository"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

var (
	OURARING_API_DOMAIN = "https://api.ouraring.com"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector())

	requestEditor := ouraring.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.OuraringAPIKey))
		return nil
	})

	ouraringClient, err := ouraring.NewClient(OURARING_API_DOMAIN, requestEditor)
	if err != nil {
		log.Critical(err.Error())
	}

	repository := repository.NewHeartRate(ouraringClient)
	repository.StartHeartRate(ctx, time.Duration(cfg.OuraringAPICallIntervalSeconds)*time.Second)
	reg.MustRegister(repository.GetHeartRateRegister())

	mux := http.NewServeMux()

	mux.Handle(
		"/metrics",
		promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Critical(err.Error())
		}
	}()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(tctx); err != nil {
		log.Critical(err.Error())
	}
}
