package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// LoggingInterceptor logs incoming gRPC requests with duration and method.
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("[gRPC] %s | error | %v | %s\n", info.FullMethod, err, duration)
	} else {
		fmt.Printf("[gRPC] %s | success | %s\n", info.FullMethod, duration)
	}
	return resp, err
}
