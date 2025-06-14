package databasesrv_test

import (
	"context"
	"errors"
	"testing"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
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
					m.On("CreateDatabase", mock.Anything, name, dspName).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
					}, nil)
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
			name: "set database error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("CreateDatabase", mock.Anything, name, dspName).
						Return(databasemodels.Database{}, errors.New("storage error"))
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
				assert.Contains(tt, err.Error(), "storage error")
			},
		},
		{
			name: "get database error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("CreateDatabase", mock.Anything, name, dspName).Return(nil)
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{}, errors.New("not found"))
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
				assert.Contains(tt, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)
			service := databasesrv.New(store)

			res, err := service.CreateDatabase(tt.args.ctx, tt.args.req.dbID, tt.args.req.dspName)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_Database(t *testing.T) {
	t.Parallel()

	const (
		name    = "databases/db123"
		dspName = "Test DB"
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
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
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
				assert.Equal(tt, dspName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{}, errors.New("not found"))
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)
			service := databasesrv.New(store)

			res, err := service.Database(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_Databases(t *testing.T) {
	t.Parallel()

	db1 := databasemodels.Database{
		Name:        "databases/db1",
		DisplayName: "DB1",
	}
	db2 := databasemodels.Database{
		Name:        "databases/db2",
		DisplayName: "DB2",
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
		wantVal func(tt require.TestingT, got interface{}, i ...interface{})
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "single page",
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
					size: 2,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				require.Len(tt, list, 2)
				assert.Equal(tt, db1, list[0])
				assert.Equal(tt, db2, list[1])
			},
			wantErr: require.NoError,
		},
		{
			name: "with pagination",
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
					size:  1,
					token: db1.Name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				require.Len(tt, list, 1)
				assert.Equal(tt, db2, list[0])
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
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Empty(tt, list)
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
			service := databasesrv.New(store)

			list, _, err := service.Databases(tt.args.ctx, tt.args.req.size, tt.args.req.token)

			tt.wantVal(t, list)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_DeleteDatabase(t *testing.T) {
	t.Parallel()

	const name = "databases/db123"

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
					m.On("DeleteDatabase", mock.Anything, name).Return(nil)
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
					m.On("DeleteDatabase", mock.Anything, name).Return(errors.New("delete failed"))
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
				assert.Contains(tt, err.Error(), "delete failed")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)
			service := databasesrv.New(store)

			err := service.DeleteDatabase(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, nil)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_UpdateDatabase(t *testing.T) {
	t.Parallel()

	const (
		name    = "databases/db123"
		oldName = "databases/db456"
		dspName = "Updated Name"
	)

	existingDB := databasemodels.Database{
		Name:        oldName,
		DisplayName: "Old Name",
	}
	updatedDB := databasemodels.Database{
		Name:        name,
		DisplayName: dspName,
	}

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
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
			name: "success update display name",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(existingDB, nil)
					m.On("UpdateDatabase", mock.Anything, oldName, dspName).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					db    databasemodels.Database
					paths []string
				}{
					db:    updatedDB,
					paths: []string{"display_name"},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, oldName, db.Name)
				assert.Equal(tt, dspName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field in paths",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(existingDB, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					db    databasemodels.Database
					paths []string
				}{
					db:    updatedDB,
					paths: []string{"invalid_field"},
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "unknown field")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(store)
			service := databasesrv.New(store)

			got, err := service.UpdateDatabase(tt.args.ctx, tt.args.req.db, tt.args.req.paths)

			tt.wantVal(t, got)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}
