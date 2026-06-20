package middleware

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor est un interceptor gRPC qui logue chaque appel.
func UnaryLoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	log.Printf("[gRPC] --> %s", info.FullMethod)

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("[gRPC] <-- %s | code=%s | duration=%v | msg=%s",
			info.FullMethod, st.Code(), duration, st.Message())
	} else {
		log.Printf("[gRPC] <-- %s | OK | duration=%v", info.FullMethod, duration)
	}

	return resp, err
}

// UnaryPanicRecoveryInterceptor rattrape les panics dans les handlers gRPC.
func UnaryPanicRecoveryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[gRPC] PANIC in %s: %v", info.FullMethod, r)
			err = grpc.Errorf(grpc.Code(err), "internal server error")
		}
	}()
	return handler(ctx, req)
}
