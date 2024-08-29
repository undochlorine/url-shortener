package shortener

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientDefaultTimeout is unary client interceptor that sets a timeout
// on the request context unless the context has timeout
func defaultTimeout(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if _, ok := ctx.Deadline(); !ok {
		c, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		ctx = c
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

// NewShortenerService Creates client ready to call remote procedures
// Create Client call methods with dot
func NewShortenerService(connString string) (ShortenerClient, func() error, error) {
	retryPolicy := `{
		"methodConfig": [{
			"name": [{"service": "shortener.Shortener"}],
			"waitForReady": true,
			"retryPolicy": {
				"MaxAttempts": 10,
				"InitialBackoff": ".1s",
				"MaxBackoff": "2.0s",
				"BackOffMultiplier": 2.0,
				"RetryableStatusCodes": [ "UNAVAILABLE", "INTERNAL" ]
			}
		}]
	}`
	conn, err := grpc.NewClient(connString, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(defaultTimeout), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		return nil, nil, err
	}

	cl := NewShortenerClient(conn)
	return cl, conn.Close, nil
}
