package core_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.With(Registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: "shelfmate_http_requests_total",
			Help: "Количество HTTP-запросов по методу, шаблону маршрута и статусу.",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDuration = promauto.With(Registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shelfmate_http_request_duration_seconds",
			Help:    "Время обработки HTTP-запроса.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)
