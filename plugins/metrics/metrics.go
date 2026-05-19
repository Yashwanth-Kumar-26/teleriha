package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

var (
	messagesProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "teleriha_messages_processed_total",
		Help: "The total number of processed messages",
	}, []string{"update_type", "status"})

	processingDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "teleriha_processing_duration_seconds",
		Help:    "Histogram of response latency for update handlers.",
		Buckets: prometheus.DefBuckets,
	}, []string{"update_type"})

)

// Server encapsulates the metrics HTTP server
type Server struct {
	port   int
	logger zerolog.Logger
}

// NewServer creates a new metrics server
func NewServer(port int, logger zerolog.Logger) *Server {
	return &Server{
		port:   port,
		logger: logger,
	}
}

// Start launches the embedded Prometheus HTTP endpoint in a goroutine.
func (s *Server) Start() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		addr := fmt.Sprintf(":%d", s.port)
		s.logger.Info().Msgf("Starting Prometheus metrics server on %s/metrics", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			s.logger.Error().Err(err).Msg("Metrics server failed")
		}
	}()
}

// Middleware creates a bot middleware to track Prometheus metrics.
func Middleware() bot.Middleware {
	return func(next bot.Handler) bot.Handler {
		return func(ctx *bot.Context) error {
			start := time.Now()
			
			updateType := "unknown"
			if ctx.Update.Message != nil {
				updateType = "message"
			} else if ctx.Update.CallbackQuery != nil {
				updateType = "callback_query"
			} else if ctx.Update.InlineQuery != nil {
				updateType = "inline_query"
			}

			err := next(ctx)
			
			duration := time.Since(start).Seconds()
			processingDuration.WithLabelValues(updateType).Observe(duration)

			status := "success"
			if err != nil {
				status = "error"
			}
			messagesProcessed.WithLabelValues(updateType, status).Inc()

			return err
		}
	}
}
