package documentgrpc_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/document/mocks"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestServerAPI_CreateDocument(t *testing.T) {
	t.Parallel()

	const (
		parent       = "databases/db/collections/coll1"
		documentID   = "doc1"
		value        = `{"key":"value"}`
		expectedName = parent + "/documents/" + documentID
	)

	type fields struct {
		setupDocumentServiceMock func(m *mocks.DocumentService)
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
					m.On("CreateDocument", mock.Anything, parent, documentID, value).
						Return(documentmodels.Document{
							Name:  expectedName,
							Value: json.RawMessage(value),
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     parent,
					DocumentId: documentID,
					Document: &dbv1.Document{
						Value: value,
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc, ok := got.(*dbv1.Document)
				require.True(tt, ok)
				assert.Equal(tt, expectedName, doc.GetName())
				assert.Equal(tt, value, doc.GetValue())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     "",
					DocumentId: documentID,
					Document: &dbv1.Document{
						Value: value,
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
					m.On("CreateDocument", mock.Anything, parent, documentID, value).
						Return(documentmodels.Document{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent:     parent,
					DocumentId: documentID,
					Document: &dbv1.Document{
						Value: value,
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
			svc := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(svc)
			server := documentgrpc.New(svc)
			resp, err := server.CreateDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_GetDocument(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1/documents/doc1"
	const value = `{"key":"value"}`

	type fields struct {
		setupDocumentServiceMock func(m *mocks.DocumentService)
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
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Document", mock.Anything, name).
						Return(documentmodels.Document{
							Name:  name,
							Value: json.RawMessage(value),
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
				require.True(tt, ok)
				assert.Equal(tt, name, doc.GetName())
				assert.Equal(tt, value, doc.GetValue())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
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
			name: "internal error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Document", mock.Anything, name).
						Return(documentmodels.Document{}, errors.New("internal error"))
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
			svc := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(svc)
			server := documentgrpc.New(svc)
			resp, err := server.GetDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_UpdateDocument(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1/documents/doc1"
	const value = `{"key":"updated_value"}`

	type fields struct {
		setupDocumentServiceMock func(m *mocks.DocumentService)
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
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("UpdateDocument", mock.Anything,
						documentmodels.Document{
							Name:  name,
							Value: json.RawMessage(value),
						},
						[]string{"value"},
					).Return(documentmodels.Document{
						Name:  name,
						Value: json.RawMessage(value),
					}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDocumentRequest{
					Document: &dbv1.Document{
						Name:  name,
						Value: value,
					},
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{"value"},
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				doc, ok := got.(*dbv1.Document)
				require.True(tt, ok)
				assert.Equal(tt, name, doc.GetName())
				assert.Equal(tt, value, doc.GetValue())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDocumentRequest{},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("UpdateDocument", mock.Anything,
						documentmodels.Document{
							Name:  name,
							Value: json.RawMessage(value),
						},
						[]string{"value"},
					).Return(documentmodels.Document{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDocumentRequest{
					Document: &dbv1.Document{
						Name:  name,
						Value: value,
					},
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{"value"},
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
			svc := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(svc)
			server := documentgrpc.New(svc)
			resp, err := server.UpdateDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_DeleteDocument(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1/documents/doc1"

	type fields struct {
		setupDocumentServiceMock func(m *mocks.DocumentService)
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
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("DeleteDocument", mock.Anything, name).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				empty, ok := got.(*emptypb.Empty)
				require.True(tt, ok)
				require.NotNil(tt, empty)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
					Name: "invalid_name",
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("DeleteDocument", mock.Anything, name).Return(errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDocumentRequest{
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
			svc := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(svc)
			server := documentgrpc.New(svc)
			resp, err := server.DeleteDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_ListDocuments(t *testing.T) {
	t.Parallel()

	const (
		parent       = "databases/db/collections/coll1"
		size   int32 = 2
		token        = "token"
		value        = `{"key":"value"}`
	)

	type fields struct {
		setupDocumentServiceMock func(m *mocks.DocumentService)
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
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Documents", mock.Anything, parent, size, token).
						Return([]documentmodels.Document{
							{
								Name:  parent + "/documents/doc1",
								Value: json.RawMessage(value),
							},
							{
								Name:  parent + "/documents/doc2",
								Value: json.RawMessage(value),
							},
						}, "", nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDocumentsRequest{
					Parent:    parent,
					PageSize:  size,
					PageToken: token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				resp, ok := got.(*dbv1.ListDocumentsResponse)
				require.True(tt, ok)
				require.Len(tt, resp.Documents, 2)
				assert.Equal(tt, parent+"/documents/doc1", resp.Documents[0].GetName())
				assert.Equal(tt, parent+"/documents/doc2", resp.Documents[1].GetName())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDocumentsRequest{
					Parent:   "",
					PageSize: 10000,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDocumentServiceMock: func(m *mocks.DocumentService) {
					m.On("Documents", mock.Anything, parent, size, token).
						Return([]documentmodels.Document{}, "", errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDocumentsRequest{
					Parent:    parent,
					PageSize:  size,
					PageToken: token,
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
			svc := mocks.NewDocumentService(t)
			tt.fields.setupDocumentServiceMock(svc)
			server := documentgrpc.New(svc)
			resp, err := server.ListDocuments(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}
