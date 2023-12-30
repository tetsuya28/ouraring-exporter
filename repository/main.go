package repository

import (
	"context"
	"time"

	"github.com/tetsuya28/ouraring-exporter/ouraring"
)

type Repository interface {
	// HeartRate
	GetHeartRateRegister() *HeartRateCollector
	StartHeartRate(ctx context.Context, durationSeconds time.Duration) error
}

type repository struct {
	ouraringClient     *ouraring.Client
	heartRateCollector *HeartRateCollector
}

func New(ouraringClient *ouraring.Client) Repository {

	return &repository{
		ouraringClient:     ouraringClient,
		heartRateCollector: newHeartRateCollector(),
	}
}
