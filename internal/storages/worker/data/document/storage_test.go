package documentstore_test

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
// 	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestStorage_Delete(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		collection string         = "users"
// 		documentID uuid.UUID      = uuid.New()
// 		content    map[string]any = map[string]any{
// 			"name": "Peter",
// 			"age":  100,
// 		}
// 		createdAt time.Time = time.Now()
// 		updatedAt time.Time = time.Now()
// 	)

// 	type fields struct {
// 		collections map[string]*documentstore.Collection
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		collection string
// 		documentID uuid.UUID
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful execution",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: documentID,
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "collection not found",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: "collection",
// 				documentID: documentID,
// 			},
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
// 			},
// 		},
// 		{
// 			name: "document not found",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: uuid.New(),
// 			},
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrDocumentNotFound.Error())
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storage := documentstore.NewStorageOf(tt.fields.collections)
// 			err := storage.Delete(tt.args.ctx, tt.args.collection, tt.args.documentID)

// 			tt.wantErr(t, err)
// 		})
// 	}
// }

// func TestStorage_Get(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		collection string         = "users"
// 		documentID uuid.UUID      = uuid.New()
// 		content    map[string]any = map[string]any{
// 			"name": "Peter",
// 			"age":  100,
// 		}
// 		createdAt time.Time = time.Now()
// 		updatedAt time.Time = time.Now()
// 	)

// 	type fields struct {
// 		collections map[string]*documentstore.Collection
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		collection string
// 		documentID uuid.UUID
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
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: documentID,
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				document, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Equal(t, documentID, document.ID)
// 				assert.Equal(t, content, document.Content)
// 				assert.Equal(t, createdAt, document.CreateTime)
// 				assert.Equal(t, updatedAt, document.UpdateTime)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "collection not found",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: "collection",
// 				documentID: documentID,
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
// 			},
// 		},
// 		{
// 			name: "document not found",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: uuid.New(),
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrDocumentNotFound.Error())
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storage := documentstore.NewStorageOf(tt.fields.collections)
// 			doc, err := storage.Get(tt.args.ctx, tt.args.collection, tt.args.documentID)

// 			tt.wantVal(t, doc)
// 			tt.wantErr(t, err)
// 		})
// 	}
// }

// func TestStorage_List(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		collection string         = "users"
// 		documentID uuid.UUID      = uuid.New()
// 		content    map[string]any = map[string]any{
// 			"name": "Peter",
// 			"age":  100,
// 		}
// 		createdAt time.Time = time.Now()
// 		updatedAt time.Time = time.Now()
// 	)

// 	type fields struct {
// 		collections map[string]*documentstore.Collection
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		collection string
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
// 			fields: fields{collections: map[string]*documentstore.Collection{
// 				collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 					documentID: documentmodels.Document{
// 						ID:         documentID,
// 						Content:    content,
// 						CreateTime: createdAt,
// 						UpdateTime: updatedAt,
// 					},
// 				}),
// 			}},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				listed, ok := got.([]documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Len(t, listed, 1)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "collection not found",
// 			fields: fields{collections: map[string]*documentstore.Collection{
// 				collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 					documentID: documentmodels.Document{
// 						ID:         documentID,
// 						Content:    content,
// 						CreateTime: createdAt,
// 						UpdateTime: updatedAt,
// 					},
// 				}),
// 			}},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: "collection",
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storage := documentstore.NewStorageOf(tt.fields.collections)
// 			listed, err := storage.List(tt.args.ctx, tt.args.collection)

// 			tt.wantVal(t, listed)
// 			tt.wantErr(t, err)
// 		})
// 	}
// }

// func TestStorage_Replace(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		collection      string    = "users"
// 		documentID      uuid.UUID = uuid.New()
// 		content                   = map[string]any{"name": "Peter", "age": 100}
// 		replacedContent           = map[string]any{"name": "Peter", "age": 10}
// 		createdAt       time.Time = time.Now()
// 		updatedAt       time.Time = time.Now()
// 	)

// 	type fields struct {
// 		collections map[string]*documentstore.Collection
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		collection string
// 		documentID uuid.UUID
// 		content    map[string]any
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
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: documentID,
// 				content:    replacedContent,
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				document, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Equal(t, documentID, document.ID)
// 				assert.Equal(t, replacedContent, document.Content)
// 			},
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "collection not found",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{
// 						documentID: documentmodels.Document{
// 							ID:         documentID,
// 							Content:    content,
// 							CreateTime: createdAt,
// 							UpdateTime: updatedAt,
// 						},
// 					}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: "collection",
// 				documentID: documentID,
// 				content:    replacedContent,
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
// 				assert.EqualError(t, err, documentstore.ErrCollectionNotFound.Error())
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storage := documentstore.NewStorageOf(tt.fields.collections)
// 			document, err := storage.Replace(tt.args.ctx, tt.args.collection, tt.args.documentID, tt.args.content)

// 			tt.wantVal(t, document)
// 			tt.wantErr(t, err)
// 		})
// 	}
// }

// func TestStorage_Set(t *testing.T) {
// 	t.Parallel()

// 	var (
// 		collection string    = "users"
// 		documentID uuid.UUID = uuid.New()
// 		content              = map[string]any{"name": "Peter", "age": 100}
// 	)

// 	type fields struct {
// 		collections map[string]*documentstore.Collection
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		collection string
// 		documentID uuid.UUID
// 		content    map[string]any
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 	}{
// 		{
// 			name: "successful set",
// 			fields: fields{
// 				collections: map[string]*documentstore.Collection{
// 					collection: documentstore.NewCollectionOf(map[uuid.UUID]documentmodels.Document{}),
// 				},
// 			},
// 			args: args{
// 				ctx:        context.Background(),
// 				collection: collection,
// 				documentID: documentID,
// 				content:    content,
// 			},
// 			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
// 				document, ok := got.(documentmodels.Document)
// 				require.True(t, ok)

// 				assert.Equal(t, documentID, document.ID)
// 				assert.Equal(t, content, document.Content)
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			storage := documentstore.NewStorageOf(tt.fields.collections)
// 			storage.Set(tt.args.ctx, tt.args.collection, tt.args.documentID, tt.args.content)
// 			document, _ := storage.Get(tt.args.ctx, tt.args.collection, tt.args.documentID)

// 			tt.wantVal(t, document)
// 		})
// 	}
// }
