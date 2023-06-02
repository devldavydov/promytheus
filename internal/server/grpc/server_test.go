package grpc

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/devldavydov/promytheus/internal/common/metric"
	pb "github.com/devldavydov/promytheus/internal/grpc"
	"github.com/devldavydov/promytheus/internal/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
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
		name     string
		respCode codes.Code
	}{
		{name: "success", respCode: codes.OK},
	} {
		tt := tt
		gs.Run(tt.name, func() {
			gs.createTestServer(nil)
			resp, err := gs.testSrv.Ping(ctx, &pb.EmptyRequest{})

			gs.NoError(err)
			gs.Equal(tt.respCode, codes.Code(resp.Result))
		})
	}

}

func (gs *GrpcServerSuite) TestUpdateMetrics() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tt := range []struct {
		name         string
		req          *pb.UpdateMetricsRequest
		hmacKey      *string
		stgInitFunc  func()
		stgCheckFunc func()
		respCode     codes.Code
		respErr      error
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
	} {
		tt := tt
		gs.Run(tt.name, func() {
			gs.createTestServer(tt.hmacKey)

			if tt.stgInitFunc != nil {
				tt.stgInitFunc()
			}

			resp, err := gs.testSrv.UpdateMetrics(ctx, tt.req)
			if tt.respErr != nil {
				gs.Equal(tt.respErr.Error(), err.Error())
			} else {
				gs.NoError(err)
			}

			gs.Equal(tt.respCode, codes.Code(resp.Result))

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
		gs.createTestServer(nil)
		resp, err := gs.testSrv.GetAllMetrics(ctx, &pb.EmptyRequest{})
		gs.NoError(err)
		gs.Equal(codes.OK, codes.Code(resp.Result))
		gs.Equal(0, len(resp.Metrics))
	})

	gs.Run("get from storage", func() {
		gs.createTestServer(nil)

		gs.stg.SetCounterMetric("counter", metric.Counter(123))
		gs.stg.SetGaugeMetric("gauge", metric.Gauge(123.123))

		resp, err := gs.testSrv.GetAllMetrics(ctx, &pb.EmptyRequest{})
		gs.NoError(err)
		gs.Equal(codes.OK, codes.Code(resp.Result))
		gs.Equal(2, len(resp.Metrics))

		gs.Equal(pb.MetricType_COUNTER, resp.Metrics[0].Type)
		gs.Equal(int64(123), resp.Metrics[0].Delta)

		gs.Equal(pb.MetricType_GAUGE, resp.Metrics[1].Type)
		gs.Equal(float64(123.123), resp.Metrics[1].Value)
	})
}

func (gs *GrpcServerSuite) createTestServer(hmacKey *string) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	gs.testSrv = &Server{storage: gs.stg, hmacKey: hmacKey, logger: gs.logger}
	srv := grpc.NewServer()
	pb.RegisterMetricServiceServer(srv, gs.testSrv)
	go func() {
		srv.Serve(lis)
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(gs.T(), err)

	gs.fTeardown = func() {
		lis.Close()
		srv.Stop()
	}

	gs.testClt = pb.NewMetricServiceClient(conn)
}

func strPointer(s string) *string { return &s }

func TestGrpcServerSuite(t *testing.T) {
	suite.Run(t, new(GrpcServerSuite))
}
