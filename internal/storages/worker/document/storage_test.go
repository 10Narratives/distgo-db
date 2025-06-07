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

func TestStorage_Get(t *testing.T) {
	t.Parallel()

	var id uuid.UUID = uuid.MustParse("c2cecd16-ed51-421e-8a4c-ccfbc4e82146")
	var badID uuid.UUID = uuid.MustParse("c2cecd16-ed51-421e-8a4c-ccfbc4e80000")

	type fields struct {
		data map[string]documentstore.Collection
	}
	type args struct {
		ctx        context.Context
		collection string
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
			name: "successful get",
			fields: fields{
				data: map[string]documentstore.Collection{
					"users": documentstore.Collection{
						id: documentmodels.Document{
							ID: id,
							Content: map[string]any{
								"fullname": "User user user",
								"email":    "email email email",
							},
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "users",
				documentID: id,
			},
			wantVal: func(tt require.TestingT, got interface{}, _ ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, id, document.ID)
				assert.Equal(t, "User user user", document.Content["fullname"])
				assert.Equal(t, "email email email", document.Content["email"])
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					"users": documentstore.Collection{
						id: documentmodels.Document{
							ID: id,
							Content: map[string]any{
								"fullname": "User user user",
								"email":    "email email email",
							},
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "teachers",
				documentID: id,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
			},
		},
		{
			name: "document not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					"users": documentstore.Collection{
						id: documentmodels.Document{
							ID: id,
							Content: map[string]any{
								"fullname": "User user user",
								"email":    "email email email",
							},
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "users",
				documentID: badID,
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

			storage := documentstore.NewOf(tt.fields.data)
			doc, err := storage.Get(tt.args.ctx,
				tt.args.collection, tt.args.documentID)
			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	var id uuid.UUID = uuid.MustParse("c2cecd16-ed51-421e-8a4c-ccfbc4e82146")

	type fields struct {
		data map[string]documentstore.Collection
	}
	type args struct {
		ctx        context.Context
		collection string
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
			name: "successful set",
			fields: fields{
				data: map[string]documentstore.Collection{
					"users": documentstore.Collection{
						id: documentmodels.Document{
							ID: id,
							Content: map[string]any{
								"fullname": "User user user",
								"email":    "email email email",
							},
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "users",
				documentID: id,
				content: map[string]any{
					"fullname": "User user user",
					"email":    "email email email",
					"age":      100,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, id, document.ID)
				assert.Equal(t, "User user user", document.Content["fullname"])
				assert.Equal(t, "email email email", document.Content["email"])
				assert.Equal(t, 100, document.Content["age"])
			},
		},
		{
			name: "successful creation",
			fields: fields{
				data: map[string]documentstore.Collection{},
			},
			args: args{
				ctx:        context.Background(),
				collection: "users",
				documentID: id,
				content: map[string]any{
					"fullname": "User user user",
					"email":    "email email email",
					"age":      100,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				document, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, id, document.ID)
				assert.Equal(t, "User user user", document.Content["fullname"])
				assert.Equal(t, "email email email", document.Content["email"])
				assert.Equal(t, 100, document.Content["age"])
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := documentstore.NewOf(tt.fields.data)
			storage.Set(tt.args.ctx, tt.args.collection, tt.args.documentID, tt.args.content)

			doc, _ := storage.Get(tt.args.ctx, tt.args.collection, tt.args.documentID)
			tt.wantVal(t, doc)
		})
	}
}

func TestStorage_List(t *testing.T) {
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
		data map[string]documentstore.Collection
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
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						id1: documentmodels.Document{
							ID:        id1,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
						id2: documentmodels.Document{
							ID:        id2,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					},
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
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						id1: documentmodels.Document{
							ID:        id1,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
						id2: documentmodels.Document{
							ID:        id2,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					},
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

			storage := documentstore.NewOf(tt.fields.data)
			docs, err := storage.List(tt.args.ctx, tt.args.collection)

			tt.wantVal(t, docs)
			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	t.Parallel()

	var (
		collection string         = "users"
		documentID uuid.UUID      = uuid.New()
		content    map[string]any = map[string]any{
			"fullname": "User Fullname",
			"email":    "email@gmail.com",
		}
		createdAt time.Time = time.Now()
		updatedAt time.Time = time.Now()
	)

	type fields struct {
		data map[string]documentstore.Collection
	}
	type args struct {
		ctx        context.Context
		collection string
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
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID,
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "collection",
				documentID: documentID,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
			},
		},
		{
			name: "document not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: updatedAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: uuid.MustParse("c2cecd16-ed51-421e-8a4c-ccfbc4e80000"),
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

			storage := documentstore.NewOf(tt.fields.data)
			err := storage.Delete(tt.args.ctx, tt.args.collection, tt.args.documentID)

			tt.wantErr(t, err)
		})
	}
}

func TestStorage_Replace(t *testing.T) {
	t.Parallel()

	var (
		documentID uuid.UUID      = uuid.New()
		collection string         = "users"
		content    map[string]any = map[string]any{
			"fullname": "User Fullname",
			"email":    "user_email@gmail.com",
		}
		updatedContent map[string]any = map[string]any{
			"fullname": "User Fullname",
		}
		createdAt time.Time = time.Now()
		// updatedAt time.Time = time.Now().Add(time.Microsecond)
	)

	type fields struct {
		data map[string]documentstore.Collection
	}
	type args struct {
		ctx        context.Context
		collection string
		documentID uuid.UUID
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
			name: "successful replace",
			fields: fields{
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: createdAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: documentID,
				content:    updatedContent,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				doc, ok := got.(documentmodels.Document)
				require.True(t, ok)

				assert.Equal(t, documentID, doc.ID)
				assert.Equal(t, updatedContent, doc.Content)
				assert.Equal(t, createdAt, doc.CreatedAt)
			},
			wantErr: require.NoError,
		},
		{
			name: "collection not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: createdAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: "collection",
				documentID: documentID,
				content:    updatedContent,
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "document not found",
			fields: fields{
				data: map[string]documentstore.Collection{
					collection: documentstore.Collection{
						documentID: documentmodels.Document{
							ID:        documentID,
							Content:   content,
							CreatedAt: createdAt,
							UpdatedAt: createdAt,
						},
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				collection: collection,
				documentID: uuid.New(),
				content:    updatedContent,
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := documentstore.NewOf(tt.fields.data)

			doc, err := storage.Replace(tt.args.ctx, tt.args.collection, tt.args.documentID, tt.args.content)

			tt.wantVal(t, doc)
			tt.wantErr(t, err)
		})
	}
}
