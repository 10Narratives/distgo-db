package databasesrv

import (
	"context"
	"errors"
	"time"

	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name DatabaseStorage --output ./mocks/
type DatabaseStorage interface {
	Database(ctx context.Context, key databasemodels.Key) (databasemodels.Database, error)
	Databases(ctx context.Context) ([]databasemodels.Database, error)
	CreateDatabase(ctx context.Context, key databasemodels.Key, displayName string) (databasemodels.Database, error)
	UpdateDatabase(ctx context.Context, key databasemodels.Key, displayName string) error
	DeleteDatabase(ctx context.Context, key databasemodels.Key) error
}

//go:generate mockery --name WAlService --output ./mocks/
type WAlService interface {
	CreateDatabaseEntry(ctx context.Context, mutation commonmodels.MutationType, key databasemodels.Key, db *databasemodels.Database) error
}

type Service struct {
	databaseStorage DatabaseStorage
	walService      WAlService
}

var _ databasegrpc.DatabaseService = &Service{}

func New(
	databaseStore DatabaseStorage,
	walService WAlService,
) *Service {
	return &Service{
		databaseStorage: databaseStore,
		walService:      walService,
	}
}

func (s *Service) CreateDatabase(ctx context.Context, databaseID string, displayName string) (databasemodels.Database, error) {
	key := databasemodels.NewKey("databases/" + databaseID)

	newDB := databasemodels.Database{
		Name:        "databases/" + databaseID,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.walService.CreateDatabaseEntry(ctx, commonmodels.MutationTypeCreate, key, &newDB); err != nil {
		return databasemodels.Database{}, errors.New("failed to create WAL entry: " + err.Error())
	}

	db, err := s.databaseStorage.CreateDatabase(ctx, key, displayName)
	if err != nil {
		return databasemodels.Database{}, err
	}

	return db, nil
}

func (s *Service) UpdateDatabase(ctx context.Context, database databasemodels.Database, paths []string) (databasemodels.Database, error) {
	key := databasemodels.NewKey(database.Name)
	existingDB, err := s.databaseStorage.Database(ctx, key)
	if err != nil {
		return databasemodels.Database{}, err
	}

	for _, path := range paths {
		switch path {
		case "display_name":
			updatedDB := existingDB
			updatedDB.DisplayName = database.DisplayName
			updatedDB.UpdatedAt = time.Now()
			if err := s.walService.CreateDatabaseEntry(ctx, commonmodels.MutationTypeUpdate, key, &updatedDB); err != nil {
				return databasemodels.Database{}, errors.New("failed to create WAL entry: " + err.Error())
			}

			if err := s.databaseStorage.UpdateDatabase(ctx, key, database.DisplayName); err != nil {
				return databasemodels.Database{}, err
			}

			existingDB.DisplayName = database.DisplayName
		default:
			return databasemodels.Database{}, status.Errorf(codes.InvalidArgument, "unknown field: %s", path)
		}
	}

	return existingDB, nil
}

func (s *Service) DeleteDatabase(ctx context.Context, name string) error {
	key := databasemodels.NewKey(name)

	if err := s.walService.CreateDatabaseEntry(ctx, commonmodels.MutationTypeDelete, key, nil); err != nil {
		return errors.New("failed to create WAL entry: " + err.Error())
	}

	if err := s.databaseStorage.DeleteDatabase(ctx, key); err != nil {
		return err
	}

	return nil
}

func (s *Service) Database(ctx context.Context, name string) (databasemodels.Database, error) {
	key := databasemodels.NewKey(name)
	return s.databaseStorage.Database(ctx, key)
}

func (s *Service) Databases(ctx context.Context, size int32, token string) ([]databasemodels.Database, string, error) {
	allDbs, err := s.databaseStorage.Databases(ctx)
	if err != nil {
		return []databasemodels.Database{}, "", err
	}

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
