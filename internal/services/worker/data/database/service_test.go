package databasesrv_test

import (
	"context"
	"errors"
	"testing"
	"time"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	"github.com/10Narratives/distgo-db/internal/services/worker/data/database/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	name        string = "databases/test_db"
	displayName string = "Test Database"
	createdAt   time.Time
	updatedAt   time.Time
	key         = databasemodels.NewKey(name)
	database    = databasemodels.Database{
		Name:        name,
		DisplayName: displayName,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	errInternalDatabaseStorage = errors.New("internal database storage error")

	db1 = databasemodels.Database{
		Name:        "databases/db1",
		DisplayName: "Database 1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db2 = databasemodels.Database{
		Name:        "databases/db2",
		DisplayName: "Database 2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db3 = databasemodels.Database{
		Name:        "databases/db3",
		DisplayName: "Database 3",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	allDbs = []databasemodels.Database{db1, db2, db3}

	errInternalWALService = errors.New("internal WAL service error")

	databaseID   = "test_db"
	databaseName = "databases/" + databaseID
	newDB        = databasemodels.Database{
		Name:        databaseName,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	existingDB = databasemodels.Database{
		Name:        "databases/test_db",
		DisplayName: "Old Display Name",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	updatedDB = databasemodels.Database{
		Name:        "databases/test_db",
		DisplayName: "New Display Name",
		CreatedAt:   existingDB.CreatedAt,
		UpdatedAt:   time.Now(),
	}
)

type fields struct {
	setupDatabaseStorageMock func(m *mocks.DatabaseStorage)
	setupWALServiceMock      func(m *mocks.WAlService)
}

func TestService_Database(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful execution",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).Return(database, nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(t, ok)

				assert.Equal(t, database, db)
			},
			wantErr: require.NoError,
		},
		{
			name: "database storage error",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).
						Return(databasemodels.Database{}, errInternalDatabaseStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			databaseStorageMock := mocks.NewDatabaseStorage(t)
			tt.fields.setupDatabaseStorageMock(databaseStorageMock)

			service := databasesrv.New(databaseStorageMock, nil)

			db, err := service.Database(tt.args.ctx, tt.args.name)

			tt.wantVal(t, db)
			tt.wantErr(t, err)

			databaseStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_Databases(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		size  int32
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful execution with pagination",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return(allDbs, nil)
				},
			},
			args: args{
				ctx:   context.Background(),
				size:  2,
				token: "",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]databasemodels.Database)
				require.True(t, ok)

				expected := []databasemodels.Database{db1, db2}
				assert.Equal(t, expected, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "empty databases list",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{}, nil)
				},
			},
			args: args{
				ctx:   context.Background(),
				size:  10,
				token: "",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]databasemodels.Database)
				require.True(t, ok)
				assert.Empty(t, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "pagination with token",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return(allDbs, nil)
				},
			},
			args: args{
				ctx:   context.Background(),
				size:  1,
				token: "databases/db1",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]databasemodels.Database)
				require.True(t, ok)

				expected := []databasemodels.Database{db2}
				assert.Equal(t, expected, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "storage error",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return(nil, errInternalDatabaseStorage)
				},
			},
			args: args{
				ctx:   context.Background(),
				size:  10,
				token: "",
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			databaseStorageMock := mocks.NewDatabaseStorage(t)
			tt.fields.setupDatabaseStorageMock(databaseStorageMock)

			service := databasesrv.New(databaseStorageMock, nil)

			page, nextPageToken, err := service.Databases(tt.args.ctx, tt.args.size, tt.args.token)

			tt.wantVal(t, page)
			tt.wantErr(t, err)

			if err == nil && len(page) > 0 && nextPageToken != "" {
				assert.Equal(t, tt.args.size, int32(len(page)))
			}

			databaseStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateDatabase(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx         context.Context
		databaseID  string
		displayName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful creation",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeCreate, key, mock.Anything).
						Return(nil)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("CreateDatabase", mock.Anything, key, displayName).
						Return(newDB, nil)
				},
			},
			args: args{
				ctx:         context.Background(),
				databaseID:  databaseID,
				displayName: displayName,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(t, ok)

				assert.Equal(t, newDB.Name, db.Name)
				assert.Equal(t, newDB.DisplayName, db.DisplayName)
				assert.WithinDuration(t, time.Now(), db.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), db.UpdatedAt, time.Second)
			},
			wantErr: require.NoError,
		},
		{
			name: "WAL service error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeCreate, key, mock.Anything).
						Return(errInternalWALService)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {},
			},
			args: args{
				ctx:         context.Background(),
				databaseID:  databaseID,
				displayName: displayName,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "failed to create WAL entry: "+errInternalWALService.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeCreate, key, mock.Anything).
						Return(nil)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("CreateDatabase", mock.Anything, key, displayName).
						Return(databasemodels.Database{}, errInternalDatabaseStorage)
				},
			},
			args: args{
				ctx:         context.Background(),
				databaseID:  databaseID,
				displayName: displayName,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			databaseStorageMock := mocks.NewDatabaseStorage(t)
			tt.fields.setupDatabaseStorageMock(databaseStorageMock)

			service := databasesrv.New(databaseStorageMock, walServiceMock)

			db, err := service.CreateDatabase(tt.args.ctx, tt.args.databaseID, tt.args.displayName)

			tt.wantVal(t, db)
			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			databaseStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateDatabase(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx      context.Context
		database databasemodels.Database
		paths    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful update of display name",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
					m.On("UpdateDatabase", mock.Anything, key, updatedDB.DisplayName).Return(nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				database: databasemodels.Database{
					Name:        existingDB.Name,
					DisplayName: "New Display Name",
				},
				paths: []string{"display_name"},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(t, ok)

				assert.Equal(t, updatedDB.DisplayName, db.DisplayName)
				assert.WithinDuration(t, time.Now(), db.UpdatedAt, time.Second)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx: context.Background(),
				database: databasemodels.Database{
					Name:        existingDB.Name,
					DisplayName: "New Display Name",
				},
				paths: []string{"unknown_field"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, status.Errorf(codes.InvalidArgument, "unknown field: unknown_field"), err.Error())
			},
		},
		{
			name: "storage error when fetching database",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).
						Return(databasemodels.Database{}, errInternalDatabaseStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx: context.Background(),
				database: databasemodels.Database{
					Name:        existingDB.Name,
					DisplayName: "New Display Name",
				},
				paths: []string{"display_name"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
		{
			name: "WAL service error",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(errInternalWALService)
				},
			},
			args: args{
				ctx: context.Background(),
				database: databasemodels.Database{
					Name:        existingDB.Name,
					DisplayName: "New Display Name",
				},
				paths: []string{"display_name"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "failed to create WAL entry: "+errInternalWALService.Error())
			},
		},
		{
			name: "storage error when updating database",
			fields: fields{
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
					m.On("UpdateDatabase", mock.Anything, key, updatedDB.DisplayName).
						Return(errInternalDatabaseStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				database: databasemodels.Database{
					Name:        existingDB.Name,
					DisplayName: "New Display Name",
				},
				paths: []string{"display_name"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			databaseStorageMock := mocks.NewDatabaseStorage(t)
			tt.fields.setupDatabaseStorageMock(databaseStorageMock)

			service := databasesrv.New(databaseStorageMock, walServiceMock)

			db, err := service.UpdateDatabase(tt.args.ctx, tt.args.database, tt.args.paths)

			tt.wantVal(t, db)
			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			databaseStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_DeleteDatabase(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful deletion",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*databasemodels.Database)(nil)).
						Return(nil)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("DeleteDatabase", mock.Anything, key).
						Return(nil)
				},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: require.NoError,
		},
		{
			name: "WAL service error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*databasemodels.Database)(nil)).
						Return(errInternalWALService)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "failed to create WAL entry: "+errInternalWALService.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateDatabaseEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*databasemodels.Database)(nil)).
						Return(nil)
				},
				setupDatabaseStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("DeleteDatabase", mock.Anything, key).
						Return(errInternalDatabaseStorage)
				},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalDatabaseStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			databaseStorageMock := mocks.NewDatabaseStorage(t)
			tt.fields.setupDatabaseStorageMock(databaseStorageMock)

			service := databasesrv.New(databaseStorageMock, walServiceMock)

			err := service.DeleteDatabase(tt.args.ctx, tt.args.name)

			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			databaseStorageMock.AssertExpectations(t)
		})
	}
}
