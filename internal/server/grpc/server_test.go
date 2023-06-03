package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/metric"
	"github.com/devldavydov/promytheus/internal/common/nettools"
	pb "github.com/devldavydov/promytheus/internal/grpc"
	"github.com/devldavydov/promytheus/internal/grpc/interceptor"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type GrpcServerSuite struct {
	suite.Suite
	testSrv   *Server
	testClt   pb.MetricServiceClient
	stg       storage.Storage
	logger    *logrus.Logger
	fTeardown func()
}

func (gs *GrpcServerSuite) SetupSuite() {
	gs.logger = logrus.New()
}

func (gs *GrpcServerSuite) SetupSubTest() {
	var err error
	gs.stg, err = storage.NewMemStorage(context.TODO(), gs.logger, storage.NewPersistSettings(0, "", false))
	require.NoError(gs.T(), err)
}

func (gs *GrpcServerSuite) TearDownSubTest() {
	gs.fTeardown()
}

func (gs *GrpcServerSuite) TestPing() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tt := range []struct {
		name          string
		trustedSubnet *net.IPNet
	}{
		{name: "success"},
		{name: "ignore subnet check", trustedSubnet: getSubnet("10.0.0.0/16")},
	} {
		tt := tt
		gs.Run(tt.name, func() {
			gs.createTestServer(nil, tt.trustedSubnet)
			_, err := gs.testClt.Ping(ctx, &pb.EmptyRequest{})

			gs.NoError(err)
		})
	}
}

func (gs *GrpcServerSuite) TestUpdateMetrics() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tt := range []struct {
		name          string
		req           *pb.UpdateMetricsRequest
		hmacKey       *string
		trustedSubnet *net.IPNet
		cltIP         *string
		stgInitFunc   func()
		stgCheckFunc  func()
		respCode      codes.Code
		respErr       error
	}{
		{
			name:     "unknown metric type",
			respCode: codes.Unimplemented,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_UNKNOWN, Id: "foo", Delta: 123},
				},
			},
			respErr: metric.ErrUnknownMetricType,
		},
		{
			name:     "empty metric name",
			respCode: codes.InvalidArgument,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_COUNTER, Id: "", Delta: 123},
				},
			},
			respErr: metric.ErrEmptyMetricName,
		},
		{
			name:     "invalid counter",
			respCode: codes.InvalidArgument,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_COUNTER, Id: "ttt", Delta: -123},
				},
			},
			respErr: errors.New("incorrect counter 'ttt': wrong metric value"),
		},
		{
			name:     "correct update",
			respCode: codes.OK,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 1},
					{Type: pb.MetricType_GAUGE, Id: "gauge1", Value: 123.123},
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 2},
					{Type: pb.MetricType_COUNTER, Id: "counter2", Delta: 3},
				},
			},
			stgInitFunc: func() {
				gs.stg.SetCounterMetric("counter2", metric.Counter(2))
			},
			stgCheckFunc: func() {
				vC, _ := gs.stg.GetCounterMetric("counter1")
				gs.Equal(metric.Counter(3), vC)

				vC, _ = gs.stg.GetCounterMetric("counter2")
				gs.Equal(metric.Counter(5), vC)

				vG, _ := gs.stg.GetGaugeMetric("gauge1")
				gs.Equal(metric.Gauge(123.123), vG)
			},
		},
		{
			name:     "correct update with hash check",
			respCode: codes.OK,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_GAUGE, Id: "Sys", Value: 13220880, Hash: "48a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{Type: pb.MetricType_COUNTER, Id: "PollCount", Delta: 5, Hash: "b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{Type: pb.MetricType_COUNTER, Id: "PollCount", Delta: 5, Hash: "b9203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
				},
			},
			hmacKey: strPointer("foobar"),
			stgCheckFunc: func() {
				vC, _ := gs.stg.GetCounterMetric("PollCount")
				gs.Equal(metric.Counter(10), vC)

				vG, _ := gs.stg.GetGaugeMetric("Sys")
				gs.Equal(metric.Gauge(13220880), vG)
			},
		},
		{
			name:     "update hash check failed",
			respCode: codes.InvalidArgument,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_GAUGE, Id: "Sys", Value: 13220880, Hash: "58a93e5dde0297029bf66cc10a1cdda9be6f858667ea885dc1b0d810032aa292"},
					{Type: pb.MetricType_COUNTER, Id: "PollCount", Delta: 5, Hash: "b8203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
					{Type: pb.MetricType_COUNTER, Id: "PollCount", Delta: 5, Hash: "b8203cac5904e73da2504aabfb77a419d3d3f9a0baee3707c55070432c6ff5a8"},
				},
			},
			hmacKey: strPointer("foobar"),
			respErr: errors.New("incorrect gauge 'Sys': metric hash check fail"),
		},
		{
			name:     "correct update with subnet check",
			respCode: codes.OK,
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 1},
					{Type: pb.MetricType_GAUGE, Id: "gauge1", Value: 123.123},
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 2},
					{Type: pb.MetricType_COUNTER, Id: "counter2", Delta: 3},
				},
			},
			trustedSubnet: getSubnet("192.168.0.0/16"),
			cltIP:         strPointer("192.168.1.1"),
			stgInitFunc: func() {
				gs.stg.SetCounterMetric("counter2", metric.Counter(2))
			},
			stgCheckFunc: func() {
				vC, _ := gs.stg.GetCounterMetric("counter1")
				gs.Equal(metric.Counter(3), vC)

				vC, _ = gs.stg.GetCounterMetric("counter2")
				gs.Equal(metric.Counter(5), vC)

				vG, _ := gs.stg.GetGaugeMetric("gauge1")
				gs.Equal(metric.Gauge(123.123), vG)
			},
		},
		{
			name: "update failed because of subnet check",
			req: &pb.UpdateMetricsRequest{
				Metrics: []*pb.Metric{
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 1},
					{Type: pb.MetricType_GAUGE, Id: "gauge1", Value: 123.123},
					{Type: pb.MetricType_COUNTER, Id: "counter1", Delta: 2},
					{Type: pb.MetricType_COUNTER, Id: "counter2", Delta: 3},
				},
			},
			trustedSubnet: getSubnet("192.168.0.0/16"),
			cltIP:         strPointer("10.10.1.1"),
			respCode:      codes.PermissionDenied,
			respErr:       errors.New("forbidden"),
		},
	} {
		tt := tt
		gs.Run(tt.name, func() {
			gs.createTestServer(tt.hmacKey, tt.trustedSubnet)

			if tt.stgInitFunc != nil {
				tt.stgInitFunc()
			}

			cltCtx := ctx
			if tt.cltIP != nil {
				md := metadata.New(map[string]string{nettools.RealIPHeader: *tt.cltIP})
				cltCtx = metadata.NewOutgoingContext(ctx, md)
			}

			_, err := gs.testClt.UpdateMetrics(cltCtx, tt.req)
			if tt.respErr != nil {
				respStatus, ok := status.FromError(err)
				gs.True(ok)
				gs.Equal(tt.respCode, respStatus.Code())
				gs.Equal(tt.respErr.Error(), respStatus.Message())
			} else {
				gs.NoError(err)
			}

			if tt.stgCheckFunc != nil {
				tt.stgCheckFunc()
			}
		})
	}
}

func (gs *GrpcServerSuite) TestGetAllMetrics() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gs.Run("empty storage", func() {
		gs.createTestServer(nil, nil)
		resp, err := gs.testClt.GetAllMetrics(ctx, &pb.EmptyRequest{})
		gs.NoError(err)
		gs.Equal(0, len(resp.Metrics))
	})

	gs.Run("get from storage with hash and ignore subnet check", func() {
		gs.createTestServer(strPointer("foobar"), getSubnet("10.0.0.0/16"))

		gs.stg.SetCounterMetric("counter", metric.Counter(123))
		gs.stg.SetGaugeMetric("gauge", metric.Gauge(123.123))

		resp, err := gs.testSrv.GetAllMetrics(ctx, &pb.EmptyRequest{})
		gs.NoError(err)
		gs.Equal(2, len(resp.Metrics))

		gs.Equal(pb.MetricType_COUNTER, resp.Metrics[0].Type)
		gs.Equal(int64(123), resp.Metrics[0].Delta)
		gs.Equal("c80d8c33875ffba1d06c517749b210aaa3ca9aceb8e2019f64626c66f117da3d", resp.Metrics[0].Hash)

		gs.Equal(pb.MetricType_GAUGE, resp.Metrics[1].Type)
		gs.Equal(float64(123.123), resp.Metrics[1].Value)
		gs.Equal("45a63e4085f263e02fd473e1bcc46f563107662a59a9c53678a59f3fc17e8b62", resp.Metrics[1].Hash)
	})
}

func (gs *GrpcServerSuite) TestGetMetric() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tt := range []struct {
		name          string
		req           *pb.GetMetricRequest
		hmacKey       *string
		trustedSubnet *net.IPNet
		stgInitFunc   func()
		resp          *pb.GetMetricResponse
		respCode      codes.Code
		respErr       error
	}{
		{
			name: "get unknown metric",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_UNKNOWN,
				Id:   "foo",
			},
			respCode: codes.Unimplemented,
			respErr:  metric.ErrUnknownMetricType,
		},
		{
			name: "empty metric name",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_COUNTER,
				Id:   "",
			},
			respCode: codes.InvalidArgument,
			respErr:  metric.ErrEmptyMetricName,
		},
		{
			name: "counter not found",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_COUNTER,
				Id:   "counter",
			},
			respCode: codes.NotFound,
			respErr:  storage.ErrMetricNotFound,
		},
		{
			name: "gauge not found",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_GAUGE,
				Id:   "gauge",
			},
			respCode: codes.NotFound,
			respErr:  storage.ErrMetricNotFound,
		},
		{
			name: "get counter",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_COUNTER,
				Id:   "counter",
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Type:  pb.MetricType_COUNTER,
					Id:    "counter",
					Delta: 123,
				},
			},
			stgInitFunc: func() {
				gs.stg.SetCounterMetric("counter", metric.Counter(123))
			},
		},
		{
			name: "get counter with hash",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_COUNTER,
				Id:   "counter",
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Type:  pb.MetricType_COUNTER,
					Id:    "counter",
					Delta: 123,
					Hash:  "c80d8c33875ffba1d06c517749b210aaa3ca9aceb8e2019f64626c66f117da3d",
				},
			},
			hmacKey: strPointer("foobar"),
			stgInitFunc: func() {
				gs.stg.SetCounterMetric("counter", metric.Counter(123))
			},
		},
		{
			name: "get gauge",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_GAUGE,
				Id:   "gauge",
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Type:  pb.MetricType_GAUGE,
					Id:    "gauge",
					Value: 123.123,
				},
			},
			stgInitFunc: func() {
				gs.stg.SetGaugeMetric("gauge", metric.Gauge(123.123))
			},
		},
		{
			name: "get gauge with hash",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_GAUGE,
				Id:   "gauge",
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Type:  pb.MetricType_GAUGE,
					Id:    "gauge",
					Value: 123.123,
					Hash:  "45a63e4085f263e02fd473e1bcc46f563107662a59a9c53678a59f3fc17e8b62",
				},
			},
			hmacKey: strPointer("foobar"),
			stgInitFunc: func() {
				gs.stg.SetGaugeMetric("gauge", metric.Gauge(123.123))
			},
		},
		{
			name: "get gauge ignore subnet check",
			req: &pb.GetMetricRequest{
				Type: pb.MetricType_GAUGE,
				Id:   "gauge",
			},
			resp: &pb.GetMetricResponse{
				Metric: &pb.Metric{
					Type:  pb.MetricType_GAUGE,
					Id:    "gauge",
					Value: 123.123,
				},
			},
			trustedSubnet: getSubnet("10.0.0.0/16"),
			stgInitFunc: func() {
				gs.stg.SetGaugeMetric("gauge", metric.Gauge(123.123))
			},
		},
	} {
		tt := tt
		gs.Run(tt.name, func() {
			gs.createTestServer(tt.hmacKey, tt.trustedSubnet)

			if tt.stgInitFunc != nil {
				tt.stgInitFunc()
			}

			resp, err := gs.testClt.GetMetric(ctx, tt.req)
			if tt.respErr != nil {
				respStatus, ok := status.FromError(err)
				gs.True(ok)
				gs.Equal(tt.respCode, respStatus.Code())
				gs.Equal(tt.respErr.Error(), respStatus.Message())
			} else {
				gs.NoError(err)
				gs.Equal(tt.resp.Metric, resp.Metric)
			}
		})
	}
}

func (gs *GrpcServerSuite) createTestServer(hmacKey *string, trustedSubnet *net.IPNet) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	var grpcSrv *grpc.Server
	grpcSrv, gs.testSrv = NewServer(gs.stg, hmacKey, trustedSubnet, gs.logger)

	go func() {
		grpcSrv.Serve(lis)
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.NewGzipClientInterceptor().Handle,
		))
	require.NoError(gs.T(), err)

	gs.fTeardown = func() {
		lis.Close()
		grpcSrv.Stop()
	}

	gs.testClt = pb.NewMetricServiceClient(conn)
}

func strPointer(s string) *string { return &s }

func getSubnet(s string) *net.IPNet {
	_, subnet, _ := net.ParseCIDR(s)
	return subnet
}

func TestGrpcServerSuite(t *testing.T) {
	suite.Run(t, new(GrpcServerSuite))
}
