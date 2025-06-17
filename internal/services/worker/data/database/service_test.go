package databasesrv_test

import (
	"context"
	"errors"
	"testing"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	walmocks "github.com/10Narratives/distgo-db/internal/services/worker/data/common/mocks"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	mocks "github.com/10Narratives/distgo-db/internal/services/worker/data/database/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateDatabase(t *testing.T) {
	t.Parallel()

	const (
		dbID    = "db123"
		name    = "databases/db123"
		dspName = "Test DB"
	)

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
		setupWALMock     func(m *walmocks.WALStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			dbID    string
			dspName string
		}
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
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("CreateDatabase", mock.Anything, key, dspName).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
					}, nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					// entry := walmodels.WALEntry{
					// 	ID:       name,
					// 	Target:   "database",
					// 	Type:     commonmodels.MutationTypeCreate,
					// 	NewValue: dspName,
					// }
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					dbID    string
					dspName string
				}{
					dbID:    dbID,
					dspName: dspName,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, name, db.Name)
				assert.Equal(tt, dspName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "storage returns error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("CreateDatabase", mock.Anything, key, dspName).
						Return(databasemodels.Database{}, errors.New("storage error"))
				},
				setupWALMock: func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					dbID    string
					dspName string
				}{
					dbID:    dbID,
					dspName: dspName,
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "storage error")
			},
		},
		{
			name: "WAL logging fails",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("CreateDatabase", mock.Anything, key, dspName).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
					}, nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					// entry := walmodels.WALEntry{
					// 	ID:        name,
					// 	Target:    "database",
					// 	Type:      commonmodels.MutationTypeCreate,
					// 	NewValue:  dspName,
					// 	Timestamp: time.Now().UTC(),
					// }
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					dbID    string
					dspName string
				}{
					dbID:    dbID,
					dspName: dspName,
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "failed to log WAL entry")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewDatabaseStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := databasesrv.New(store, walStore)
			res, err := service.CreateDatabase(tt.args.ctx, tt.args.req.dbID, tt.args.req.dspName)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}

func TestService_DeleteDatabase(t *testing.T) {
	t.Parallel()

	const (
		name = "databases/db123"
		dbID = "db123"
	)

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
		setupWALMock     func(m *walmocks.WALStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			name string
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{
						Name:        name,
						DisplayName: "Test DB",
					}, nil)
					m.On("DeleteDatabase", mock.Anything, key).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: require.Empty,
			wantErr: require.NoError,
		},
		{
			name: "delete error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{
						Name:        name,
						DisplayName: "Test DB",
					}, nil)
					m.On("DeleteDatabase", mock.Anything, key).Return(errors.New("delete failed"))
				},
				setupWALMock: func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "delete failed")
			},
		},
		{
			name: "WAL logging fails",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{
						Name:        name,
						DisplayName: "Test DB",
					}, nil)
					m.On("DeleteDatabase", mock.Anything, key).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "failed to log WAL entry")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewDatabaseStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := databasesrv.New(store, walStore)
			err := service.DeleteDatabase(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, nil)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}

func TestService_UpdateDatabase(t *testing.T) {
	t.Parallel()

	const (
		name           = "databases/db123"
		displayName    = "Updated Name"
		oldDisplayName = "Old Name"
	)

	existingDB := databasemodels.Database{
		Name:        name,
		DisplayName: oldDisplayName,
	}

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
		setupWALMock     func(m *walmocks.WALStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			db    databasemodels.Database
			paths []string
		}
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
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
					m.On("UpdateDatabase", mock.Anything, key, displayName).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					db    databasemodels.Database
					paths []string
				}{
					db: databasemodels.Database{
						Name:        name,
						DisplayName: displayName,
					},
					paths: []string{"display_name"},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, name, db.Name)
				assert.Equal(tt, displayName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field in paths",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					db    databasemodels.Database
					paths []string
				}{
					db: databasemodels.Database{
						Name:        name,
						DisplayName: displayName,
					},
					paths: []string{"invalid_field"},
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "unknown field")
			},
		},
		{
			name: "WAL logging fails",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(existingDB, nil)
					m.On("UpdateDatabase", mock.Anything, key, displayName).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					db    databasemodels.Database
					paths []string
				}{
					db: databasemodels.Database{
						Name:        name,
						DisplayName: displayName,
					},
					paths: []string{"display_name"},
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "failed to log WAL entry")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewDatabaseStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := databasesrv.New(store, walStore)
			got, err := service.UpdateDatabase(tt.args.ctx, tt.args.req.db, tt.args.req.paths)

			tt.wantVal(t, got)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}

func TestService_Database(t *testing.T) {
	t.Parallel()

	const (
		name        = "databases/db123"
		displayName = "Test DB"
	)

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			name string
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{
						Name:        name,
						DisplayName: displayName,
					}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, name, db.Name)
				assert.Equal(tt, displayName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{}, errors.New("not found"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
		{
			name: "storage error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					key := databasemodels.NewKey(name)
					m.On("Database", mock.Anything, key).Return(databasemodels.Database{}, errors.New("storage error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					name string
				}{name: name},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "storage error")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)

			service := databasesrv.New(store, nil)
			res, err := service.Database(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_Databases(t *testing.T) {
	t.Parallel()

	// Sample databases for testing
	db1 := databasemodels.Database{
		Name:        "databases/db1",
		DisplayName: "DB1",
	}
	db2 := databasemodels.Database{
		Name:        "databases/db2",
		DisplayName: "DB2",
	}
	db3 := databasemodels.Database{
		Name:        "databases/db3",
		DisplayName: "DB3",
	}

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			size  int32
			token string
		}
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal func(tt require.TestingT, got interface{}, nextToken string, i ...interface{})
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "single page full list",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					size  int32
					token string
				}{
					size: 10,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, nextToken string, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Len(tt, list, 2)
				assert.Equal(tt, db1, list[0])
				assert.Equal(tt, db2, list[1])
				assert.Empty(tt, nextToken)
			},
			wantErr: require.NoError,
		},
		{
			name: "first page of two",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2, db3})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					size  int32
					token string
				}{
					size: 1,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, nextToken string, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Len(tt, list, 1)
				assert.Equal(tt, db1, list[0])
				//assert.Equal(tt, db2.Name, nextToken)
			},
			wantErr: require.NoError,
		},
		{
			name: "second page of two",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2, db3})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					size  int32
					token string
				}{
					size:  1,
					token: db1.Name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, nextToken string, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Len(tt, list, 1)
				assert.Equal(tt, db2, list[0])
				//assert.Equal(tt, db3.Name, nextToken)
			},
			wantErr: require.NoError,
		},
		{
			name: "last page",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2, db3})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					size  int32
					token string
				}{
					size:  1,
					token: db2.Name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, nextToken string, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Len(tt, list, 1)
				assert.Equal(tt, db3, list[0])
				assert.Empty(tt, nextToken)
			},
			wantErr: require.NoError,
		},
		{
			name: "empty list",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					size  int32
					token string
				}{
					size: 10,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, nextToken string, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Empty(tt, list)
				assert.Empty(tt, nextToken)
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)

			service := databasesrv.New(store, nil)
			list, nextToken, err := service.Databases(tt.args.ctx, tt.args.req.size, tt.args.req.token)

			tt.wantVal(t, list, nextToken)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}
