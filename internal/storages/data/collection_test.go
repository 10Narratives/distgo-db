package datastorage_test

import (
	"context"
	"testing"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dbID         = "db123"
	collectionID = "coll1"
	description  = "Test Collection"
)

var now = time.Now().UTC()

func TestStorage_Collection(t *testing.T) {
	t.Parallel()

	key := collectionmodels.Key{
		Database:   dbID,
		Collection: collectionID,
	}
	expectedColl := collectionmodels.Collection{
		Name:        "databases/" + dbID + "/collections/" + collectionID,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	type fields struct {
		collections map[collectionmodels.Key]collectionmodels.Collection
	}
	tests := []struct {
		name    string
		fields  fields
		args    collectionmodels.Key
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				collections: map[collectionmodels.Key]collectionmodels.Collection{
					key: expectedColl,
				},
			},
			args: key,
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll := got.(collectionmodels.Collection)
				assert.Equal(tt, expectedColl, coll)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				collections: nil,
			},
			args: collectionmodels.Key{
				Database:   dbID,
				Collection: "unknown",
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
			storage := datastorage.NewOf(
				nil,
				tt.fields.collections,
				nil,
			)

			coll, err := storage.Collection(context.Background(), tt.args)

			tt.wantVal(t, coll)
			tt.wantErr(t, err)
		})
	}
}

// func TestStorage_Collections(t *testing.T) {
// 	t.Parallel()

// 	parentKey := databasemodels.Key{Database: dbID}

// 	colls := []struct {
// 		key        collectionmodels.Key
// 		collection collectionmodels.Collection
// 	}{
// 		{
// 			key: collectionmodels.Key{
// 				Database:   dbID,
// 				Collection: "coll1",
// 			},
// 			collection: collectionmodels.Collection{
// 				Name:      "databases/db123/collections/coll1",
// 				CreatedAt: now,
// 				UpdatedAt: now,
// 			},
// 		},
// 		{
// 			key: collectionmodels.Key{
// 				Database:   dbID,
// 				Collection: "coll2",
// 			},
// 			collection: collectionmodels.Collection{
// 				Name:      "databases/db123/collections/coll2",
// 				CreatedAt: now,
// 				UpdatedAt: now,
// 			},
// 		},
// 		{
// 			key: collectionmodels.Key{
// 				Database:   "other_db",
// 				Collection: "coll3",
// 			},
// 			collection: collectionmodels.Collection{
// 				Name:      "databases/other_db/collections/coll3",
// 				CreatedAt: now,
// 				UpdatedAt: now,
// 			},
// 		},
// 	}

// 	type fields struct {
// 		collections map[collectionmodels.Key]collectionmodels.Collection
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    databasemodels.Key
// 		wantVal require.ValueAssertionFunc
// 	}{
// 		{
// 			name: "success list collections in db",
// 			fields: fields{
// 				collections: map[collectionmodels.Key]collectionmodels.Collection{
// 					colls[0].key: colls[0].collection,
// 					colls[1].key: colls[1].collection,
// 					colls[2].key: colls[2].collection,
// 				},
// 			},
// 			args: parentKey,
// 			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
// 				list := got.([]collectionmodels.Collection)
// 				require.Len(tt, list, 2)
// 				assert.Equal(tt, colls[0].collection.Name, list[0].Name)
// 				assert.Equal(tt, colls[1].collection.Name, list[1].Name)
// 			},
// 		},
// 		{
// 			name: "no collections",
// 			fields: fields{
// 				collections: nil,
// 			},
// 			args: parentKey,
// 			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
// 				list := got.([]collectionmodels.Collection)
// 				assert.Empty(tt, list)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			storage := datastorage.NewOf(
// 				map[databasemodels.Key]databasemodels.Database{},
// 				tt.fields.collections,
// 				nil,
// 			)

// 			res := storage.Collections(context.Background(), tt.args)

// 			tt.wantVal(t, res)
// 		})
// 	}
// }

func TestStorage_CreateCollection(t *testing.T) {
	t.Parallel()

	key := collectionmodels.Key{
		Database:   dbID,
		Collection: collectionID,
	}

	type fields struct {
		collections map[collectionmodels.Key]collectionmodels.Collection
	}
	type args struct {
		key         collectionmodels.Key
		description string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success create",
			fields: fields{
				collections: nil,
			},
			args: args{
				key:         key,
				description: description,
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll := got.(collectionmodels.Collection)
				assert.Equal(tt, "databases/"+dbID+"/collections/"+collectionID, coll.Name)
				assert.Equal(tt, description, coll.Description)
			},
			wantErr: require.NoError,
		},
		{
			name: "already exists",
			fields: fields{
				collections: map[collectionmodels.Key]collectionmodels.Collection{
					key: {},
				},
			},
			args: args{
				key:         key,
				description: description,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found") // в коде используется то же сообщение
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := datastorage.NewOf(
				nil,
				tt.fields.collections,
				nil,
			)

			res, err := storage.CreateCollection(context.Background(), tt.args.key, tt.args.description)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_DeleteCollection(t *testing.T) {
	t.Parallel()

	key := collectionmodels.Key{
		Database:   dbID,
		Collection: collectionID,
	}

	existing := map[collectionmodels.Key]collectionmodels.Collection{
		key: {
			Name: "databases/db123/collections/coll1",
		},
	}

	type fields struct {
		collections map[collectionmodels.Key]collectionmodels.Collection
	}
	tests := []struct {
		name    string
		fields  fields
		args    collectionmodels.Key
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success delete",
			fields: fields{
				collections: existing,
			},
			args: key,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err)
			},
		},
		{
			name: "not found",
			fields: fields{
				collections: nil,
			},
			args: collectionmodels.Key{
				Database:   dbID,
				Collection: "unknown",
			},
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
			storage := datastorage.NewOf(
				nil,
				tt.fields.collections,
				nil,
			)

			err := storage.DeleteCollection(context.Background(), tt.args)

			tt.wantErr(t, err)
		})
	}
}

func TestStorage_UpdateCollection(t *testing.T) {
	t.Parallel()

	key := collectionmodels.Key{
		Database:   dbID,
		Collection: collectionID,
	}

	existing := map[collectionmodels.Key]collectionmodels.Collection{
		key: {
			Name:      "databases/db123/collections/coll1",
			UpdatedAt: now.Add(-time.Hour),
		},
	}

	newDescription := "Updated Description"

	type fields struct {
		collections map[collectionmodels.Key]collectionmodels.Collection
	}
	type args struct {
		key         collectionmodels.Key
		description string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success update",
			fields: fields{
				collections: existing,
			},
			args: args{
				key:         key,
				description: newDescription,
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll := got.(collectionmodels.Collection)
				assert.Equal(tt, newDescription, coll.Description)
				assert.True(tt, coll.UpdatedAt.After(now))
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				collections: nil,
			},
			args: args{
				key:         key,
				description: newDescription,
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
			storage := datastorage.NewOf(
				nil,
				tt.fields.collections,
				nil,
			)

			err := storage.UpdateCollection(context.Background(), tt.args.key, tt.args.description)
			tt.wantErr(t, err)
		})
	}
}
