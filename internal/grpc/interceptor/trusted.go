package interceptor

import (
	"context"
	"net"

	"github.com/devldavydov/promytheus/internal/common/nettools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TrustedSubnetIncerceptor struct {
	trustedSubnet    *net.IPNet
	protectedMethods map[string]bool
}

func NewTrustedSubnetInterceptor(trustedSubnet *net.IPNet, protectedMethods []string) *TrustedSubnetIncerceptor {
	pm := make(map[string]bool, len(protectedMethods))
	for _, p := range protectedMethods {
		pm[p] = true
	}
	return &TrustedSubnetIncerceptor{trustedSubnet: trustedSubnet, protectedMethods: pm}
}

func (t *TrustedSubnetIncerceptor) Handle(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if t.trustedSubnet == nil || !t.protectedMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	vals := md.Get(nettools.RealIPHeader)
	if len(vals) == 0 || !t.trustedSubnet.Contains(net.ParseIP(vals[0])) {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	return handler(ctx, req)
}
