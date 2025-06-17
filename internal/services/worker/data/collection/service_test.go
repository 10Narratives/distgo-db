package collectionsrv_test

import (
	"context"
	"errors"
	"testing"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	mocks "github.com/10Narratives/distgo-db/internal/services/worker/data/collection/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateCollection(t *testing.T) {
	t.Parallel()

	const (
		parent       = "databases/db"
		collectionID = "coll1"
		name         = parent + "/collections/" + collectionID
		description  = "My test collection"
	)

	type fields struct {
		setupStorageMock func(m *mocks.CollectionStorage)
	}
	type args struct {
		ctx context.Context
		req struct {
			parent       string
			collectionID string
			description  string
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
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("CreateCollection", mock.Anything, key, description).
						Return(collectionmodels.Collection{
							Name:        name,
							Description: description,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent       string
					collectionID string
					description  string
				}{
					parent:       parent,
					collectionID: collectionID,
					description:  description,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll, ok := got.(collectionmodels.Collection)
				require.True(tt, ok)
				assert.Equal(tt, name, coll.Name)
				assert.Equal(tt, description, coll.Description)
			},
			wantErr: require.NoError,
		},
		{
			name: "storage returns error",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("CreateCollection", mock.Anything, key, description).
						Return(collectionmodels.Collection{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					parent       string
					collectionID string
					description  string
				}{
					parent:       parent,
					collectionID: collectionID,
					description:  description,
				},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				require.Error(tt, err)
				assert.Contains(tt, err.Error(), "internal error")
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewCollectionStorage(t)
			tt.fields.setupStorageMock(store)
			service := collectionsrv.New(store)

			res, err := service.CreateCollection(
				tt.args.ctx,
				tt.args.req.parent,
				tt.args.req.collectionID,
				tt.args.req.description,
			)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_GetCollection(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1"

	type fields struct {
		setupStorageMock func(m *mocks.CollectionStorage)
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
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("Collection", mock.Anything, key).
						Return(collectionmodels.Collection{
							Name: name,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct{ name string }{name: name},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll, ok := got.(collectionmodels.Collection)
				require.True(tt, ok)
				assert.Equal(tt, name, coll.Name)
			},
			wantErr: require.NoError,
		},
		{
			name: "not found",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("Collection", mock.Anything, key).
						Return(collectionmodels.Collection{}, errors.New("not found"))
				},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewCollectionStorage(t)
			tt.fields.setupStorageMock(store)
			service := collectionsrv.New(store)

			res, err := service.Collection(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, res)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_ListCollections(t *testing.T) {
	t.Parallel()

	const (
		parent = "databases/db"
		size   = int32(2)
		token  = "token"
	)

	type fields struct {
		setupStorageMock func(m *mocks.CollectionStorage)
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
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful list",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {
					parentKey := databasemodels.NewKey(parent)
					m.On("Collections", mock.Anything, parentKey).
						Return([]collectionmodels.Collection{
							{Name: parent + "/collections/coll1"},
							{Name: parent + "/collections/coll2"},
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
				list, ok := got.([]collectionmodels.Collection)
				require.True(tt, ok)
				require.Len(tt, list, 2)
				assert.Equal(tt, parent+"/collections/coll1", list[0].Name)
				assert.Equal(tt, parent+"/collections/coll2", list[1].Name)
			},
			wantErr: require.NoError,
		},
		{
			name: "with next page token",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {
					parentKey := databasemodels.NewKey(parent)
					m.On("Collections", mock.Anything, parentKey).
						Return([]collectionmodels.Collection{
							{Name: parent + "/collections/coll1"},
							{Name: parent + "/collections/coll2"},
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
				list, ok := got.([]collectionmodels.Collection)
				require.True(tt, ok)
				require.Len(tt, list, 2)
				assert.Equal(tt, parent+"/collections/coll1", list[0].Name)
				assert.Equal(tt, parent+"/collections/coll2", list[1].Name)
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewCollectionStorage(t)
			tt.fields.setupStorageMock(store)
			service := collectionsrv.New(store)

			list, _, err := service.Collections(tt.args.ctx, tt.args.req.parent, tt.args.req.size, tt.args.req.token)

			tt.wantVal(t, list)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_DeleteCollection(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1"

	type fields struct {
		setupStorageMock func(m *mocks.CollectionStorage)
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
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("DeleteCollection", mock.Anything, key).Return(nil)
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
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("DeleteCollection", mock.Anything, key).Return(errors.New("not found"))
				},
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewCollectionStorage(t)
			tt.fields.setupStorageMock(store)
			service := collectionsrv.New(store)

			err := service.DeleteCollection(tt.args.ctx, tt.args.req.name)

			tt.wantVal(t, nil)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}

func TestService_UpdateCollection(t *testing.T) {
	t.Parallel()

	const name = "databases/db/collections/coll1"
	const newDescription = "Updated description"

	type fields struct {
		setupStorageMock func(m *mocks.CollectionStorage)
	}
	type args struct {
		ctx context.Context
		req struct {
			collection collectionmodels.Collection
			paths      []string
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
			name: "update description",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(name)
					m.On("UpdateCollection", mock.Anything, key, newDescription).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					collection collectionmodels.Collection
					paths      []string
				}{
					collection: collectionmodels.Collection{
						Name:        name,
						Description: newDescription,
					},
					paths: []string{"description"},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				updated, ok := got.(collectionmodels.Collection)
				require.True(tt, ok)
				assert.Equal(tt, newDescription, updated.Description)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field",
			fields: fields{
				setupStorageMock: func(m *mocks.CollectionStorage) {},
			},
			args: args{
				ctx: context.Background(),
				req: struct {
					collection collectionmodels.Collection
					paths      []string
				}{
					collection: collectionmodels.Collection{
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			store := mocks.NewCollectionStorage(t)
			tt.fields.setupStorageMock(store)
			service := collectionsrv.New(store)

			updated, err := service.UpdateCollection(tt.args.ctx, tt.args.req.collection, tt.args.req.paths)

			tt.wantVal(t, updated)
			tt.wantErr(t, err)
			store.AssertExpectations(t)
		})
	}
}
