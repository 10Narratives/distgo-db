package databasesrv

import (
	"context"

	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name DatabaseStorage --output ./mocks/
type DatabaseStorage interface {
	Database(ctx context.Context, name string) (databasemodels.Database, error)
	Databases(ctx context.Context) []databasemodels.Database
	SetDatabase(ctx context.Context, name, displayName string) error
	DeleteDatabase(ctx context.Context, name string) error
}

type Service struct {
	databaseStore DatabaseStorage
}

var _ databasegrpc.DatabaseService = &Service{}

func New(databaseStore DatabaseStorage) *Service {
	return &Service{
		databaseStore: databaseStore,
	}
}

func (s *Service) CreateDatabase(ctx context.Context, databaseID string, displayName string) (databasemodels.Database, error) {
	err := s.databaseStore.SetDatabase(ctx, "databases/"+databaseID, displayName)
	if err != nil {
		return databasemodels.Database{}, err
	}
	return s.databaseStore.Database(ctx, "databases/"+databaseID)
}

func (s *Service) Database(ctx context.Context, name string) (databasemodels.Database, error) {
	return s.databaseStore.Database(ctx, name)
}

func (s *Service) Databases(ctx context.Context, size int32, token string) ([]databasemodels.Database, string, error) {
	allDbs := s.databaseStore.Databases(ctx)

	if len(allDbs) == 0 {
		return []databasemodels.Database{}, "", nil
	}

	startIndex := 0
	if token != "" {
		for i, db := range allDbs {
			if db.Name == token {
				startIndex = i + 1
				break
			}
		}
	}
	endIndex := startIndex + int(size)
	if endIndex > len(allDbs) {
		endIndex = len(allDbs)
	}

	page := allDbs[startIndex:endIndex]

	var nextPageToken string
	if endIndex < len(allDbs) {
		nextPageToken = page[len(page)-1].Name
	}

	return page, nextPageToken, nil
}

func (s *Service) DeleteDatabase(ctx context.Context, name string) error {
	return s.databaseStore.DeleteDatabase(ctx, name)
}

func (s *Service) UpdateDatabase(ctx context.Context, database databasemodels.Database, paths []string) (databasemodels.Database, error) {
	existingDB, err := s.databaseStore.Database(ctx, database.Name)
	if err != nil {
		return databasemodels.Database{}, err
	}

	for _, path := range paths {
		switch path {
		case "display_name":
			existingDB.DisplayName = database.DisplayName
		default:
			return databasemodels.Database{}, status.Errorf(codes.InvalidArgument, "unknown field: %s", path)
		}
	}

	err = s.databaseStore.SetDatabase(ctx, existingDB.Name, existingDB.DisplayName)
	if err != nil {
		return databasemodels.Database{}, err
	}

	return existingDB, nil
}
