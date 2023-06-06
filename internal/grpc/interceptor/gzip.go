package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

type GzipClientInterceptor struct {
	gzip grpc.CallOption
}

func NewGzipClientInterceptor() *GzipClientInterceptor {
	return &GzipClientInterceptor{gzip: grpc.UseCompressor(gzip.Name)}
}

func (g *GzipClientInterceptor) Handle(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	opts = append([]grpc.CallOption{g.gzip}, opts...)
	return invoker(ctx, method, req, reply, cc, opts...)
}
