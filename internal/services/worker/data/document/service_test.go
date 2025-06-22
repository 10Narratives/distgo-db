package documentsrv_test

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"
// 	"time"

// 	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
// 	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
// 	"github.com/10Narratives/distgo-db/internal/services/worker/data/document/mocks"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	name    string = "databases/test_db/collections/users/documents/1"
// 	content        = map[string]interface{}{
// 		"field1": "value1",
// 		"field2": 42,
// 	}
// 	contentJSON, _           = json.Marshal(content)
// 	createdAt      time.Time = time.Now()
// 	updatedAt      time.Time = time.Now()
// 	key                      = documentmodels.NewKey(name)
// 	document                 = documentmodels.Document{
// 		Name:      name,
// 		ID:        "1",
// 		Value:     contentJSON,
// 		CreatedAt: createdAt,
// 		UpdatedAt: updatedAt,
// 		Parent:    "databases/test_db/collections/users",
// 	}
// 	errInternalDocumentStorage = errors.New("internal document storage error")

// 	parent string = "databases/test_db/collections/users"

// 	doc1Content, _ = json.Marshal(map[string]interface{}{
// 		"field1": "value1",
// 		"field2": 42,
// 	})
// 	doc2Content, _ = json.Marshal(map[string]interface{}{
// 		"field3": "value3",
// 		"field4": 84,
// 	})
// 	doc3Content, _ = json.Marshal(map[string]interface{}{
// 		"field5": "value5",
// 	})

// 	allDocs = []documentmodels.Document{
// 		{
// 			Name:      parent + "/documents/doc1",
// 			ID:        "doc1",
// 			Value:     doc1Content,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			Parent:    parent,
// 		},
// 		{
// 			Name:      parent + "/documents/doc2",
// 			ID:        "doc2",
// 			Value:     doc2Content,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			Parent:    parent,
// 		},
// 		{
// 			Name:      parent + "/documents/doc3",
// 			ID:        "doc3",
// 			Value:     doc3Content,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 			Parent:    parent,
// 		},
// 	}
// 	parentKey = collectionmodels.NewKey(parent)
// )

// type fields struct {
// 	setupDocumentStorageMock func(m *mocks.DocumentStorage)
// 	setupWALServiceMock      func(m *mocks.WALService)
// }

// func TestService_Document(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		ctx  context.Context
// 		name string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful execution",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Document", mock.Anything, key).Return(document, nil)
// 				},
// 			},
// 			args: args{
// 				ctx:  context.Background(),
// 				name: name,
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				doc, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Equal(t, document.Name, doc.Name)
// 				assert.Equal(t, document.ID, doc.ID)
// 				assert.JSONEq(t, string(document.Value), string(doc.Value)) // Сравнение JSON
// 				assert.WithinDuration(t, document.CreatedAt, doc.CreatedAt, time.Second)
// 				assert.WithinDuration(t, document.UpdatedAt, doc.UpdatedAt, time.Second)
// 				assert.Equal(t, document.Parent, doc.Parent)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "document storage error",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Document", mock.Anything, key).
// 						Return(documentmodels.Document{}, errInternalDocumentStorage)
// 				},
// 			},
// 			args: args{
// 				ctx:  context.Background(),
// 				name: name,
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, errInternalDocumentStorage.Error())
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			documentStorageMock := mocks.NewDocumentStorage(t)
// 			tt.fields.setupDocumentStorageMock(documentStorageMock)

// 			service := documentsrv.New(documentStorageMock, nil)

// 			doc, err := service.Document(tt.args.ctx, tt.args.name)

// 			tt.wantVal(t, doc)
// 			tt.wantErr(t, err)

// 			documentStorageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_Documents(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		ctx    context.Context
// 		parent string
// 		size   int32
// 		token  string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful execution with pagination",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Documents", mock.Anything, parentKey).Return(allDocs, nil)
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				parent: parent,
// 				size:   2,
// 				token:  "",
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				docs, ok := got.([]documentmodels.Document)
// 				require.True(t, ok)

// 				expected := allDocs[:2]
// 				assert.Equal(t, expected, docs)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "empty documents list",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Documents", mock.Anything, parentKey).Return([]documentmodels.Document{}, nil)
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				parent: parent,
// 				size:   10,
// 				token:  "",
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				docs, ok := got.([]documentmodels.Document)
// 				require.True(t, ok)
// 				assert.Empty(t, docs)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "pagination with token",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Documents", mock.Anything, parentKey).Return(allDocs, nil)
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				parent: parent,
// 				size:   1,
// 				token:  parent + "/documents/doc1",
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				docs, ok := got.([]documentmodels.Document)
// 				require.True(t, ok)

// 				expected := []documentmodels.Document{allDocs[1]}
// 				assert.Equal(t, expected, docs)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "storage error",
// 			fields: fields{
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("Documents", mock.Anything, parentKey).Return(nil, errInternalDocumentStorage)
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				parent: parent,
// 				size:   10,
// 				token:  "",
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, errInternalDocumentStorage.Error())
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			documentStorageMock := mocks.NewDocumentStorage(t)
// 			tt.fields.setupDocumentStorageMock(documentStorageMock)

// 			service := documentsrv.New(documentStorageMock, nil)

// 			docs, nextPageToken, err := service.Documents(tt.args.ctx, tt.args.parent, tt.args.size, tt.args.token)

// 			tt.wantVal(t, docs)
// 			tt.wantErr(t, err)

// 			if err == nil && len(docs) > 0 && nextPageToken != "" {
// 				assert.Equal(t, tt.args.size, int32(len(docs)))
// 			}

// 			documentStorageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_CreateDocument(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		ctx        context.Context
// 		parent     string
// 		documentID string
// 		value      string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful creation",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("CreateDocument", mock.Anything, mock.Anything, "test-value").
// 						Return(documentmodels.Document{
// 							Name:      "databases/test_db/collections/users/documents/doc1",
// 							ID:        "doc1",
// 							Value:     []byte("test-value"),
// 							CreatedAt: time.Now(),
// 							UpdatedAt: time.Now(),
// 							Parent:    "databases/test_db/collections/users",
// 						}, nil)
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				parent:     "databases/test_db/collections/users",
// 				documentID: "doc1",
// 				value:      "test-value",
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				doc, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Equal(t, "databases/test_db/collections/users/documents/doc1", doc.Name)
// 				//assert.JSONEq(t, `"test-value"`, string(doc.Value))
// 				assert.WithinDuration(t, time.Now(), doc.CreatedAt, time.Second)
// 				assert.WithinDuration(t, time.Now(), doc.UpdatedAt, time.Second)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "WAL service error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
// 						Return(errors.New("WAL error"))
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				parent:     "databases/test_db/collections/users",
// 				documentID: "doc1",
// 				value:      "test-value",
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "failed to create WAL entry: WAL error")
// 			},
// 		},
// 		{
// 			name: "storage error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("CreateDocument", mock.Anything, mock.Anything, "test-value").
// 						Return(documentmodels.Document{}, errors.New("storage error"))
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				parent:     "databases/test_db/collections/users",
// 				documentID: "doc1",
// 				value:      "test-value",
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "storage error")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			walServiceMock := mocks.NewWALService(t)
// 			tt.fields.setupWALServiceMock(walServiceMock)

// 			documentStorageMock := mocks.NewDocumentStorage(t)
// 			tt.fields.setupDocumentStorageMock(documentStorageMock)

// 			service := documentsrv.New(documentStorageMock, walServiceMock)

// 			doc, err := service.CreateDocument(tt.args.ctx, tt.args.parent, tt.args.documentID, tt.args.value)

// 			tt.wantVal(t, doc)
// 			tt.wantErr(t, err)

// 			walServiceMock.AssertExpectations(t)
// 			documentStorageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_DeleteDocument(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		ctx  context.Context
// 		name string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful deletion",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeDelete, mock.Anything, (*documentmodels.Document)(nil)).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("DeleteDocument", mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 			},
// 			args: args{
// 				ctx:  context.Background(),
// 				name: "databases/test_db/collections/users/documents/doc1",
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "WAL service error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeDelete, mock.Anything, (*documentmodels.Document)(nil)).
// 						Return(errors.New("WAL error"))
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {},
// 			},
// 			args: args{
// 				ctx:  context.Background(),
// 				name: "databases/test_db/collections/users/documents/doc1",
// 			},
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "failed to create WAL entry: WAL error")
// 			},
// 		},
// 		{
// 			name: "storage error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeDelete, mock.Anything, (*documentmodels.Document)(nil)).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("DeleteDocument", mock.Anything, mock.Anything).
// 						Return(errors.New("storage error"))
// 				},
// 			},
// 			args: args{
// 				ctx:  context.Background(),
// 				name: "databases/test_db/collections/users/documents/doc1",
// 			},
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "storage error")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			walServiceMock := mocks.NewWALService(t)
// 			tt.fields.setupWALServiceMock(walServiceMock)

// 			documentStorageMock := mocks.NewDocumentStorage(t)
// 			tt.fields.setupDocumentStorageMock(documentStorageMock)

// 			service := documentsrv.New(documentStorageMock, walServiceMock)

// 			err := service.DeleteDocument(tt.args.ctx, tt.args.name)

// 			tt.wantErr(t, err)

// 			walServiceMock.AssertExpectations(t)
// 			documentStorageMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_UpdateDocument(t *testing.T) {
// 	t.Parallel()

// 	type args struct {
// 		ctx      context.Context
// 		document documentmodels.Document
// 		paths    []string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful update of value",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeUpdate, mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("UpdateDocument", mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				document: documentmodels.Document{
// 					Name:      "databases/test_db/collections/users/documents/doc1",
// 					ID:        "doc1",
// 					Value:     []byte("new-value"),
// 					CreatedAt: time.Now(),
// 					UpdatedAt: time.Now(),
// 					Parent:    "databases/test_db/collections/users",
// 				},
// 				paths: []string{"value"},
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				doc, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				//assert.JSONEq(t, `"new-value"`, string(doc.Value))
// 				assert.WithinDuration(t, time.Now(), doc.UpdatedAt, time.Second)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "WAL service error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeUpdate, mock.Anything, mock.Anything).
// 						Return(errors.New("WAL error"))
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				document: documentmodels.Document{
// 					Name:      "databases/test_db/collections/users/documents/doc1",
// 					ID:        "doc1",
// 					Value:     []byte("new-value"),
// 					CreatedAt: time.Now(),
// 					UpdatedAt: time.Now(),
// 					Parent:    "databases/test_db/collections/users",
// 				},
// 				paths: []string{"value"},
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "failed to create WAL entry: WAL error")
// 			},
// 		},
// 		{
// 			name: "storage error",
// 			fields: fields{
// 				setupWALServiceMock: func(m *mocks.WALService) {
// 					m.On("CreateDocumentEntry", mock.Anything, commonmodels.MutationTypeUpdate, mock.Anything, mock.Anything).
// 						Return(nil)
// 				},
// 				setupDocumentStorageMock: func(m *mocks.DocumentStorage) {
// 					m.On("UpdateDocument", mock.Anything, mock.Anything).
// 						Return(errors.New("storage error"))
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				document: documentmodels.Document{
// 					Name:      "databases/test_db/collections/users/documents/doc1",
// 					ID:        "doc1",
// 					Value:     []byte("new-value"),
// 					CreatedAt: time.Now(),
// 					UpdatedAt: time.Now(),
// 					Parent:    "databases/test_db/collections/users",
// 				},
// 				paths: []string{"value"},
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, "storage error")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			walServiceMock := mocks.NewWALService(t)
// 			tt.fields.setupWALServiceMock(walServiceMock)

// 			documentStorageMock := mocks.NewDocumentStorage(t)
// 			tt.fields.setupDocumentStorageMock(documentStorageMock)

// 			service := documentsrv.New(documentStorageMock, walServiceMock)

// 			doc, err := service.UpdateDocument(tt.args.ctx, tt.args.document, tt.args.paths)

// 			tt.wantVal(t, doc)
// 			tt.wantErr(t, err)

// 			walServiceMock.AssertExpectations(t)
// 			documentStorageMock.AssertExpectations(t)
// 		})
// 	}
// }
