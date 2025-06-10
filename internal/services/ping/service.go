package ping

import (
	"context"

	pingmodel "payslip-generation-system/internal/entity/ping"
	pingrepo "payslip-generation-system/internal/repositories/ping"
)

//go:generate mockgen -source=service.go -package=mock -destination=mock/service_mock.go
type PingServiceProvider interface {
	Ping(ctx context.Context) (pingmodel.PingPong, error)
}

type pingService struct {
	pingRepo pingrepo.PingRepositoryProvider
}

func NewPingService(
	pingRepo pingrepo.PingRepositoryProvider,
) PingServiceProvider {
	return &pingService{
		pingRepo: pingRepo,
	}
}