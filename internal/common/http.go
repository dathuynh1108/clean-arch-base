package common

import "context"

type HTTPDelivery interface {
	ServeHTTP(host, port string) error
}

type StartStopper interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
