package collectionsrv_test

import (
	"context"
	"errors"
	"testing"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	"github.com/10Narratives/distgo-db/internal/services/worker/data/collection/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	name        string    = "databases/test_db/collections/users"
	description string    = "The awesome test collection"
	createdAt   time.Time = time.Now().UTC()
	updatedAt   time.Time = time.Now().UTC()
	key                   = collectionmodels.NewKey(name)
	collection            = collectionmodels.Collection{
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	errInternalCollectionStorage        = errors.New("internal collection storage error")
	parent                       string = "databases/test_db"
	parentKey                           = databasemodels.NewKey(parent)
	collectionID                 string = "users"
	errInternalWALService               = errors.New("internal wal service error")
)

type fields struct {
	setupCollectionStorageMock func(m *mocks.CollectionStorage)
	setupWALServiceMock        func(m *mocks.WAlService)
}

func TestService_Collection(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		name string
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
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).Return(collection, nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				col, ok := got.(collectionmodels.Collection)
				require.True(t, ok)

				assert.Equal(t, collection, col)
			},
			wantErr: require.NoError,
		},
		{
			name: "collection storage error",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).
						Return(collectionmodels.Collection{}, errInternalCollectionStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collectionStorageMock := mocks.NewCollectionStorage(t)
			tt.fields.setupCollectionStorageMock(collectionStorageMock)

			walServiceMocks := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMocks)

			service := collectionsrv.New(collectionStorageMock, walServiceMocks)

			collection, err := service.Collection(tt.args.ctx, tt.args.name)

			tt.wantVal(t, collection)
			tt.wantErr(t, err)

			collectionStorageMock.AssertExpectations(t)
			walServiceMocks.AssertExpectations(t)
		})
	}
}

func TestService_Collections(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		parent string
		size   int32
		token  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful execution with pagination",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					allCollections := []collectionmodels.Collection{
						{
							Name:        "databases/test_db/collections/users",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
						{
							Name:        "databases/test_db/collections/orders",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
						{
							Name:        "databases/test_db/collections/products",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
					}
					m.On("Collections", mock.Anything, parentKey).Return(allCollections, nil)
				},
			},
			args: args{
				ctx:    context.Background(),
				parent: parent,
				size:   2,
				token:  "",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]collectionmodels.Collection)
				require.True(t, ok)

				expected := []collectionmodels.Collection{
					{
						Name:        "databases/test_db/collections/users",
						Description: description,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
					{
						Name:        "databases/test_db/collections/orders",
						Description: description,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
				}
				assert.Equal(t, expected, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "empty collections list",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collections", mock.Anything, parentKey).Return([]collectionmodels.Collection{}, nil)
				},
			},
			args: args{
				ctx:    context.Background(),
				parent: parent,
				size:   10,
				token:  "",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]collectionmodels.Collection)
				require.True(t, ok)
				assert.Empty(t, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "pagination with token",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					allCollections := []collectionmodels.Collection{
						{
							Name:        "databases/test_db/collections/users",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
						{
							Name:        "databases/test_db/collections/orders",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
						{
							Name:        "databases/test_db/collections/products",
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						},
					}
					m.On("Collections", mock.Anything, parentKey).Return(allCollections, nil)
				},
			},
			args: args{
				ctx:    context.Background(),
				parent: parent,
				size:   1,
				token:  "databases/test_db/collections/users",
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				result, ok := got.([]collectionmodels.Collection)
				require.True(t, ok)

				expected := []collectionmodels.Collection{
					{
						Name:        "databases/test_db/collections/orders",
						Description: description,
						CreatedAt:   createdAt,
						UpdatedAt:   updatedAt,
					},
				}
				assert.Equal(t, expected, result)
			},
			wantErr: require.NoError,
		},
		{
			name: "storage error",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collections", mock.Anything, parentKey).Return(nil, errInternalCollectionStorage)
				},
			},
			args: args{
				ctx:    context.Background(),
				parent: parent,
				size:   10,
				token:  "",
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			collectionStorageMock := mocks.NewCollectionStorage(t)
			tt.fields.setupCollectionStorageMock(collectionStorageMock)

			service := collectionsrv.New(collectionStorageMock, nil)

			page, nextPageToken, err := service.Collections(tt.args.ctx, tt.args.parent, tt.args.size, tt.args.token)

			tt.wantVal(t, page)
			tt.wantErr(t, err)

			if err == nil && len(page) > 0 && nextPageToken != "" {
				assert.Equal(t, tt.args.size, int32(len(page)))
			}

			collectionStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateCollection(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx          context.Context
		parent       string
		collectionID string
		description  string
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
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
						Return(nil)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(parent + "/collections/" + collectionID)
					m.On("CreateCollection", mock.Anything, key, description).
						Return(collectionmodels.Collection{
							Name:        parent + "/collections/" + collectionID,
							Description: description,
							CreatedAt:   createdAt,
							UpdatedAt:   updatedAt,
						}, nil)
				},
			},
			args: args{
				ctx:          context.Background(),
				parent:       parent,
				collectionID: collectionID,
				description:  description,
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				coll, ok := got.(collectionmodels.Collection)
				require.True(t, ok)

				assert.Equal(t, parent+"/collections/"+collectionID, coll.Name)
				assert.Equal(t, description, coll.Description)
				assert.WithinDuration(t, time.Now(), coll.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), coll.UpdatedAt, time.Second)
			},
			wantErr: require.NoError,
		},
		{
			name: "WAL service error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
						Return(errInternalWALService)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {},
			},
			args: args{
				ctx:          context.Background(),
				parent:       parent,
				collectionID: collectionID,
				description:  description,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.Contains(t, err.Error(), errInternalWALService.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeCreate, mock.Anything, mock.Anything).
						Return(nil)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					key := collectionmodels.NewKey(parent + "/collections/" + collectionID)
					m.On("CreateCollection", mock.Anything, key, description).
						Return(collectionmodels.Collection{}, errInternalCollectionStorage)
				},
			},
			args: args{
				ctx:          context.Background(),
				parent:       parent,
				collectionID: collectionID,
				description:  description,
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			collectionStorageMock := mocks.NewCollectionStorage(t)
			tt.fields.setupCollectionStorageMock(collectionStorageMock)

			service := collectionsrv.New(collectionStorageMock, walServiceMock)

			coll, err := service.CreateCollection(tt.args.ctx, tt.args.parent, tt.args.collectionID, tt.args.description)

			tt.wantVal(t, coll)
			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			collectionStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_DeleteCollection(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		name string
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
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*collectionmodels.Collection)(nil)).
						Return(nil)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("DeleteCollection", mock.Anything, key).
						Return(nil)
				},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: require.NoError,
		},
		{
			name: "WAL service error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*collectionmodels.Collection)(nil)).
						Return(errInternalWALService)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "failed to create WAL entry: "+errInternalWALService.Error())
			},
		},
		{
			name: "storage error",
			fields: fields{
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeDelete, key, (*collectionmodels.Collection)(nil)).
						Return(nil)
				},
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("DeleteCollection", mock.Anything, key).
						Return(errInternalCollectionStorage)
				},
			},
			args: args{
				ctx:  context.Background(),
				name: name,
			},
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			collectionStorageMock := mocks.NewCollectionStorage(t)
			tt.fields.setupCollectionStorageMock(collectionStorageMock)

			service := collectionsrv.New(collectionStorageMock, walServiceMock)

			err := service.DeleteCollection(tt.args.ctx, tt.args.name)

			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			collectionStorageMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateCollection(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		collection collectionmodels.Collection
		paths      []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "successful update of description",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).
						Return(collection, nil)
					m.On("UpdateCollection", mock.Anything, key, "New Description").
						Return(nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				collection: collectionmodels.Collection{
					Name:        name,
					Description: "New Description",
					CreatedAt:   createdAt,
					UpdatedAt:   time.Now().UTC(),
				},
				paths: []string{"description"},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				updatedColl, ok := got.(collectionmodels.Collection)
				require.True(t, ok)

				assert.Equal(t, "New Description", updatedColl.Description)
				assert.WithinDuration(t, time.Now(), updatedColl.UpdatedAt, time.Second)
			},
			wantErr: require.NoError,
		},
		{
			name: "unknown field",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {},
				setupWALServiceMock:        func(m *mocks.WAlService) {},
			},
			args: args{
				ctx: context.Background(),
				collection: collectionmodels.Collection{
					Name:        name,
					Description: "New Description",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				paths: []string{"unknown_field"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "unknown field: unknown_field")
			},
		},
		{
			name: "storage error when fetching collection",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).
						Return(collectionmodels.Collection{}, errInternalCollectionStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {},
			},
			args: args{
				ctx: context.Background(),
				collection: collectionmodels.Collection{
					Name:        name,
					Description: "New Description",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				paths: []string{"description"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
		{
			name: "WAL service error",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).
						Return(collection, nil)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(errInternalWALService)
				},
			},
			args: args{
				ctx: context.Background(),
				collection: collectionmodels.Collection{
					Name:        name,
					Description: "New Description",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				paths: []string{"description"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, "failed to create WAL entry: "+errInternalWALService.Error())
			},
		},
		{
			name: "storage error when updating collection",
			fields: fields{
				setupCollectionStorageMock: func(m *mocks.CollectionStorage) {
					m.On("Collection", mock.Anything, key).
						Return(collection, nil)
					m.On("UpdateCollection", mock.Anything, key, "New Description").
						Return(errInternalCollectionStorage)
				},
				setupWALServiceMock: func(m *mocks.WAlService) {
					m.On("CreateCollectionEntry", mock.Anything, commonmodels.MutationTypeUpdate, key, mock.Anything).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				collection: collectionmodels.Collection{
					Name:        name,
					Description: "New Description",
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				paths: []string{"description"},
			},
			wantVal: require.Empty,
			wantErr: func(tt require.TestingT, err error, i ...interface{}) {
				assert.EqualError(t, err, errInternalCollectionStorage.Error())
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walServiceMock := mocks.NewWAlService(t)
			tt.fields.setupWALServiceMock(walServiceMock)

			collectionStorageMock := mocks.NewCollectionStorage(t)
			tt.fields.setupCollectionStorageMock(collectionStorageMock)

			service := collectionsrv.New(collectionStorageMock, walServiceMock)

			updatedColl, err := service.UpdateCollection(tt.args.ctx, tt.args.collection, tt.args.paths)

			tt.wantVal(t, updatedColl)
			tt.wantErr(t, err)

			walServiceMock.AssertExpectations(t)
			collectionStorageMock.AssertExpectations(t)
		})
	}
}
