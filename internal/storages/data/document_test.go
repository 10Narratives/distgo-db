package datastorage_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_CreateDocument(t *testing.T) {
	t.Parallel()

	const (
		dbID   = "db"
		collID = "coll"
		docID  = "doc"
		value  = `{"key": "value"}`
	)

	parentKey := collectionmodels.Key{
		Database:   dbID,
		Collection: collID,
	}

	existingColl := collectionmodels.Collection{
		Name:      "databases/" + dbID + "/collections/" + collID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	type fields struct {
		databases   map[databasemodels.Key]databasemodels.Database
		collections map[collectionmodels.Key]collectionmodels.Collection
		documents   map[documentmodels.Key]documentmodels.Document
	}
	type args struct {
		key   documentmodels.Key
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful create",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{Database: dbID}: {},
				},
				collections: map[collectionmodels.Key]collectionmodels.Collection{
					parentKey: existingColl,
				},
			},
			args: args{
				key: documentmodels.Key{
					Database:   dbID,
					Collection: collID,
					Document:   docID,
				},
				value: value,
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc := got.(documentmodels.Document)
				assert.Equal(tt, "databases/"+dbID+"/collections/"+collID+"/documents/"+docID, doc.Name)
				assert.Equal(tt, json.RawMessage(value), doc.Value)
				assert.NotZero(tt, doc.CreatedAt)
				assert.Equal(tt, doc.CreatedAt, doc.UpdatedAt)
			},
			wantErr: require.NoError,
		},
		{
			name: "missing collection",
			fields: fields{
				databases: map[databasemodels.Key]databasemodels.Database{
					{Database: dbID}: {},
				},
			},
			args: args{
				key: documentmodels.Key{
					Database:   dbID,
					Collection: collID,
					Document:   docID,
				},
				value: value,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
		{
			name: "invalid key",
			fields: fields{
				collections: map[collectionmodels.Key]collectionmodels.Collection{
					parentKey: existingColl,
				},
			},
			args: args{
				key: documentmodels.Key{
					Collection: collID,
					Document:   docID,
				},
				value: value,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "invalid key")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := datastorage.NewOf(
				tt.fields.databases,
				tt.fields.collections,
				tt.fields.documents,
			)

			doc, err := storage.CreateDocument(context.Background(), tt.args.key, tt.args.value)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_DeleteDocument(t *testing.T) {
	t.Parallel()

	key := documentmodels.Key{
		Database:   "db",
		Collection: "coll",
		Document:   "doc",
	}

	existingDoc := documentmodels.Document{
		Name: "databases/db/collections/coll/documents/doc",
	}

	type fields struct {
		documents map[documentmodels.Key]documentmodels.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    documentmodels.Key
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success delete",
			fields: fields{
				documents: map[documentmodels.Key]documentmodels.Document{
					key: existingDoc,
				},
			},
			args: key,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err)
			},
		},
		{
			name: "not found",
			fields: fields{
				documents: nil,
			},
			args: key,
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
			storage := datastorage.NewOf(nil, nil, tt.fields.documents)
			err := storage.DeleteDocument(context.Background(), tt.args)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Document(t *testing.T) {
	t.Parallel()

	key := documentmodels.Key{
		Database:   "db",
		Collection: "coll",
		Document:   "doc",
	}

	expected := documentmodels.Document{
		Name: "databases/db/collections/coll/documents/doc",
	}

	type fields struct {
		documents map[documentmodels.Key]documentmodels.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    documentmodels.Key
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "success get",
			fields: fields{
				documents: map[documentmodels.Key]documentmodels.Document{
					key: expected,
				},
			},
			args: key,
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc := got.(documentmodels.Document)
				assert.Equal(tt, expected, doc)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				documents: nil,
			},
			args:    key,
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
			storage := datastorage.NewOf(nil, nil, tt.fields.documents)
			doc, err := storage.Document(context.Background(), tt.args)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Documents(t *testing.T) {
	t.Parallel()

	parentKey := collectionmodels.Key{
		Database:   "db",
		Collection: "coll",
	}

	docs := []struct {
		key      documentmodels.Key
		document documentmodels.Document
	}{
		{
			key: documentmodels.Key{
				Database:   "db",
				Collection: "coll",
				Document:   "doc1",
			},
			document: documentmodels.Document{
				Name: "databases/db/collections/coll/documents/doc1",
			},
		},
		{
			key: documentmodels.Key{
				Database:   "db",
				Collection: "coll",
				Document:   "doc2",
			},
			document: documentmodels.Document{
				Name: "databases/db/collections/coll/documents/doc2",
			},
		},
		{
			key: documentmodels.Key{
				Database:   "db",
				Collection: "other_coll",
				Document:   "doc3",
			},
			document: documentmodels.Document{
				Name: "databases/db/collections/other_coll/documents/doc3",
			},
		},
	}

	type fields struct {
		documents map[documentmodels.Key]documentmodels.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    collectionmodels.Key
		wantVal require.ValueAssertionFunc
	}{
		{
			name: "list documents in collection",
			fields: fields{
				documents: map[documentmodels.Key]documentmodels.Document{
					docs[0].key: docs[0].document,
					docs[1].key: docs[1].document,
					docs[2].key: docs[2].document,
				},
			},
			args: parentKey,
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list := got.([]documentmodels.Document)
				require.Len(tt, list, 2)
				assert.Equal(tt, docs[0].document.Name, list[1].Name)
				assert.Equal(tt, docs[1].document.Name, list[0].Name)
			},
		},
		{
			name: "empty list",
			fields: fields{
				documents: nil,
			},
			args: parentKey,
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list := got.([]documentmodels.Document)
				assert.Empty(tt, list)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			storage := datastorage.NewOf(nil, nil, tt.fields.documents)
			res := storage.Documents(context.Background(), tt.args)
			tt.wantVal(t, res)
		})
	}
}

func TestStorage_UpdateDocument(t *testing.T) {
	t.Parallel()

	key := documentmodels.Key{
		Database:   "db",
		Collection: "coll",
		Document:   "doc",
	}

	existing := documentmodels.Document{
		Name:      "databases/db/collections/coll/documents/doc",
		Value:     []byte(`{"key": "old_value"}`),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Parent:    "databases/db/collections/coll",
	}

	updated := documentmodels.Document{
		Name:      existing.Name,
		Value:     []byte(`{"key": "new_value"}`),
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Time{},
		Parent:    existing.Parent,
	}

	type fields struct {
		documents map[documentmodels.Key]documentmodels.Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    documentmodels.Document
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful update",
			fields: fields{
				documents: map[documentmodels.Key]documentmodels.Document{
					key: existing,
				},
			},
			args: updated,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.NoError(tt, err)
			},
		},
		{
			name: "document not found",
			fields: fields{
				documents: nil,
			},
			args: updated,
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
			storage := datastorage.NewOf(nil, nil, tt.fields.documents)
			err := storage.UpdateDocument(context.Background(), tt.args)
			tt.wantErr(t, err)
		})
	}
}
