package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tetsuya28/ouraring-exporter/ouraring"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

type PersonalInfoCollector struct {
	desc         *prometheus.Desc
	personalInfo *ouraring.PersonalInfoResponse
}

func newPersonalInfoCollector() *PersonalInfoCollector {
	return &PersonalInfoCollector{
		desc: prometheus.NewDesc(
			"ouraring_exporter_usercollection_personal_info",
			"",
			[]string{
				"age",
				"weight",
				"height",
				"biological_sex",
				"email",
			},
			nil,
		),
	}
}

func (c *PersonalInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.desc
}

func (c *PersonalInfoCollector) Collect(ch chan<- prometheus.Metric) {
	metric, err := prometheus.NewConstMetric(
		c.desc,
		prometheus.GaugeValue,
		float64(1),
		strconv.Itoa(*c.personalInfo.Age),
		strconv.Itoa(int(*c.personalInfo.Weight)),
		strconv.Itoa(int(*c.personalInfo.Height)),
		*c.personalInfo.BiologicalSex,
		*c.personalInfo.Email,
	)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	ch <- metric
}

func (r repository) GetPersonalInfoCollector() *PersonalInfoCollector {
	return r.personalInfoCollector
}

func (r repository) StartPersonalInfoCollector(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)

	go func() {
		for ; true; <-ticker.C {
			resp, err := r.ouraringClient.SinglePersonalInfoDocumentV2UsercollectionPersonalInfoGet(ctx)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			info, err := ouraring.ParseSinglePersonalInfoDocumentV2UsercollectionPersonalInfoGetResponse(resp)
			if err != nil {
				log.Warn(err.Error())
				continue
			}
			r.personalInfoCollector.personalInfo = info.JSON200
		}
	}()

	return nil
}
