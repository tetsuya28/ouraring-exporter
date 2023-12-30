package repository

import (
	"context"
	"time"

	"github.com/tetsuya28/ouraring-exporter/ouraring"
)

type Repository interface {
	// HeartRate
	GetHeartRateCollector() *HeartRateCollector
	StartHeartRateCollector(ctx context.Context, durationSeconds time.Duration) error

	// PersonalInfo
	GetPersonalInfoCollector() *PersonalInfoCollector
	StartPersonalInfoCollector(ctx context.Context) error
}

type repository struct {
	ouraringClient        *ouraring.Client
	heartRateCollector    *HeartRateCollector
	personalInfoCollector *PersonalInfoCollector
}

func New(ouraringClient *ouraring.Client) Repository {
	return &repository{
		ouraringClient:        ouraringClient,
		heartRateCollector:    newHeartRateCollector(),
		personalInfoCollector: newPersonalInfoCollector(),
	}
}
