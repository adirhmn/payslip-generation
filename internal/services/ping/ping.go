package ping

import (
	"context"

	pingmodel "payslip-generation-system/internal/entity/ping"
)

func (s *pingService) Ping(ctx context.Context) (pingmodel.PingPong, error) {
	err := s.pingRepo.Ping(ctx)
	if err != nil {
		return pingmodel.PingPong{
			Message: pingmodel.ErrorMessage,
		}, err
	}

	return pingmodel.PingPong{
		Message: pingmodel.SuccessMessage,
	}, nil
}