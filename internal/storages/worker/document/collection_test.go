package documentstore_test

import (
	"context"
	"testing"
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_Create(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"name": "Name",
			"age":  120,
		}
	)

	type fields struct {
		documents map[uuid.UUID]documentmodels.Document
	}
	type args struct {
		ctx        context.Context
		documentID uuid.UUID
		content    map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
	}{
		{
			name: "successful creation",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
				content:    content,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, content, document.Content)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collection := documentstore.NewCollectionOf(tt.fields.documents)
			document, _ := collection.Create(tt.args.ctx, tt.args.documentID, tt.args.content)

			tt.wantVal(t, document)
		})
	}
}

func TestCollection_Document(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"name": "Name",
			"age":  120,
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		documents map[uuid.UUID]documentmodels.Document
	}
	type args struct {
		ctx        context.Context
		documentID uuid.UUID
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
				documents: map[uuid.UUID]documentmodels.Document{
					documentID: documentmodels.Document{
						ID:        documentID,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
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
			name: "document not found",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{
					documentID: documentmodels.Document{
						ID:        documentID,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				documentID: uuid.New(),
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

			collection := documentstore.NewCollectionOf(tt.fields.documents)
			doc, err := collection.Document(tt.args.ctx, tt.args.documentID)

			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}

func TestCollection_Documents(t *testing.T) {
	t.Parallel()

	var (
		documentID1 uuid.UUID      = uuid.New()
		documentID2 uuid.UUID      = uuid.New()
		content     map[string]any = map[string]any{
			"name": "Name",
			"age":  120,
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		documents map[uuid.UUID]documentmodels.Document
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
			name: "successful execution",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{
					documentID1: documentmodels.Document{
						ID:        documentID1,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
					documentID2: documentmodels.Document{
						ID:        documentID2,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				listed, ok := got.([]documentmodels.Document)
				require.True(t, ok)

				assert.Len(t, listed, 2)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collection := documentstore.NewCollectionOf(tt.fields.documents)
			documents, _ := collection.Documents(tt.args.ctx)

			tt.wantVal(t, documents)
		})
	}
}

func TestCollection_Replace(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"name": "Name",
			"age":  120,
		}
		replacedContent map[string]any = map[string]any{
			"fullname": "Name",
			"age":      12,
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		documents map[uuid.UUID]documentmodels.Document
	}
	type args struct {
		ctx        context.Context
		documentID uuid.UUID
		content    map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
	}{
		{
			name: "successful execution - replaced old content",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{
					documentID: documentmodels.Document{
						ID:        documentID,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
				content:    replacedContent,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, replacedContent, document.Content)
			},
		},
		{
			name: "successful execution - new document created",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
				content:    replacedContent,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, document.ID)
				assert.Equal(t, replacedContent, document.Content)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collection := documentstore.NewCollectionOf(tt.fields.documents)
			document, _ := collection.Replace(tt.args.ctx, tt.args.documentID, tt.args.content)

			tt.wantVal(t, document)
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"name": "Name",
			"age":  120,
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		documents map[uuid.UUID]documentmodels.Document
	}
	type args struct {
		ctx        context.Context
		documentID uuid.UUID
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
				documents: map[uuid.UUID]documentmodels.Document{
					documentID: documentmodels.Document{
						ID:        documentID,
						Content:   content,
						CreatedAt: createdAt,
						UpdatedAt: updatedAt,
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
			},
			wantErr: require.NoError,
		},
		{
			name: "document not found",
			fields: fields{
				documents: map[uuid.UUID]documentmodels.Document{},
			},
			args: args{
				ctx:        context.Background(),
				documentID: documentID,
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

			collection := documentstore.NewCollectionOf(tt.fields.documents)
			err := collection.Delete(tt.args.ctx, tt.args.documentID)

			tt.wantErr(t, err)
		})
	}
}
