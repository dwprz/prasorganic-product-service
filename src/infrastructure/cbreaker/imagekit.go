package cbreaker

import (
	"errors"
	"time"

	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/imagekit-developer/imagekit-go/api"
	"github.com/sony/gobreaker/v2"
)

var ImageKit *gobreaker.CircuitBreaker[any]

func init() {
	settings := gobreaker.Settings{
		Name:        "imagekit-restful",
		MaxRequests: 3,
		Interval:    1 * time.Minute,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {

			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.8
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}

			if errors.Is(err, api.ErrNotFound) {
				return true
			}

			return false
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Logger.Infof("circuit breaker %v from status %v to %v", name, from, to)
		},
	}

	ImageKit = gobreaker.NewCircuitBreaker[any](settings)
}
