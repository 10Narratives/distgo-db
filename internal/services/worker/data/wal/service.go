package walsrv

// import (
// 	"context"
// 	"time"

// 	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// )

// //go:generate mockery --name WALStorage --output ./mocks/
// type WALStorage interface {
// 	Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error)
// 	Truncate(ctx context.Context, before time.Time) error
// }

// type Service struct {
// 	walStorage WALStorage
// }

// var _ walgrpc.WALService = &Service{}

// func New(walStorage WALStorage) *Service {
// 	return &Service{
// 		walStorage: walStorage,
// 	}
// }

// func (s *Service) TruncateWAL(ctx context.Context, before time.Time) error {
// 	return s.walStorage.Truncate(ctx, before)
// }

// func (s *Service) WALEntries(ctx context.Context, size int32, token string, from time.Time, to time.Time) ([]walmodels.WALEntry, string, error) {
// 	return s.walStorage.Entries(ctx, size, token, from, to)
// }
