package publisher

import (
	"context"
	"fmt"
	"time"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/common/nettools"
	pb "github.com/devldavydov/promytheus/internal/grpc"
	"github.com/devldavydov/promytheus/internal/grpc/interceptor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GRPCPublisher is a gRPC metric publisher.
type GRPCPublisher struct {
	serverAddress        nettools.Address
	hmacKey              *string
	metricsChan          <-chan metric.Metrics
	logger               *logrus.Logger
	failedCounterMetrics metric.Metrics
	hostIP               string
	threadID             int
	shutdownTimeout      time.Duration
	tlsCredentials       credentials.TransportCredentials
}

// GRPCPublisher constructor.
func NewGRPCPublisher(
	serverAddress nettools.Address,
	metricsChan <-chan metric.Metrics,
	threadID int,
	logger *logrus.Logger,
	extra PublisherExtraSettings,
) *GRPCPublisher {
	shutdownTimeout := _defaultShutdownTimeout
	if extra.ShutdownTimeout != nil {
		shutdownTimeout = *extra.ShutdownTimeout
	}

	return &GRPCPublisher{
		serverAddress:   serverAddress,
		hmacKey:         extra.HmacKey,
		metricsChan:     metricsChan,
		threadID:        threadID,
		shutdownTimeout: shutdownTimeout,
		hostIP:          extra.HostIP.String(),
		tlsCredentials:  extra.EncrSettings.TLSCredentials,
		logger:          logger,
	}
}

func (g *GRPCPublisher) Publish() {
	for metricsToSend := range g.metricsChan {
		g.processMetrics([]metric.Metrics{metricsToSend, g.failedCounterMetrics})
	}
	// If channel closed, try to send failed metrics and exit
	g.shutdown()
	g.logger.Infof("gRPC publisher[%d] thread shutdown due to context closed", g.threadID)
}

func (g *GRPCPublisher) processMetrics(metricsList []metric.Metrics) {
	var counterMetricsToSend = make(metric.Metrics)

	g.logger.Debugf("gRPC publisher[%d] publishing metrics: %+v", g.threadID, metricsList)

	metricReq := make([]metric.MetricsDTO, 0, totalMetrics(metricsList))

	iterateMetrics(metricsList, func(name string, value metric.MetricValue) {
		metricReq = append(metricReq, prepareMetric(name, value, g.hmacKey))

		if value.TypeName() == metric.CounterTypeName {
			curVal, ok := counterMetricsToSend[name]
			if !ok {
				counterMetricsToSend[name] = value
			} else {
				counterMetricsToSend[name] = curVal.(metric.Counter) + value.(metric.Counter)
			}
		}
	})

	if err := g.publishMetrics(metricReq); err != nil {
		g.logger.Errorf("gRPC publisher[%d] failed to publish: %v", g.threadID, err)
		g.failedCounterMetrics = counterMetricsToSend
		return
	}

	g.failedCounterMetrics = nil
}

func (g *GRPCPublisher) publishMetrics(metricReq []metric.MetricsDTO) error {
	opts := []grpc.DialOption{grpc.WithUnaryInterceptor(interceptor.NewGzipClientInterceptor().Handle)}
	if g.tlsCredentials != nil {
		opts = append([]grpc.DialOption{grpc.WithTransportCredentials(g.tlsCredentials)}, opts...)
	}

	conn, err := grpc.Dial(g.serverAddress.String(), opts...)
	if err != nil {
		return fmt.Errorf("gRPC publisher[%d] failed to create connection: %w", g.threadID, err)
	}
	defer conn.Close()

	clnt := pb.NewMetricServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), _defaultRequestTimeout)
	defer cancel()

	updMetrics := make([]*pb.Metric, 0, len(metricReq))
	for _, mReq := range metricReq {
		updMetric := &pb.Metric{Id: mReq.ID}
		if mReq.Hash != nil {
			updMetric.Hash = *mReq.Hash
		}

		if mReq.MType == metric.CounterTypeName {
			updMetric.Type = pb.MetricType_COUNTER
			updMetric.Delta = *mReq.Delta
		} else if mReq.MType == metric.GaugeTypeName {
			updMetric.Type = pb.MetricType_GAUGE
			updMetric.Value = *mReq.Value
		}

		updMetrics = append(updMetrics, updMetric)
	}

	_, err = clnt.UpdateMetrics(ctx, &pb.UpdateMetricsRequest{Metrics: updMetrics})
	if err != nil {
		return fmt.Errorf("gRPC publisher[%d] failed to publish metrics: %w", g.threadID, err)
	}

	return nil
}

func (g *GRPCPublisher) shutdown() {
	if g.failedCounterMetrics == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.shutdownTimeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if g.failedCounterMetrics == nil {
				return
			}

			g.processMetrics([]metric.Metrics{g.failedCounterMetrics})
		case <-ctx.Done():
			return
		}
	}
}
