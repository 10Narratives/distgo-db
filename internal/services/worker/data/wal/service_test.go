package walsrv_test

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// 	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
// 	"github.com/10Narratives/distgo-db/internal/services/worker/data/wal/mocks"
// )

// func TestService_Entries(t *testing.T) {
// 	t.Parallel()

// 	now := time.Now().UTC()
// 	testEntries := []walmodels.WALEntry{
// 		{
// 			ID:        uuid.New(),
// 			Entity:    walmodels.EntityTypeDatabase,
// 			Mutation:  commonmodels.MutationTypeCreate,
// 			Timestamp: now,
// 		},
// 	}

// 	tests := []struct {
// 		name          string
// 		size          int32
// 		token         string
// 		from          time.Time
// 		to            time.Time
// 		mockSetup     func(*mocks.WALStorage)
// 		expectedError error
// 	}{
// 		{
// 			name:  "success",
// 			size:  10,
// 			token: "token",
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Entries", mock.Anything, int32(10), "token", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
// 					Return(testEntries, "next-token", nil)
// 			},
// 		},
// 		{
// 			name:  "invalid page size - zero",
// 			size:  0,
// 			token: "token",
// 			mockSetup: func(m *mocks.WALStorage) {
// 				// No expectations - should fail before calling storage
// 			},
// 			expectedError: walsrv.ErrInvalidPageSize,
// 		},
// 		{
// 			name:  "invalid page size - too large",
// 			size:  1001,
// 			token: "token",
// 			mockSetup: func(m *mocks.WALStorage) {
// 				// No expectations - should fail before calling storage
// 			},
// 			expectedError: walsrv.ErrInvalidPageSize,
// 		},
// 		{
// 			name:  "storage error",
// 			size:  10,
// 			token: "token",
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Entries", mock.Anything, int32(10), "token", mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
// 					Return(nil, "", errors.New("storage error"))
// 			},
// 			expectedError: walsrv.ErrStorageOperation,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storageMock := mocks.NewWALStorage(t)
// 			if tt.mockSetup != nil {
// 				tt.mockSetup(storageMock)
// 			}

// 			service := walsrv.New(storageMock)
// 			entries, nextToken, err := service.Entries(context.Background(), tt.size, tt.token, tt.from, tt.to)

// 			if tt.expectedError != nil {
// 				require.ErrorIs(t, err, tt.expectedError)
// 				return
// 			}

// 			require.NoError(t, err)
// 			require.Equal(t, testEntries, entries)
// 			require.Equal(t, "next-token", nextToken)
// 			storageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_Truncate(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name          string
// 		before        time.Time
// 		mockSetup     func(*mocks.WALStorage)
// 		expectedError error
// 	}{
// 		{
// 			name:   "success",
// 			before: time.Now(),
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Truncate", mock.Anything, mock.AnythingOfType("time.Time")).
// 					Return(nil)
// 			},
// 		},
// 		{
// 			name:          "invalid before time - zero",
// 			before:        time.Time{},
// 			mockSetup:     func(m *mocks.WALStorage) {},
// 			expectedError: walsrv.ErrInvalidBeforeTime,
// 		},
// 		{
// 			name:   "storage error",
// 			before: time.Now(),
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Truncate", mock.Anything, mock.AnythingOfType("time.Time")).
// 					Return(errors.New("storage error"))
// 			},
// 			expectedError: walsrv.ErrStorageOperation,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storageMock := mocks.NewWALStorage(t)
// 			if tt.mockSetup != nil {
// 				tt.mockSetup(storageMock)
// 			}

// 			service := walsrv.New(storageMock)
// 			err := service.Truncate(context.Background(), tt.before)

// 			if tt.expectedError != nil {
// 				require.ErrorIs(t, err, tt.expectedError)
// 				return
// 			}

// 			require.NoError(t, err)
// 			storageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_Append(t *testing.T) {
// 	t.Parallel()

// 	testPayload := walmodels.DatabasePayload{
// 		Key: databasemodels.Key{Database: "test-db"},
// 	}

// 	tests := []struct {
// 		name          string
// 		entityType    walmodels.EntityType
// 		mutation      commonmodels.MutationType
// 		payload       interface{}
// 		mockSetup     func(*mocks.WALStorage)
// 		expectedError error
// 	}{
// 		{
// 			name:       "success",
// 			entityType: walmodels.EntityTypeDatabase,
// 			mutation:   commonmodels.MutationTypeCreate,
// 			payload:    testPayload,
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Append", mock.Anything, mock.AnythingOfType("walmodels.WALEntry")).
// 					Return(nil)
// 			},
// 		},
// 		{
// 			name:          "payload marshal error",
// 			entityType:    walmodels.EntityTypeDatabase,
// 			mutation:      commonmodels.MutationTypeCreate,
// 			payload:       make(chan int), // Unmarshalable type
// 			mockSetup:     func(m *mocks.WALStorage) {},
// 			expectedError: walsrv.ErrPayloadMarshal,
// 		},
// 		{
// 			name:       "storage error",
// 			entityType: walmodels.EntityTypeDatabase,
// 			mutation:   commonmodels.MutationTypeCreate,
// 			payload:    testPayload,
// 			mockSetup: func(m *mocks.WALStorage) {
// 				m.On("Append", mock.Anything, mock.AnythingOfType("walmodels.WALEntry")).
// 					Return(errors.New("storage error"))
// 			},
// 			expectedError: walsrv.ErrStorageOperation,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storageMock := mocks.NewWALStorage(t)
// 			if tt.mockSetup != nil {
// 				tt.mockSetup(storageMock)
// 			}

// 			service := walsrv.New(storageMock)
// 			err := service.Append(context.Background(), tt.entityType, tt.mutation, tt.payload)

// 			if tt.expectedError != nil {
// 				require.ErrorIs(t, err, tt.expectedError)
// 				return
// 			}

// 			require.NoError(t, err)
// 			storageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestCreateDatabaseEntry(t *testing.T) {
// 	t.Parallel()

// 	storageMock := mocks.NewWALStorage(t)
// 	service := walsrv.New(storageMock)

// 	key := databasemodels.Key{Database: "test-db"}
// 	db := &databasemodels.Database{Name: "test-db"}

// 	storageMock.On("Append", mock.Anything, mock.MatchedBy(func(entry walmodels.WALEntry) bool {
// 		var payload walmodels.DatabasePayload
// 		err := json.Unmarshal(entry.Payload, &payload)
// 		return err == nil &&
// 			payload.Key.Database == key.Database &&
// 			payload.Database.Name == db.Name
// 	})).Return(nil)

// 	err := service.CreateDatabaseEntry(context.Background(), commonmodels.MutationTypeCreate, key, db)
// 	require.NoError(t, err)
// 	storageMock.AssertExpectations(t)
// }

// func TestCreateCollectionEntry(t *testing.T) {
// 	t.Parallel()

// 	storageMock := mocks.NewWALStorage(t)
// 	service := walsrv.New(storageMock)

// 	key := collectionmodels.Key{Database: "test-db", Collection: "test-collection"}
// 	coll := &collectionmodels.Collection{Name: "test-collection"}

// 	storageMock.On("Append", mock.Anything, mock.MatchedBy(func(entry walmodels.WALEntry) bool {
// 		var payload walmodels.CollectionPayload
// 		err := json.Unmarshal(entry.Payload, &payload)
// 		return err == nil &&
// 			payload.Key.Collection == key.Collection &&
// 			payload.Collection.Name == coll.Name
// 	})).Return(nil)

// 	err := service.CreateCollectionEntry(context.Background(), commonmodels.MutationTypeCreate, key, coll)
// 	require.NoError(t, err)
// 	storageMock.AssertExpectations(t)
// }

// func TestCreateDocumentEntry(t *testing.T) {
// 	t.Parallel()

// 	storageMock := mocks.NewWALStorage(t)
// 	service := walsrv.New(storageMock)

// 	key := documentmodels.Key{Database: "test-db", Collection: "test-collection", Document: "test-doc"}
// 	doc := &documentmodels.Document{ID: "test-doc"}

// 	storageMock.On("Append", mock.Anything, mock.MatchedBy(func(entry walmodels.WALEntry) bool {
// 		var payload walmodels.DocumentPayload
// 		err := json.Unmarshal(entry.Payload, &payload)
// 		return err == nil &&
// 			payload.Key.Document == key.Document &&
// 			payload.Document.ID == doc.ID
// 	})).Return(nil)

// 	err := service.CreateDocumentEntry(context.Background(), commonmodels.MutationTypeCreate, key, doc)
// 	require.NoError(t, err)
// 	storageMock.AssertExpectations(t)
// }
