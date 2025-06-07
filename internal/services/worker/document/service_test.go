package documentsrv_test

import (
	"context"
	"testing"
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
	"github.com/10Narratives/distgo-db/internal/services/worker/document/mocks"
	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	var (
		collection string         = "users"
		documentID uuid.UUID      = uuid.MustParse("287dcccf-3fb7-44cf-9832-f2866d24d6e1")
		content    map[string]any = map[string]any{
			"fullname": "Ivan Petrov",
			"email":    "example@gmail.com",
		}
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentStorage)
	}
	type args struct {
		ctx        context.Context
		collection string
		content    map[string]any
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful document creation",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Set", mock.Anything, collection, mock.Anything, content)
					m.On("Get", mock.Anything, collection, mock.Anything).
						Return(documentmodels.Document{
							ID:      documentID,
							Content: content,
						}, nil)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				content:    content,
			},
			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, content["fullname"], document.Content["fullname"])
				assert.Equal(t, content["email"], document.Content["email"])
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentStorage(t)
			tt.fields.mockSetup(mock)

			service := documentsrv.New(mock)

			doc, err := service.Create(tt.args.ctx, tt.args.collection, tt.args.content)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)

			mock.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	t.Parallel()

	var (
		collection string         = "users"
		documentID uuid.UUID      = uuid.MustParse("287dcccf-3fb7-44cf-9832-f2866d24d6e1")
		content    map[string]any = map[string]any{
			"fullname": "Ivan Petrov",
			"email":    "example@gmail.com",
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentStorage)
	}

	type args struct {
		ctx        context.Context
		collection string
		documentID string
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
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Get", mock.Anything, collection, documentID).
						Return(documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						}, nil)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, content, document.Content)
				assert.Equal(t, createdAt, document.CreatedAt)
				assert.Equal(t, updatedAt, document.UpdatedAt)
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Get", mock.Anything, collection, documentID).
						Return(documentmodels.Document{}, documentstore.ErrCollectionNotFound)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
			},
		},
		{
			name: "document not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Get", mock.Anything, collection, documentID).
						Return(documentmodels.Document{}, documentstore.ErrDocumentNotFound)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrDocumentNotFound.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentStorage(t)
			tt.fields.mockSetup(mock)

			service := documentsrv.New(mock)

			doc, err := service.Get(tt.args.ctx, tt.args.collection, tt.args.documentID)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)

			mock.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()

	var (
		id1        uuid.UUID      = uuid.New()
		id2        uuid.UUID      = uuid.New()
		collection string         = "users"
		content    map[string]any = map[string]any{
			"fullname": "User Fullname",
			"email":    "user_email@gmail.com",
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentStorage)
	}
	type args struct {
		ctx        context.Context
		collection string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful list",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("List", mock.Anything, collection).
						Return([]documentmodels.Document{
							documentmodels.Document{
								ID:        id1,
								Content:   content,
								CreatedAt: createdAt,
								UpdatedAt: updatedAt,
							},
							documentmodels.Document{
								ID:        id2,
								Content:   content,
								CreatedAt: createdAt,
								UpdatedAt: updatedAt,
							},
						}, nil)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				list, ok := got.([]documentmodels.Document)
				require.True(t, ok)

				assert.Len(t, list, 2)
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("List", mock.Anything, "collection").
						Return([]documentmodels.Document{}, documentstore.ErrCollectionNotFound)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "collection",
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentStorage(t)
			tt.fields.mockSetup(mock)

			service := documentsrv.New(mock)
			docs, err := service.List(tt.args.ctx, tt.args.collection)

			tt.wantVal(t, docs)
			tt.wantErr(t, err)

			mock.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	t.Parallel()

	var (
		collection string    = "users"
		documentID uuid.UUID = uuid.New()
	)

	type fields struct {
		mockSetup func(m *mocks.DocumentStorage)
	}
	type args struct {
		ctx        context.Context
		collection string
		documentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful deletion",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Delete", mock.Anything, collection, documentID).Return(nil)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Delete", mock.Anything, collection, documentID).Return(documentstore.ErrCollectionNotFound)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
			},
		},
		{
			name: "document not found",
			fields: fields{
				mockSetup: func(m *mocks.DocumentStorage) {
					m.On("Delete", mock.Anything, collection, documentID).Return(documentstore.ErrDocumentNotFound)
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID.String(),
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrDocumentNotFound.Error())
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewDocumentStorage(t)
			tt.fields.mockSetup(mock)

			service := documentsrv.New(mock)
			err := service.Delete(tt.args.ctx, tt.args.collection, tt.args.documentID)

			tt.wantErr(t, err)

			mock.AssertExpectations(t)
		})
	}
}
