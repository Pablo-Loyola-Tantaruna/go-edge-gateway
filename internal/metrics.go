package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gopherguard_requests_total",
			Help: "Número total de peticiones procesadas por el Gateway",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.SimpleHistogram(
			"gopherguard_request_duration_seconds",
			"Latencia de las peticiones en segundos",
			[]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		),
		[]string{"method", "endpoint"},
	)
)
