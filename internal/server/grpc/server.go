package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/devldavydov/promytheus/internal/common/hash"
	"github.com/devldavydov/promytheus/internal/common/metric"
	pb "github.com/devldavydov/promytheus/internal/grpc"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Server struct {
	pb.UnimplementedMetricServiceServer
	storage storage.Storage
	hmacKey *string
	logger  *logrus.Logger
}

func NewServer(stg storage.Storage, hmacKey *string, logger *logrus.Logger) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterMetricServiceServer(srv, &Server{storage: stg, hmacKey: hmacKey, logger: logger})
	return srv
}

func (s *Server) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	resp := &pb.UpdateMetricsResponse{}

	metrics, err := s.parseUpdateRequest(in.Metrics)
	if err != nil {
		resp.Result = getResult(err)
		return resp, err
	}

	err = s.storage.SetMetrics(metrics)
	resp.Result = getResult(err)
	return resp, err
}

func (s *Server) GetAllMetrics(ctx context.Context, in *pb.EmptyRequest) (*pb.GetAllMetricsResponse, error) {
	resp := &pb.GetAllMetricsResponse{}

	metrics, err := s.storage.GetAllMetrics()
	if err != nil {
		resp.Result = getResult(err)
		return resp, err
	}

	resMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, item := range metrics {
		res := &pb.Metric{Id: item.MetricName}

		switch item.Value.TypeName() {
		case metric.CounterTypeName:
			res.Type = pb.MetricType_COUNTER
			res.Delta = *item.Value.(metric.Counter).IntP()
		case metric.GaugeTypeName:
			res.Type = pb.MetricType_GAUGE
			res.Value = *item.Value.(metric.Gauge).FloatP()
		}

		resMetrics = append(resMetrics, res)
	}

	resp.Result = getResult(nil)
	resp.Metrics = resMetrics

	return resp, nil
}

func (s *Server) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.PingResponse, error) {
	resp := &pb.PingResponse{Result: getResult(nil)}
	if !s.storage.Ping() {
		resp.Result = uint32(codes.Internal)
	}
	return resp, nil
}

func (s *Server) parseUpdateRequest(inMetrics []*pb.Metric) ([]storage.StorageItem, error) {
	metrics := make([]storage.StorageItem, 0, len(inMetrics))
	for _, inMetric := range inMetrics {
		if inMetric.Id == "" {
			return nil, metric.ErrEmptyMetricName
		}

		stMetric := storage.StorageItem{MetricName: inMetric.Id}
		if inMetric.Type == pb.MetricType_COUNTER {
			val, err := metric.NewCounterFromIntP(&inMetric.Delta)
			if err != nil {
				return nil, fmt.Errorf("incorrect %s '%s': %w", metric.CounterTypeName, inMetric.Id, metric.ErrWrongMetricValue)
			}
			stMetric.Value = val

			if err = s.hmacCheck(inMetric.Hash, inMetric.Id, val); err != nil {
				return nil, fmt.Errorf("incorrect %s '%s': %w", metric.CounterTypeName, inMetric.Id, err)
			}

		} else if inMetric.Type == pb.MetricType_GAUGE {
			val, err := metric.NewGaugeFromFloatP(&inMetric.Value)
			if err != nil {
				return nil, fmt.Errorf("incorrect %s '%s': %w", metric.GaugeTypeName, inMetric.Id, metric.ErrWrongMetricValue)
			}
			stMetric.Value = val

			if err = s.hmacCheck(inMetric.Hash, inMetric.Id, val); err != nil {
				return nil, fmt.Errorf("incorrect %s '%s': %w", metric.GaugeTypeName, inMetric.Id, err)
			}
		} else {
			return nil, metric.ErrUnknownMetricType
		}
		metrics = append(metrics, stMetric)
	}

	return metrics, nil
}

func (s *Server) hmacCheck(reqHash string, reqID string, value metric.MetricValue) error {
	if s.hmacKey == nil {
		return nil
	}

	if !hash.HmacEqual(reqHash, value.Hmac(reqID, *s.hmacKey)) {
		return metric.ErrMetricHashCheck
	}
	return nil
}

func getResult(err error) uint32 {
	if err == nil {
		return uint32(codes.OK)
	}

	var res codes.Code

	if errors.Is(err, metric.ErrUnknownMetricType) {
		res = codes.Unimplemented
	} else if errors.Is(err, metric.ErrMetricHashCheck) ||
		errors.Is(err, metric.ErrEmptyMetricName) ||
		errors.Is(err, metric.ErrWrongMetricValue) {
		res = codes.InvalidArgument
	} else {
		res = codes.Internal
	}

	return uint32(res)
}
