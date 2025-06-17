package documentsrv_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	walmocks "github.com/10Narratives/distgo-db/internal/services/worker/data/common/mocks"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
	mocks "github.com/10Narratives/distgo-db/internal/services/worker/data/document/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateDocument(t *testing.T) {
	t.Parallel()

	const (
		parent     = "databases/db"
		documentID = "doc1"
		name       = parent + "/documents/" + documentID
		value      = `{"key": "value"}`
	)

	type fields struct {
		setupStorageMock func(m *mocks.DocumentStorage)
		setupWALMock     func(m *walmocks.WALStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			parent string
			id     string
			value  string
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
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("CreateDocument", mock.Anything, key, value).
						Return(documentmodels.Document{
							Name:  name,
							Value: json.RawMessage(value),
						}, nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					id     string
					value  string
				}{
					parent: parent,
					id:     documentID,
					value:  value,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc, ok := got.(documentmodels.Document)
				require.True(tt, ok)
				assert.Equal(tt, name, doc.Name)
				assert.Equal(tt, json.RawMessage(value), doc.Value)
			},
			wantErr: require.NoError,
		},
		{
			name: "storage returns error",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("CreateDocument", mock.Anything, key, value).
						Return(documentmodels.Document{}, errors.New("internal error"))
				},
				setupWALMock: func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					id     string
					value  string
				}{
					parent: parent,
					id:     documentID,
					value:  value,
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "internal error")
			},
		},
		{
			name: "WAL logging fails",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("CreateDocument", mock.Anything, key, value).
						Return(documentmodels.Document{
							Name:  name,
							Value: json.RawMessage(value),
						}, nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					id     string
					value  string
				}{
					parent: parent,
					id:     documentID,
					value:  value,
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

			store := mocks.NewDocumentStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := documentsrv.New(store, walStore)
			res, err := service.CreateDocument(
				tt.args.ctx,
				tt.args.req.parent,
				tt.args.req.id,
				tt.args.req.value,
			)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}
func TestService_Document(t *testing.T) {
	t.Parallel()
	const name = "databases/db/documents/doc1"

	type fields struct {
		setupStorageMock func(m *mocks.DocumentStorage)
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
			name: "successful get",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).
						Return(documentmodels.Document{
							Name: name,
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
				doc, ok := got.(documentmodels.Document)
				require.True(tt, ok)
				assert.Equal(tt, name, doc.Name)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).
						Return(documentmodels.Document{}, errors.New("not found"))
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
			store := mocks.NewDocumentStorage(t)
			tt.fields.setupStorageMock(store)
			service := documentsrv.New(store, nil)
			res, err := service.Document(tt.args.ctx, tt.args.req.name)
			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_Documents(t *testing.T) {
	t.Parallel()
	const (
		parent = "databases/db/collections/coll1"
		size   = int32(2)
		token  = "token"
	)

	type fields struct {
		setupStorageMock func(m *mocks.DocumentStorage)
	}
	type args struct {
		ctx context.Context
		req struct {
			parent string
			size   int32
			token  string
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
			name: "successful list",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := collectionmodels.NewKey(parent)
					m.On("Documents", mock.Anything, key).
						Return([]documentmodels.Document{
							{Name: parent + "/documents/doc1"},
							{Name: parent + "/documents/doc2"},
						})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					size   int32
					token  string
				}{
					parent: parent,
					size:   size,
					token:  token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]documentmodels.Document)
				require.True(tt, ok)
				require.Len(tt, list, 2)
				assert.Equal(tt, parent+"/documents/doc1", list[0].Name)
				assert.Equal(tt, parent+"/documents/doc2", list[1].Name)
			},
			wantErr: require.NoError,
		},
		{
			name: "with next page token",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := collectionmodels.NewKey(parent)
					m.On("Documents", mock.Anything, key).
						Return([]documentmodels.Document{
							{Name: parent + "/documents/doc1"},
							{Name: parent + "/documents/doc2"},
							{Name: parent + "/documents/doc3"},
						})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					size   int32
					token  string
				}{
					parent: parent,
					size:   size,
					token:  token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]documentmodels.Document)
				require.True(tt, ok)
				require.Len(tt, list, 2)
				assert.Equal(tt, parent+"/documents/doc1", list[0].Name)
				assert.Equal(tt, parent+"/documents/doc2", list[1].Name)
			},
			wantErr: require.NoError,
		},
		{
			name: "empty list",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := collectionmodels.NewKey(parent)
					m.On("Documents", mock.Anything, key).Return([]documentmodels.Document{})
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent string
					size   int32
					token  string
				}{
					parent: parent,
					size:   size,
					token:  token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				list, ok := got.([]documentmodels.Document)
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
			store := mocks.NewDocumentStorage(t)
			tt.fields.setupStorageMock(store)
			service := documentsrv.New(store, nil)
			list, _, err := service.Documents(tt.args.ctx, tt.args.req.parent, tt.args.req.size, tt.args.req.token)
			tt.wantVal(t, list)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_DeleteDocument(t *testing.T) {
	t.Parallel()

	const name = "databases/db/documents/doc1"
	existingDoc := documentmodels.Document{
		Name:  name,
		Value: json.RawMessage(`{"key": "value"}`),
	}

	type fields struct {
		setupStorageMock func(m *mocks.DocumentStorage)
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
			name: "successful delete",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).Return(existingDoc, nil)
					m.On("DeleteDocument", mock.Anything, key).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct{ name string }{name: name},
			},
			wantVal: require.Empty,
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).Return(documentmodels.Document{}, errors.New("not found"))
				},
				setupWALMock: func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct{ name string }{name: name},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "not found")
			},
		},
		{
			name: "WAL logging fails",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).Return(existingDoc, nil)
					m.On("DeleteDocument", mock.Anything, key).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct{ name string }{name: name},
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

			store := mocks.NewDocumentStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := documentsrv.New(store, walStore)
			err := service.DeleteDocument(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, nil)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}

func TestService_UpdateDocument(t *testing.T) {
	t.Parallel()

	const (
		name     = "databases/db/documents/doc1"
		oldValue = `{"key": "old_value"}`
		newValue = `{"key": "updated_value"}`
	)

	existingDoc := documentmodels.Document{
		Name:  name,
		Value: json.RawMessage(oldValue),
	}

	type fields struct {
		setupStorageMock func(m *mocks.DocumentStorage)
		setupWALMock     func(m *walmocks.WALStorage)
	}

	type args struct {
		ctx context.Context
		req struct {
			document documentmodels.Document
			paths    []string
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
			name: "update value",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).Return(existingDoc, nil)
					m.On("UpdateDocument", mock.Anything, mock.Anything).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					document documentmodels.Document
					paths    []string
				}{
					document: documentmodels.Document{
						Name:  name,
						Value: json.RawMessage(newValue),
					},
					paths: []string{"value"},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				updated, ok := got.(documentmodels.Document)
				require.True(tt, ok)
				assert.Equal(tt, newValue, string(updated.Value))
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field",
			fields: fields{
				setupStorageMock: func(m *mocks.DocumentStorage) {},
				setupWALMock:     func(m *walmocks.WALStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					document documentmodels.Document
					paths    []string
				}{
					document: documentmodels.Document{
						Name: name,
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
				setupStorageMock: func(m *mocks.DocumentStorage) {
					key := documentmodels.NewKey(name)
					m.On("Document", mock.Anything, key).Return(existingDoc, nil)
					m.On("UpdateDocument", mock.Anything, mock.Anything).Return(nil)
				},
				setupWALMock: func(m *walmocks.WALStorage) {
					m.On("LogEntry", mock.Anything, mock.Anything).Return(errors.New("WAL error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					document documentmodels.Document
					paths    []string
				}{
					document: documentmodels.Document{
						Name:  name,
						Value: json.RawMessage(newValue),
					},
					paths: []string{"value"},
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

			store := mocks.NewDocumentStorage(t)
			walStore := walmocks.NewWALStorage(t)
			tt.fields.setupStorageMock(store)
			tt.fields.setupWALMock(walStore)

			service := documentsrv.New(store, walStore)
			updated, err := service.UpdateDocument(tt.args.ctx, tt.args.req.document, tt.args.req.paths)

			tt.wantVal(t, updated)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
			walStore.AssertExpectations(t)
		})
	}
}
