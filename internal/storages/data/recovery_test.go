package datastorage_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// 	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestStorage_Recover(t *testing.T) {
// 	t.Parallel()

// 	const (
// 		dbID   = "db"
// 		collID = "coll"
// 		docID  = "doc"
// 	)

// 	type fields struct {
// 		databases   map[databasemodels.Key]databasemodels.Database
// 		collections map[collectionmodels.Key]collectionmodels.Collection
// 		documents   map[documentmodels.Key]documentmodels.Document
// 	}

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wal     []walmodels.WALEntry
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name:   "successful recovery of database",
// 			fields: fields{},
// 			wal: []walmodels.WALEntry{
// 				{
// 					ID:        dbID,
// 					Target:    "database",
// 					Type:      commonmodels.MutationTypeCreate,
// 					NewValue:  `{"name": "databases/test_db", "display_name": "Test DB"}`,
// 					Timestamp: time.Now().UTC(),
// 				},
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
// 				storage, ok := got.(*datastorage.Storage)
// 				require.True(t, ok)

// 				key := databasemodels.NewKey("databases/test_db")
// 				database, err := storage.Database(context.Background(), key)

// 				require.NoError(t, err)

// 				assert.Equal(t, "databases/test_db", database.Name)
// 				assert.Equal(t, "Test DB", database.DisplayName)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		// {
// 		// 	name:   "successful recovery of collection",
// 		// 	fields: fields{},
// 		// 	wal: []walmodels.WALEntry{
// 		// 		{
// 		// 			ID:        collID,
// 		// 			Target:    "collection",
// 		// 			Type:      commonmodels.MutationTypeCreate,
// 		// 			NewValue:  `{"name": "databases/db/collections/coll", "description": "Test Collection"}`,
// 		// 			Timestamp: time.Now().UTC(),
// 		// 		},
// 		// 	},
// 		// 	wantVal: func(tt require.TestingT, storage *datastorage.Storage) {
// 		// 		key := documentmodels.NewKey(collID)
// 		// 		val, ok := storage.Collections.Load(key)
// 		// 		require.True(tt, ok)
// 		// 		coll := val.(documentmodels.Collection)
// 		// 		assert.Equal(tt, "databases/db/collections/coll", coll.Name)
// 		// 		assert.Equal(tt, "Test Collection", coll.Description)
// 		// 	},
// 		// 	wantErr: require.NoError,
// 		// },
// 		// {
// 		// 	name:   "successful recovery of document",
// 		// 	fields: fields{},
// 		// 	wal: []walmodels.WALEntry{
// 		// 		{
// 		// 			ID:        docID,
// 		// 			Target:    "document",
// 		// 			Type:      commonmodels.MutationTypeCreate,
// 		// 			NewValue:  `{"name": "databases/db/collections/coll/documents/doc", "value": {"key": "value"}}`,
// 		// 			Timestamp: time.Now().UTC(),
// 		// 		},
// 		// 	},
// 		// 	wantVal: func(tt require.TestingT, storage *datastorage.Storage) {
// 		// 		key := documentmodels.NewKey(docID)
// 		// 		val, ok := storage.Documents.Load(key)
// 		// 		require.True(tt, ok)
// 		// 		doc := val.(documentmodels.Document)
// 		// 		var expectedValue map[string]interface{}
// 		// 		json.Unmarshal([]byte(`{"key": "value"}`), &expectedValue)
// 		// 		assert.Equal(tt, "databases/db/collections/coll/documents/doc", doc.Name)
// 		// 		assert.Equal(tt, expectedValue, json.RawMessage(doc.Value))
// 		// 	},
// 		// 	wantErr: require.NoError,
// 		// },
// 		// {
// 		// 	name:   "unknown target in WAL entry",
// 		// 	fields: fields{},
// 		// 	wal: []walmodels.WALEntry{
// 		// 		{
// 		// 			ID:        "unknown_id",
// 		// 			Target:    "unknown_target",
// 		// 			Type:      commonmodels.MutationTypeCreate,
// 		// 			NewValue:  `{}`,
// 		// 			Timestamp: time.Now().UTC(),
// 		// 		},
// 		// 	},
// 		// 	wantVal: require.Empty,
// 		// 	wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 		// 		require.Error(tt, err)
// 		// 		assert.Contains(tt, err.Error(), "unknown target")
// 		// 	},
// 		// },
// 		// {
// 		// 	name:   "invalid JSON in WAL entry",
// 		// 	fields: fields{},
// 		// 	wal: []walmodels.WALEntry{
// 		// 		{
// 		// 			ID:        dbID,
// 		// 			Target:    "database",
// 		// 			Type:      commonmodels.MutationTypeCreate,
// 		// 			NewValue:  `invalid_json`,
// 		// 			Timestamp: time.Now().UTC(),
// 		// 		},
// 		// 	},
// 		// 	wantVal: require.Empty,
// 		// 	wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 		// 		require.Error(tt, err)
// 		// 		assert.Contains(tt, err.Error(), "failed to unmarshal")
// 		// 	},
// 		// },
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			storage := datastorage.NewOf(
// 				tt.fields.databases,
// 				tt.fields.collections,
// 				tt.fields.documents,
// 			)
// 			err := storage.Recover(context.Background(), tt.wal)
// 			tt.wantErr(t, err)
// 			if tt.wantVal != nil {
// 				tt.wantVal(t, storage)
// 			}
// 		})
// 	}
// }
