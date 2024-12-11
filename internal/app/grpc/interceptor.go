package grpc

import (
	"context"
	"github.com/korol8484/shortener/internal/app/usecase"
	"github.com/korol8484/shortener/internal/app/user/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"slices"
	"strings"
)

// JwtInterceptor - add or read auth info from request
func JwtInterceptor(jwt *usecase.Jwt) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "could not load metadata")
		}

		if len(md.Get(jwt.GetTokenName())) > 0 {
			tokenFromMD := md.Get(jwt.GetTokenName())[0]

			cl, err := jwt.LoadClaims(tokenFromMD)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "not valid token: %s", err)
			}

			return handler(util.SetUserIDToCtx(ctx, cl.UserID), req)
		}

		user, token, err := jwt.CreateNewToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "can't create token: %s", err)
		}

		md.Set(jwt.GetTokenName(), token)
		ctx = util.SetUserIDToCtx(ctx, user.ID)

		if err = grpc.SetHeader(ctx, md); err != nil {
			return nil, err
		}

		return handler(metadata.NewOutgoingContext(ctx, md), req)
	}
}

// IPInterceptor - check access by ip for methods
func IPInterceptor(CIDR string, methods []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !slices.Contains(methods, info.FullMethod) {
			return handler(ctx, req)
		}

		if CIDR == "" {
			return nil, status.Errorf(codes.Unauthenticated, "cidr not set in config")
		}

		ip := net.ParseIP(getIP(ctx))
		if ip == nil {
			return nil, status.Errorf(codes.Unauthenticated, "x-real-ip not set in request")
		}

		_, ipNet, err := net.ParseCIDR(CIDR)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "can't parse cidr: %s", err)
		}

		if !ipNet.Contains(ip) {
			return nil, status.Errorf(codes.Unauthenticated, "x-real-ip not in cidr")
		}

		return handler(ctx, req)
	}
}

func getIP(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-real-ip")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIP := ips[0]
				return clientIP
			}
		}
	}

	return ""
}
