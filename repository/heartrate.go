package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tetsuya28/ouraring-exporter/ouraring"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type HeartRateCollector struct {
	desc      *prometheus.Desc
	heartRate []ouraring.HeartRateModel
}

func newHeartRateCollector() *HeartRateCollector {
	return &HeartRateCollector{
		desc: prometheus.NewDesc(
			"ouraring_exporter_usercollection_heartrate",
			"",
			[]string{
				"source",
			},
			nil,
		),
	}
}

func (c *HeartRateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *HeartRateCollector) Collect(ch chan<- prometheus.Metric) {
	for _, h := range c.heartRate {
		ch <- prometheus.NewMetricWithTimestamp(
			h.Timestamp,
			prometheus.MustNewConstMetric(
				c.desc,
				prometheus.GaugeValue,
				float64(h.Bpm),
				string(h.Source),
			),
		)
	}
}

func (r repository) GetHeartRateCollector() *HeartRateCollector {
	return r.heartRateCollector
}

func (r repository) StartHeartRateCollector(ctx context.Context, durationSeconds time.Duration) error {
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

			r.heartRateCollector.heartRate = b.JSON200.Data
		}
	}()

	return nil
}
