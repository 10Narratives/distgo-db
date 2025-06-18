package databasesrv

import (
	"context"

	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	commonsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockery --name DatabaseStorage --output ./mocks/
type DatabaseStorage interface {
	Database(ctx context.Context, key databasemodels.Key) (databasemodels.Database, error)
	Databases(ctx context.Context) []databasemodels.Database
	CreateDatabase(ctx context.Context, key databasemodels.Key, displayName string) (databasemodels.Database, error)
	UpdateDatabase(ctx context.Context, key databasemodels.Key, displayName string) error
	DeleteDatabase(ctx context.Context, key databasemodels.Key) error
}

type Service struct {
	databaseStore DatabaseStorage
	walStorage    commonsrv.WALStorage
}

var _ databasegrpc.DatabaseService = &Service{}

func New(
	databaseStore DatabaseStorage,
	walStorage commonsrv.WALStorage,
) *Service {
	return &Service{
		databaseStore: databaseStore,
		walStorage:    walStorage,
	}
}

func (s *Service) CreateDatabase(ctx context.Context, databaseID string, displayName string) (databasemodels.Database, error) {
	key := databasemodels.NewKey("databases/" + databaseID)

	// Create the database in the storage
	db, err := s.databaseStore.CreateDatabase(ctx, key, displayName)
	if err != nil {
		return databasemodels.Database{}, err
	}

	// // Log the creation in WAL
	// entry := walmodels.WALEntry{
	// 	ID:        db.Name,
	// 	Target:    "database",
	// 	Type:      commonmodels.MutationTypeCreate,
	// 	NewValue:  db.DisplayName,
	// 	Timestamp: time.Now(),
	// }
	// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
	// 	return databasemodels.Database{}, status.Errorf(codes.Internal, "failed to log WAL entry: %v", err)
	// }

	return db, nil
}

func (s *Service) Database(ctx context.Context, name string) (databasemodels.Database, error) {
	key := databasemodels.NewKey(name)
	return s.databaseStore.Database(ctx, key)
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
	key := databasemodels.NewKey(name)

	// db, err := s.databaseStore.Database(ctx, key)
	// if err != nil {
	// 	return err
	// }

	if err := s.databaseStore.DeleteDatabase(ctx, key); err != nil {
		return err
	}

	// entry := walmodels.WALEntry{
	// 	ID:        db.Name,
	// 	Target:    "database",
	// 	Type:      commonmodels.MutationTypeDelete,
	// 	OldValue:  db.DisplayName,
	// 	Timestamp: time.Now(),
	// }
	// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
	// 	return status.Errorf(codes.Internal, "failed to log WAL entry: %v", err)
	// }

	return nil
}

func (s *Service) UpdateDatabase(ctx context.Context, database databasemodels.Database, paths []string) (databasemodels.Database, error) {
	key := databasemodels.NewKey(database.Name)
	existingDB, err := s.databaseStore.Database(ctx, key)
	if err != nil {
		return databasemodels.Database{}, err
	}

	for _, path := range paths {
		switch path {
		case "display_name":
			// Update the database in storage
			if err := s.databaseStore.UpdateDatabase(ctx, key, database.DisplayName); err != nil {
				return databasemodels.Database{}, err
			}

			// // Log the update in WAL
			// entry := walmodels.WALEntry{
			// 	ID:        existingDB.Name,
			// 	Target:    "database",
			// 	Type:      commonmodels.MutationTypeUpdate,
			// 	OldValue:  existingDB.DisplayName,
			// 	NewValue:  database.DisplayName,
			// 	Timestamp: time.Now(),
			// }
			// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
			// 	return databasemodels.Database{}, status.Errorf(codes.Internal, "failed to log WAL entry: %v", err)
			// }

			existingDB.DisplayName = database.DisplayName
		default:
			return databasemodels.Database{}, status.Errorf(codes.InvalidArgument, "unknown field: %s", path)
		}
	}

	return existingDB, nil
}
