package documentstore_test

import (
	"context"
	"testing"

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
