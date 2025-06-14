package collectiongrpc_test

import (
	"context"
	"errors"
	"testing"

	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection/mocks"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestServerAPI_CreateCollection(t *testing.T) {
	t.Parallel()
	var (
		parent       = "databases/database_001"
		collectionID = "coll1"
	)
	type fields struct {
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.CreateCollectionRequest
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
					m.On("CreateCollection", mock.Anything, parent, collectionID).
						Return(collectionmodels.Collection{
							Name: parent + "/collections/" + collectionID,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateCollectionRequest{
					Parent:       parent,
					CollectionId: collectionID,
					Collection: &dbv1.Collection{
						Name: parent + "/collections/" + collectionID,
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll, ok := got.(*dbv1.Collection)
				require.True(tt, ok)
				assert.Equal(tt, parent+"/collections/"+collectionID, coll.GetName())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateCollectionRequest{
					Parent:       "",
					CollectionId: collectionID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("CreateCollection", mock.Anything, parent, collectionID).
						Return(collectionmodels.Collection{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateCollectionRequest{
					Parent:       parent,
					CollectionId: collectionID,
					Collection: &dbv1.Collection{
						Name: parent + "/collections/" + collectionID,
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
			svc := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(svc)
			server := collectiongrpc.New(svc)
			resp, err := server.CreateCollection(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_DeleteCollection(t *testing.T) {
	t.Parallel()
	var (
		name = "databases/ви/collections/coll1"
	)
	type fields struct {
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.DeleteCollectionRequest
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
					m.On("DeleteCollection", mock.Anything, name).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteCollectionRequest{
					Name: name,
				},
			},
			wantVal: require.Empty,
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteCollectionRequest{
					Name: "invalid_name",
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("DeleteCollection", mock.Anything, name).Return(errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteCollectionRequest{
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
			svc := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(svc)
			server := collectiongrpc.New(svc)
			resp, err := server.DeleteCollection(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_GetCollection(t *testing.T) {
	t.Parallel()
	var (
		name = "databases/db/collections/coll1"
	)
	type fields struct {
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.GetCollectionRequest
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
					m.On("Collection", mock.Anything, name).Return(collectionmodels.Collection{
						Name: name,
					}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetCollectionRequest{
					Name: name,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll, ok := got.(*dbv1.Collection)
				require.True(tt, ok)
				assert.Equal(tt, name, coll.GetName())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetCollectionRequest{
					Name: "invalid_name",
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collection", mock.Anything, name).Return(collectionmodels.Collection{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetCollectionRequest{
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
			svc := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(svc)
			server := collectiongrpc.New(svc)
			resp, err := server.GetCollection(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_ListCollections(t *testing.T) {
	t.Parallel()
	var (
		parent       = "databases/db"
		size   int32 = 2
		token        = "token"
	)
	type fields struct {
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.ListCollectionsRequest
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
					m.On("Collections", mock.Anything, parent, size, token).
						Return([]collectionmodels.Collection{
							{Name: parent + "/collections/coll1"},
							{Name: parent + "/collections/coll2"},
						}, "", nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListCollectionsRequest{
					Parent:    parent,
					PageSize:  size,
					PageToken: token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				resp, ok := got.(*dbv1.ListCollectionsResponse)
				require.True(tt, ok)
				require.Len(tt, resp.Collections, 2)
				assert.Equal(tt, parent+"/collections/coll1", resp.Collections[0].GetName())
				assert.Equal(tt, parent+"/collections/coll2", resp.Collections[1].GetName())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListCollectionsRequest{
					Parent:    "",
					PageSize:  10000,
					PageToken: token,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("Collections", mock.Anything, parent, size, token).
						Return([]collectionmodels.Collection{}, "", errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListCollectionsRequest{
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
			svc := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(svc)
			server := collectiongrpc.New(svc)
			resp, err := server.ListCollections(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}

func TestServerAPI_UpdateCollection(t *testing.T) {
	t.Parallel()
	var (
		name = "databases/db/collections/coll1"
	)
	type fields struct {
		setupCollectionServiceMock func(m *mocks.CollectionService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.UpdateCollectionRequest
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
					m.On("UpdateCollection", mock.Anything, collectionmodels.Collection{
						Name: name,
					}, []string{"display_name"}).
						Return(collectionmodels.Collection{
							Name: name,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateCollectionRequest{
					Collection: &dbv1.Collection{
						Name: name,
					},
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{"display_name"},
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i ...interface{}) {
				coll, ok := got.(*dbv1.Collection)
				require.True(tt, ok)
				assert.Equal(tt, name, coll.GetName())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateCollectionRequest{},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupCollectionServiceMock: func(m *mocks.CollectionService) {
					m.On("UpdateCollection", mock.Anything, collectionmodels.Collection{
						Name: name,
					}, []string{"display_name"}).
						Return(collectionmodels.Collection{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateCollectionRequest{
					Collection: &dbv1.Collection{
						Name: name,
					},
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{"display_name"},
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
			svc := mocks.NewCollectionService(t)
			tt.fields.setupCollectionServiceMock(svc)
			server := collectiongrpc.New(svc)
			resp, err := server.UpdateCollection(tt.args.ctx, tt.args.req)
			tt.wantVal(t, resp)
			tt.wantErr(t, err)
			svc.AssertExpectations(t)
		})
	}
}
