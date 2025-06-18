package datastorage_test

import (
	"context"
	"testing"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateDatabase(t *testing.T) {
	t.Parallel()

	type fields struct {
		databases map[databasemodels.Key]databasemodels.Database
	}
	type args struct {
		ctx         context.Context
		key         databasemodels.Key
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
			name: "success create new database",
			fields: fields{
				databases: nil,
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db123",
				},
				displayName: "Test DB",
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, "databases/db123", db.Name)
				assert.Equal(tt, "Test DB", db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "error when database already exists",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{
						Database: "db123",
					}: {
						Name:        "databases/db123",
						DisplayName: "Existing DB",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db123",
				},
				displayName: "New DB",
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "already exists")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := datastorage.NewOf(
				tt.fields.databases,
				map[collectionmodels.Key]collectionmodels.Collection{},
				map[documentmodels.Key]documentmodels.Document{},
			)
			created, err := storage.CreateDatabase(tt.args.ctx, tt.args.key, tt.args.displayName)

			tt.wantVal(t, created)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Database(t *testing.T) {
	t.Parallel()

	type fields struct {
		databases map[databasemodels.Key]databasemodels.Database
	}
	type args struct {
		ctx context.Context
		key databasemodels.Key
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success get existing database",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{
						Database: "db123",
					}: {
						Name:        "databases/db123",
						DisplayName: "Test DB",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db123",
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Equal(tt, "databases/db123", db.Name)
				assert.Equal(tt, "Test DB", db.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "error when database not found",
			fields: fields{
				databases: nil,
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db456",
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				db, ok := got.(databasemodels.Database)
				require.True(tt, ok)
				assert.Empty(tt, db)
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := datastorage.NewOf(
				tt.fields.databases,
				map[collectionmodels.Key]collectionmodels.Collection{},
				map[documentmodels.Key]documentmodels.Document{},
			)
			database, err := storage.Database(tt.args.ctx, tt.args.key)

			tt.wantVal(t, database)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Databases(t *testing.T) {
	t.Parallel()

	type fields struct {
		databases map[databasemodels.Key]databasemodels.Database
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
	}{
		{
			name: "returns empty list when no databases",
			fields: fields{
				databases: nil,
			},
			args: args{
				ctx: context.Background(),
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				assert.Empty(tt, list)
			},
		},
		{
			name: "returns all databases in storage",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{
						Database: "db1",
					}: {
						Name:        "databases/db1",
						DisplayName: "DB 1",
					},
					{
						Database: "db2",
					}: {
						Name:        "databases/db2",
						DisplayName: "DB 2",
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]databasemodels.Database)
				require.True(tt, ok)
				require.Len(tt, list, 2)

				names := make([]string, len(list))
				for _, db := range list {
					names = append(names, db.Name)
				}

				assert.Contains(tt, names, "databases/db1")
				assert.Contains(tt, names, "databases/db2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := datastorage.NewOf(
				tt.fields.databases,
				map[collectionmodels.Key]collectionmodels.Collection{},
				map[documentmodels.Key]documentmodels.Document{},
			)
			databases, _ := storage.Databases(tt.args.ctx)

			tt.wantVal(t, databases)
		})
	}
}

func TestStorage_DeleteDatabase(t *testing.T) {
	t.Parallel()

	type fields struct {
		databases map[databasemodels.Key]databasemodels.Database
	}
	type args struct {
		ctx context.Context
		key databasemodels.Key
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success delete existing database",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{
						Database: "db123",
					}: {
						Name:        "databases/db123",
						DisplayName: "Test DB",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db123",
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "error when database not found",
			fields: fields{
				databases: nil,
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db456",
				},
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := datastorage.NewOf(
				tt.fields.databases,
				map[collectionmodels.Key]collectionmodels.Collection{},
				map[documentmodels.Key]documentmodels.Document{},
			)
			err := storage.DeleteDatabase(tt.args.ctx, tt.args.key)

			tt.wantErr(t, err)
		})
	}
}

func TestStorage_UpdateDatabase(t *testing.T) {
	t.Parallel()

	type fields struct {
		databases map[databasemodels.Key]databasemodels.Database
	}
	type args struct {
		ctx         context.Context
		key         databasemodels.Key
		displayName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success update existing database",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{
						Database: "db123",
					}: {
						Name:        "databases/db123",
						DisplayName: "Old Name",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db123",
				},
				displayName: "New Name",
			},
			wantErr: require.NoError,
		},
		{
			name: "error when database not found",
			fields: fields{
				databases: nil,
			},
			args: args{
				ctx: context.Background(),
				key: databasemodels.Key{
					Database: "db456",
				},
				displayName: "New Name",
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := datastorage.NewOf(
				tt.fields.databases,
				map[collectionmodels.Key]collectionmodels.Collection{},
				map[documentmodels.Key]documentmodels.Document{},
			)
			err := storage.UpdateDatabase(tt.args.ctx, tt.args.key, tt.args.displayName)

			tt.wantErr(t, err)
		})
	}
}
