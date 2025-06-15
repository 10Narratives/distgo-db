package databasegrpc_test

import (
	"context"
	"errors"
	"testing"

	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/database/mocks"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestServerAPI_CreateDatabase(t *testing.T) {
	t.Parallel()

	var (
		databaseID  string = "database_001"
		displayName string = "MEPHI"
	)

	type fields struct {
		setupDatabaseServiceMock func(m *mocks.DatabaseService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.CreateDatabaseRequest
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
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("CreateDatabase", mock.Anything, databaseID, displayName).
						Return(databasemodels.Database{
							Name:        "databases/" + databaseID,
							DisplayName: displayName,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDatabaseRequest{
					DatabaseId: databaseID,
					Database: &dbv1.Database{
						Name:        "databases/" + databaseID,
						DisplayName: displayName,
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				database, ok := got.(*dbv1.Database)
				require.True(t, ok)

				assert.Equal(t, "databases/"+databaseID, database.Name)
				assert.Equal(t, displayName, database.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDatabaseRequest{
					DatabaseId: "databaseID!@#$%^",
					Database: &dbv1.Database{
						Name:        "databases/" + databaseID,
						DisplayName: displayName,
					},
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("CreateDatabase", mock.Anything, databaseID, displayName).
						Return(databasemodels.Database{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CreateDatabaseRequest{
					DatabaseId: databaseID,
					Database: &dbv1.Database{
						Name:        "databases/" + databaseID,
						DisplayName: displayName,
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

			databaseSrv := mocks.NewDatabaseService(t)
			tt.fields.setupDatabaseServiceMock(databaseSrv)

			serverAPI := databasegrpc.New(databaseSrv)
			resp, err := serverAPI.CreateDatabase(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			databaseSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_DeleteDatabase(t *testing.T) {
	t.Parallel()

	var (
		databaseID string = "database_001"
	)

	type fields struct {
		setupDatabaseServiceMock func(m *mocks.DatabaseService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.DeleteDatabaseRequest
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
			fields: fields{setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
				m.On("DeleteDatabase", mock.Anything, "databases/"+databaseID).
					Return(nil)
			}},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDatabaseRequest{
					Name: "databases/" + databaseID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.NoError,
		},
		{
			name:   "validation error",
			fields: fields{setupDatabaseServiceMock: func(m *mocks.DatabaseService) {}},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDatabaseRequest{
					Name: "databases_wrong/" + databaseID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
				m.On("DeleteDatabase", mock.Anything, "databases/"+databaseID).
					Return(errors.New("internal"))
			}},
			args: args{
				ctx: context.Background(),
				req: &dbv1.DeleteDatabaseRequest{
					Name: "databases/" + databaseID,
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

			databaseSrv := mocks.NewDatabaseService(t)
			tt.fields.setupDatabaseServiceMock(databaseSrv)

			serverAPI := databasegrpc.New(databaseSrv)
			resp, err := serverAPI.DeleteDatabase(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			databaseSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_GetDatabase(t *testing.T) {
	t.Parallel()

	var (
		databaseID  string = "database_001"
		displayName string = "MEPHI"
	)

	type fields struct {
		setupDatabaseServiceMock func(m *mocks.DatabaseService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.GetDatabaseRequest
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
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("Database", mock.Anything, "databases/"+databaseID).
						Return(databasemodels.Database{
							Name:        "databases/" + databaseID,
							DisplayName: displayName,
						}, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDatabaseRequest{
					Name: "databases/" + databaseID,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				database, ok := got.(*dbv1.Database)
				require.True(t, ok)

				assert.Equal(t, "databases/"+databaseID, database.Name)
				assert.Equal(t, displayName, database.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDatabaseRequest{
					Name: "databases_wrong/" + databaseID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("Database", mock.Anything, "databases/"+databaseID).
						Return(databasemodels.Database{}, errors.New("internal"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.GetDatabaseRequest{
					Name: "databases/" + databaseID,
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

			databaseSrv := mocks.NewDatabaseService(t)
			tt.fields.setupDatabaseServiceMock(databaseSrv)

			serverAPI := databasegrpc.New(databaseSrv)
			resp, err := serverAPI.GetDatabase(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			databaseSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_ListDatabases(t *testing.T) {
	t.Parallel()

	var (
		databaseID1  string                  = "database_001"
		databaseID2  string                  = "database_002"
		displayName1 string                  = "MEPHI#1"
		displayName2 string                  = "MEPHI#2"
		database1    databasemodels.Database = databasemodels.Database{
			Name:        "databases/" + databaseID1,
			DisplayName: displayName1,
		}
		database2 databasemodels.Database = databasemodels.Database{
			Name:        "databases/" + databaseID2,
			DisplayName: displayName2,
		}
		size  int32  = 2
		token string = "token"
	)

	type fields struct {
		setupDatabaseServiceMock func(m *mocks.DatabaseService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.ListDatabasesRequest
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
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("Databases", mock.Anything, size, token).
						Return([]databasemodels.Database{
							database1,
							database2,
						}, "", nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDatabasesRequest{
					PageSize:  size,
					PageToken: token,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				resp, ok := got.(*dbv1.ListDatabasesResponse)
				require.True(tt, ok)

				require.Len(tt, resp.Databases, 2)

				require.Equal(tt, "databases/database_001", resp.Databases[0].Name)
				require.Equal(tt, "MEPHI#1", resp.Databases[0].DisplayName)

				require.Equal(tt, "databases/database_002", resp.Databases[1].Name)
				require.Equal(tt, "MEPHI#2", resp.Databases[1].DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDatabasesRequest{
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
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("Databases", mock.Anything, size, token).
						Return([]databasemodels.Database{}, "", errors.New("internal"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListDatabasesRequest{
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

			databaseSrv := mocks.NewDatabaseService(t)
			tt.fields.setupDatabaseServiceMock(databaseSrv)

			serverAPI := databasegrpc.New(databaseSrv)
			resp, err := serverAPI.ListDatabases(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			databaseSrv.AssertExpectations(t)
		})
	}
}

func TestServerAPI_UpdateDatabase(t *testing.T) {
	t.Parallel()

	var (
		databaseID  string                  = "database_001"
		displayName string                  = "MEPHI"
		paths       []string                = make([]string, 0)
		database    databasemodels.Database = databasemodels.Database{
			Name:        "databases/" + databaseID,
			DisplayName: displayName,
		}
	)

	type fields struct {
		setupDatabaseServiceMock func(m *mocks.DatabaseService)
	}
	type args struct {
		ctx context.Context
		req *dbv1.UpdateDatabaseRequest
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
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("UpdateDatabase", mock.Anything, database, paths).
						Return(database, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDatabaseRequest{
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{},
					},
					Database: &dbv1.Database{
						Name:        "databases/" + databaseID,
						DisplayName: displayName,
					},
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				database, ok := got.(*dbv1.Database)
				require.True(t, ok)

				assert.Equal(t, "databases/"+databaseID, database.Name)
				assert.Equal(t, displayName, database.DisplayName)
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDatabaseRequest{
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{},
					},
				},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupDatabaseServiceMock: func(m *mocks.DatabaseService) {
					m.On("UpdateDatabase", mock.Anything, database, paths).
						Return(databasemodels.Database{}, errors.New("internal"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.UpdateDatabaseRequest{
					UpdateMask: &fieldmaskpb.FieldMask{
						Paths: []string{},
					},
					Database: &dbv1.Database{
						Name:        "databases/" + databaseID,
						DisplayName: displayName,
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

			databaseSrv := mocks.NewDatabaseService(t)
			tt.fields.setupDatabaseServiceMock(databaseSrv)

			serverAPI := databasegrpc.New(databaseSrv)
			resp, err := serverAPI.UpdateDatabase(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			databaseSrv.AssertExpectations(t)
		})
	}
}
