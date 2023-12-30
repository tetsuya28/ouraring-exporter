package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tetsuya28/ouraring-exporter/ouraring"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type Repository interface {
	GetHeartRateRegister() *Collector
	StartHeartRate(ctx context.Context, durationSeconds time.Duration) error
}

type repository struct {
	ouraringClient    *ouraring.Client
	heartrateRegister *Collector
}

type Collector struct {
	metric    *prometheus.Desc
	heartRate []ouraring.HeartRateModel
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.metric
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, h := range c.heartRate {
		ch <- prometheus.NewMetricWithTimestamp(
			h.Timestamp,
			prometheus.MustNewConstMetric(
				c.metric,
				prometheus.GaugeValue,
				float64(h.Bpm),
				string(h.Source),
			),
		)
	}
}

func NewHeartRate(ouraringClient *ouraring.Client) Repository {
	heartrate := &Collector{
		metric: prometheus.NewDesc(
			"ouraring_exporter_usercollection_heartrate",
			"",
			[]string{
				"source",
			},
			nil,
		),
	}

	return &repository{
		ouraringClient:    ouraringClient,
		heartrateRegister: heartrate,
	}
}

func (r repository) GetHeartRateRegister() *Collector {
	return r.heartrateRegister
}

func (r repository) StartHeartRate(ctx context.Context, durationSeconds time.Duration) error {
	ticker := time.NewTicker(durationSeconds)

	go func() {
		for ; true; <-ticker.C {
			end := time.Now()
			start := end.Add(-1 * 24 * time.Hour)

			resp, err := r.ouraringClient.MultipleHeartRateDocumentsV2UsercollectionHeartrateGet(
				ctx,
				&ouraring.MultipleHeartRateDocumentsV2UsercollectionHeartrateGetParams{
					StartDatetime: &start,
					EndDatetime:   &end,
				},
			)
			if err != nil {
				log.Warn(err.Error())
				continue
			}

			b, err := ouraring.ParseMultipleHeartRateDocumentsV2UsercollectionHeartrateGetResponse(resp)
			if err != nil {
				log.Warn(err.Error())
				continue
			}

			if b.StatusCode() != http.StatusOK {
				log.Warn(err.Error())
				continue
			}

			r.heartrateRegister.heartRate = b.JSON200.Data
		}
	}()

	return nil
}
