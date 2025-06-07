package documentgrpc_test

import (
	"context"
	"testing"
	"time"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/document"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/document/mocks"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestServerAPI_CreateDocument(t *testing.T) {
	t.Parallel()

	var (
		collection string         = "projects/my-project/databases/main-db"
		documentID uuid.UUID      = uuid.MustParse("287dcccf-3fb7-44cf-9832-f2866d24d6e1")
		content    map[string]any = map[string]any{
			"fullname": "Ivan Petrov",
			"email":    "example@gmail.com",
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentService)
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
			name: "successful call",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
					m.On("Create", mock.Anything, collection, content).
						Return(documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent: collection,
					Content: func() *structpb.Struct {
						s, _ := structpb.NewStruct(content)
						return s
					}(),
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(*dbv1.Document)
				require.True(t, ok)

				assert.Equal(t, documentID.String(), document.Name)
				assert.Equal(t, content, document.Content.AsMap())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent: "users",
					Content: func() *structpb.Struct {
						s, _ := structpb.NewStruct(content)
						return s
					}(),
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid CreateDocumentRequest.Parent: value does not match regex pattern \"projects/.*/databases/.*\"")
			},
		},
		{
			name: "internal error",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
					m.On("Create", mock.Anything, collection, content).
						Return(documentmodels.Document{}, documentstore.ErrCollectionNotFound)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDocumentRequest{
					Parent: collection,
					Content: func() *structpb.Struct {
						s, _ := structpb.NewStruct(content)
						return s
					}(),
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = Internal desc = collection not found")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentService(t)
			tt.fields.mockSetup(mock)

			serverAPI := documentgrpc.New(mock)
			doc, err := serverAPI.CreateDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}

func TestServerAPI_GetDocument(t *testing.T) {
	t.Parallel()

	var (
		collection string         = "projects/my-project/databases/main-db"
		documentID uuid.UUID      = uuid.MustParse("287dcccf-3fb7-44cf-9832-f2866d24d6e1")
		content    map[string]any = map[string]any{
			"fullname": "Ivan Petrov",
			"email":    "example@gmail.com",
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentService)
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
			name: "successful get",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
					m.On("Get", mock.Anything, collection, documentID.String()).
						Return(documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Collection: collection,
					DocumentId: documentID.String(),
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(*dbv1.Document)
				require.True(t, ok)

				assert.Equal(t, documentID.String(), document.Name)
				assert.Equal(t, content, document.Content.AsMap())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error collection regexp mismatch",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Collection: "collection",
					DocumentId: documentID.String(),
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid GetDocumentRequest.Collection: value does not match regex pattern \"projects/.*/databases/.*\"")
			},
		},
		{
			name: "validation error document id regexp mismatch",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Collection: collection,
					DocumentId: "documentID.String()",
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid GetDocumentRequest.DocumentId: value must be a valid UUID | caused by: invalid uuid format")
			},
		},
		{
			name: "internal error collection not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
					m.On("Get", mock.Anything, collection, documentID.String()).
						Return(documentmodels.Document{}, documentstore.ErrCollectionNotFound)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Collection: collection,
					DocumentId: documentID.String(),
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = Internal desc = collection not found")
			},
		},
		{
			name: "internal error document not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentService) {
					m.On("Get", mock.Anything, collection, documentID.String()).
						Return(documentmodels.Document{}, documentstore.ErrDocumentNotFound)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDocumentRequest{
					Collection: collection,
					DocumentId: documentID.String(),
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "rpc error: code = Internal desc = document not found")
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentService(t)
			tt.fields.mockSetup(mock)

			serverAPI := documentgrpc.New(mock)
			doc, err := serverAPI.GetDocument(tt.args.ctx, tt.args.req)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}
