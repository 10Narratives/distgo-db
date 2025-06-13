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

	tests := []struct {
		name    string
		fields  fields
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful creation",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("SetDatabase", mock.Anything, name, dspName).Return(nil)
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
					}, nil)
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
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
					m.On("SetDatabase", mock.Anything, name, dspName).Return(errors.New("storage error"))
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "get database error",
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("SetDatabase", mock.Anything, name, dspName).Return(nil)
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{}, errors.New("not found"))
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(storage)
			service := databasesrv.New(storage)

			got, err := service.CreateDatabase(context.Background(), dbID, dspName)
			tt.wantVal(t, got)
			tt.wantErr(t, err)
			storage.AssertExpectations(t)
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

	tests := []struct {
		name    string
		args    string
		fields  fields
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: name,
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{
						Name:        name,
						DisplayName: dspName,
					}, nil)
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				db := got.(databasemodels.Database)
				assert.Equal(tt, name, db.Name)
				assert.Equal(tt, dspName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			args: name,
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(databasemodels.Database{}, errors.New("not found"))
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(storage)
			service := databasesrv.New(storage)

			got, err := service.Database(context.Background(), tt.args)
			tt.wantVal(t, got)
			tt.wantErr(t, err)
			storage.AssertExpectations(t)
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

	tests := []struct {
		name string
		args struct {
			size  int32
			token string
		}
		fields  fields
		wantVal func(tt require.TestingT, got []databasemodels.Database, token string, _ ...interface{})
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "single page",
			args: struct {
				size  int32
				token string
			}{size: 2, token: ""},
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2})
				},
			},
			wantVal: func(tt require.TestingT, got []databasemodels.Database, token string, _ ...interface{}) {
				require.Len(tt, got, 2)
				assert.Equal(tt, db1, got[0])
				assert.Equal(tt, db2, got[1])
				assert.Empty(tt, token)
			},
			wantErr: require.NoError,
		},
		{
			name: "pagination",
			args: struct {
				size  int32
				token string
			}{size: 1, token: "databases/db1"},
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Databases", mock.Anything).Return([]databasemodels.Database{db1, db2})
				},
			},
			wantVal: func(tt require.TestingT, got []databasemodels.Database, token string, _ ...interface{}) {
				require.Len(tt, got, 1)
				assert.Equal(tt, db2, got[0])
				assert.Empty(tt, token)
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(storage)
			service := databasesrv.New(storage)

			dbs, token, err := service.Databases(context.Background(), tt.args.size, tt.args.token)
			tt.wantVal(t, dbs, token)
			tt.wantErr(t, err)
			storage.AssertExpectations(t)
		})
	}
}

func TestService_DeleteDatabase(t *testing.T) {
	t.Parallel()

	const name = "databases/db123"

	type fields struct {
		setupStorageMock func(m *mocks.DatabaseStorage)
	}

	tests := []struct {
		name    string
		args    string
		fields  fields
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: name,
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("DeleteDatabase", mock.Anything, name).Return(nil)
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "delete error",
			args: name,
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("DeleteDatabase", mock.Anything, name).Return(errors.New("delete failed"))
				},
			},
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(storage)
			service := databasesrv.New(storage)

			err := service.DeleteDatabase(context.Background(), tt.args)
			tt.wantErr(t, err)
			storage.AssertExpectations(t)
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

	tests := []struct {
		name string
		args struct {
			db    databasemodels.Database
			paths []string
		}
		fields  fields
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success update display name",
			args: struct {
				db    databasemodels.Database
				paths []string
			}{
				db: databasemodels.Database{
					Name:        name,
					DisplayName: dspName,
				},
				paths: []string{"display_name"},
			},
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(existingDB, nil)
					m.On("SetDatabase", mock.Anything, oldName, dspName).Return(nil)
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				db := got.(databasemodels.Database)
				assert.Equal(tt, oldName, db.Name)
				assert.Equal(tt, dspName, db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field in paths",
			args: struct {
				db    databasemodels.Database
				paths []string
			}{
				db:    updatedDB,
				paths: []string{"invalid_field"},
			},
			fields: fields{
				setupStorageMock: func(m *mocks.DatabaseStorage) {
					m.On("Database", mock.Anything, name).Return(existingDB, nil)
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := mocks.NewDatabaseStorage(t)
			tt.fields.setupStorageMock(storage)
			service := databasesrv.New(storage)

			got, err := service.UpdateDatabase(context.Background(), tt.args.db, tt.args.paths)
			tt.wantVal(t, got)
			tt.wantErr(t, err)
			storage.AssertExpectations(t)
		})
	}
}
