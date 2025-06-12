package documentgrpc_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/document/mocks"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServerAPI_CreateDocument(t *testing.T) {
	t.Parallel()

	var (
		parent     string    = "databases/mdb/collections/pillars"
		documentID string    = "host_1012"
		value      string    = "{}"
		createdAt  time.Time = time.Now().UTC()
		updatedAt  time.Time = time.Now().UTC()
	)

	type fields struct {
		setupDocumentServiceMock   func(m *mocks.DocumentService)
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.CreateDocumentRequest
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
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Create", mock.Anything, parent, documentID, value).
						Return(documentmodels.Document{
							Name:      strings.Join([]string{parent, documentID}, "/"),
							ID:        documentID,
							Value:     json.RawMessage(value),
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						}, nil)
				},
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, parent).
						Return(collectionmodels.Collection{}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     parent,
					DocumentId: documentID,
					Document: &dbv1.Document{
						Name:      "databases/mdb/collections/pillars/host#1012",
						Id:        documentID,
						Value:     value,
						CreatedAt: timestamppb.New(createdAt),
						UpdatedAt: timestamppb.New(updatedAt),
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(*dbv1.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.Id)
				assert.Equal(t, value, document.Value)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock:   func(m *mocks.DocumentService) {},
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     "parent",
					DocumentId: documentID,
					Document: &dbv1.Document{
						Name:      "databases/mdb/collections/pillars/host#1012",
						Id:        documentID,
						Value:     value,
						CreatedAt: timestamppb.New(createdAt),
						UpdatedAt: timestamppb.New(updatedAt),
					},
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "collection not found",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, parent).
						Return(collectionmodels.Collection{}, errors.New("collection is not found"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     parent,
					DocumentId: documentID,
					Document: &dbv1.Document{
						Name:      "databases/mdb/collections/pillars/host#1012",
						Id:        documentID,
						Value:     value,
						CreatedAt: timestamppb.New(createdAt),
						UpdatedAt: timestamppb.New(updatedAt),
					},
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Create", mock.Anything, parent, documentID, value).
						Return(documentmodels.Document{}, errors.New("internal"))
				},
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, parent).
						Return(collectionmodels.Collection{}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     parent,
					DocumentId: documentID,
					Document: &dbv1.Document{
						Name:      "databases/mdb/collections/pillars/host#1012",
						Id:        documentID,
						Value:     value,
						CreatedAt: timestamppb.New(createdAt),
						UpdatedAt: timestamppb.New(updatedAt),
					},
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

			documentSrv := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(documentSrv)

			collectionSrv := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(collectionSrv)

			serverAPI := documentgrpc.New(documentSrv, collectionSrv)

			resp, err := serverAPI.CreateDocument(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			documentSrv.AssertExpectations(t)
			collectionSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_DeleteDocument(t *testing.T) {
	t.Parallel()

	var (
		parent     string = "databases/mdb/collections/pillars"
		collection string = "pillars"
		documentID string = "host_1012"
	)

	type fields struct {
		setupDocumentServiceMock   func(m *mocks.DocumentService)
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.DeleteDocumentRequest
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
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Delete", mock.Anything, collection, documentID).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: parent + "/documents/" + documentID,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
				setupDocumentServiceMock:   func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: "parent" + "/documents/" + documentID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "collection not found",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{}, errors.New("collection not found"))
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: parent + "/documents/" + documentID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "document not found",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Delete", mock.Anything, collection, documentID).
						Return(errors.New("document not found"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: parent + "/documents/" + documentID,
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

			documentSrv := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(documentSrv)

			collectionSrv := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(collectionSrv)

			serverAPI := documentgrpc.New(documentSrv, collectionSrv)

			resp, err := serverAPI.DeleteDocument(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			documentSrv.AssertExpectations(t)
			collectionSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_GetDocument(t *testing.T) {
	t.Parallel()

	var (
		parent     string    = "databases/mdb/collections/pillars"
		collection string    = "pillars"
		documentID string    = "host_1012"
		name       string    = parent + "/documents/" + documentID
		value      string    = "{}"
		createdAt  time.Time = time.Now().UTC()
		updatedAt  time.Time = createdAt
	)

	type fields struct {
		setupDocumentServiceMock   func(m *mocks.DocumentService)
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.GetDocumentRequest
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
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{
							Name: collection,
						}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Document", mock.Anything, collection, documentID).
						Return(documentmodels.Document{
							Name:      name,
							ID:        documentID,
							Value:     json.RawMessage(value),
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Name: name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc, ok := got.(*dbv1.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, doc.Id)
				assert.Equal(t, value, doc.Value)
				assert.WithinDuration(t, createdAt, doc.CreatedAt.AsTime(), time.Second)
				assert.WithinDuration(t, updatedAt, doc.UpdatedAt.AsTime(), time.Second)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
				setupDocumentServiceMock:   func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Name: "invalid_name",
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "document not found",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{
							Name: collection,
						}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Document", mock.Anything, collection, documentID).
						Return(documentmodels.Document{}, errors.New("not found"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Name: name,
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

			documentSrv := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(documentSrv)

			collectionSrv := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(collectionSrv)

			serverAPI := documentgrpc.New(documentSrv, collectionSrv)

			resp, err := serverAPI.GetDocument(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			documentSrv.AssertExpectations(t)
			collectionSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_ListDocuments(t *testing.T) {
	t.Parallel()

	var (
		parent     string    = "databases/mdb/collections/pillars"
		collection string    = "pillars"
		pageSize   int32     = 10
		pageToken  string    = ""
		createdAt  time.Time = time.Now().UTC()
		updatedAt  time.Time = createdAt
	)

	type fields struct {
		setupDocumentServiceMock   func(m *mocks.DocumentService)
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.ListDocumentsRequest
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
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{
							Name: collection,
						}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Documents", mock.Anything, collection, pageSize, pageToken).
						Return([]documentmodels.Document{
							{
								Name:      parent + "/documents/host_1012",
								ID:        "host_1012",
								Value:     json.RawMessage("{}"),
								CreatedAt: createdAt,
								UpdatedAt: updatedAt,
							},
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDocumentsRequest{
					Parent:    parent,
					PageSize:  pageSize,
					PageToken: pageToken,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				resp, ok := got.(*dbv1.ListDocumentsResponse)
				require.True(t, ok)
				assert.Len(t, resp.Documents, 1)
				assert.Equal(t, "host_1012", resp.Documents[0].Id)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
				setupDocumentServiceMock:   func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDocumentsRequest{
					Parent: "",
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

			documentSrv := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(documentSrv)

			collectionSrv := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(collectionSrv)

			serverAPI := documentgrpc.New(documentSrv, collectionSrv)

			resp, err := serverAPI.ListDocuments(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			documentSrv.AssertExpectations(t)
			collectionSrv.AssertExpectations(t)
		})
	}
}
func TestServerAPI_UpdateDocument(t *testing.T) {
	t.Parallel()

	var (
		name       string    = "databases/mdb/collections/pillars/documents/host_1012"
		collection string    = "pillars"
		createdAt  time.Time = time.Now().UTC()
		documentID string    = "host_1012"
	)

	type fields struct {
		setupDocumentServiceMock   func(m *mocks.DocumentService)
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.UpdateDocumentRequest
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
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, collection).
						Return(collectionmodels.Collection{
							Name: collection,
						}, nil)
				},
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Update", mock.Anything, collection, documentID, mock.Anything, mock.Anything).
						Return(documentmodels.Document{
							Name:      name,
							ID:        documentID,
							Value:     json.RawMessage(`{"status":"active"}`),
							CreatedAt: createdAt,
							UpdatedAt: time.Now(),
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDocumentRequest{
					Document: &dbv1.Document{
						Name:  name,
						Value: `{"status":"active"}`,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"value"}},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc, ok := got.(*dbv1.Document)
				require.True(t, ok)
				assert.Contains(t, doc.Value, `"status":"active"`)
			},
			wantErr: require.NoError,
		},
		{
			name: "update_mask validation failed",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
				setupDocumentServiceMock:   func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDocumentRequest{
					Document: &dbv1.Document{
						Name: "invalid",
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"id"}},
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

			documentSrv := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(documentSrv)

			collectionSrv := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(collectionSrv)

			serverAPI := documentgrpc.New(documentSrv, collectionSrv)

			resp, err := serverAPI.UpdateDocument(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			documentSrv.AssertExpectations(t)
			collectionSrv.AssertExpectations(t)
		})
	}
}
