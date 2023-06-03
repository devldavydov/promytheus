package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/devldavydov/promytheus/internal/common/hash"
	"github.com/devldavydov/promytheus/internal/common/metric"
	pb "github.com/devldavydov/promytheus/internal/grpc"
	"github.com/devldavydov/promytheus/internal/grpc/interceptor"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	_ "google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedMetricServiceServer
	storage storage.Storage
	hmacKey *string
	logger  *logrus.Logger
}

func NewServer(stg storage.Storage, hmacKey *string, trustedSubnet *net.IPNet, logger *logrus.Logger) (*grpc.Server, *Server) {
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.NewTrustedSubnetInterceptor(trustedSubnet, []string{"/grpc.MetricService/UpdateMetrics"}).Handle),
	)
	srv := &Server{storage: stg, hmacKey: hmacKey, logger: logger}
	pb.RegisterMetricServiceServer(grpcSrv, srv)
	return grpcSrv, srv
}

func (s *Server) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.EmptyResponse, error) {
	metrics, err := s.parseUpdateRequest(in.Metrics)
	if err != nil {
		return nil, getErrorStatus(err)
	}

	if err = s.storage.SetMetrics(metrics); err != nil {
		return nil, getErrorStatus(err)
	}
	return &pb.EmptyResponse{}, nil
}

func (s *Server) GetAllMetrics(ctx context.Context, in *pb.EmptyRequest) (*pb.GetAllMetricsResponse, error) {
	metrics, err := s.storage.GetAllMetrics()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
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
		if s.hmacKey != nil {
			res.Hash = item.Value.Hmac(item.MetricName, *s.hmacKey)
		}

		resMetrics = append(resMetrics, res)
	}

	return &pb.GetAllMetricsResponse{Metrics: resMetrics}, nil
}

func (s *Server) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	if in.Id == "" {
		return nil, getErrorStatus(metric.ErrEmptyMetricName)
	}

	resp := &pb.GetMetricResponse{Metric: &pb.Metric{Id: in.Id}}
	var val metric.MetricValue
	if in.Type == pb.MetricType_COUNTER {
		cnt, err := s.storage.GetCounterMetric(in.Id)
		if err != nil {
			return nil, getErrorStatus(err)
		}
		resp.Metric.Type = pb.MetricType_COUNTER
		resp.Metric.Delta = *cnt.IntP()
		val = cnt
	} else if in.Type == pb.MetricType_GAUGE {
		gg, err := s.storage.GetGaugeMetric(in.Id)
		if err != nil {
			return nil, getErrorStatus(err)
		}
		resp.Metric.Type = pb.MetricType_GAUGE
		resp.Metric.Value = *gg.FloatP()
		val = gg
	} else {
		return nil, getErrorStatus(metric.ErrUnknownMetricType)
	}

	if s.hmacKey != nil {
		resp.Metric.Hash = val.Hmac(in.Id, *s.hmacKey)
	}

	return resp, nil
}

func (s *Server) Ping(ctx context.Context, in *pb.EmptyRequest) (*pb.EmptyResponse, error) {
	if !s.storage.Ping() {
		return nil, status.Errorf(codes.Internal, "ping failed")
	}
	return &pb.EmptyResponse{}, nil
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

func getErrorStatus(err error) error {
	if err == nil {
		return nil
	}

	var code codes.Code
	var msg string

	switch {
	case errors.Is(err, metric.ErrUnknownMetricType):
		code, msg = codes.Unimplemented, err.Error()
	case errors.Is(err, metric.ErrMetricHashCheck):
		code, msg = codes.InvalidArgument, err.Error()
	case errors.Is(err, metric.ErrEmptyMetricName):
		code, msg = codes.InvalidArgument, err.Error()
	case errors.Is(err, metric.ErrWrongMetricValue):
		code, msg = codes.InvalidArgument, err.Error()
	case errors.Is(err, storage.ErrMetricNotFound):
		code, msg = codes.NotFound, err.Error()
	default:
		code, msg = codes.NotFound, "internal error"
	}

	return status.Error(code, msg)
}
