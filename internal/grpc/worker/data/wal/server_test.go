package walgrpc_test

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"testing"
// 	"time"

// 	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
// 	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal/mocks"
// 	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// 	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// func TestServerAPI_ListWALEntries(t *testing.T) {
// 	t.Parallel()

// 	now := time.Now()
// 	testUUID := uuid.New()
// 	testTimeProto := timestamppb.New(now)

// 	type fields struct {
// 		setupWalServiceMock func(m *mocks.WALService)
// 	}
// 	type args struct {
// 		ctx context.Context
// 		req *dbv1.ListWALEntriesRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "success with database entries",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					dbPayload := walmodels.DatabasePayload{
// 						Key: databasemodels.Key{Database: "db1"},
// 					}
// 					payloadData, _ := json.Marshal(dbPayload)

// 					m.On("Entries", mock.Anything, int32(10), "", now.UTC().Add(-time.Hour), now.UTC()).
// 						Return([]walmodels.WALEntry{
// 							{
// 								ID:        testUUID,
// 								Entity:    walmodels.EntityTypeDatabase,
// 								Mutation:  commonmodels.MutationTypeCreate,
// 								Timestamp: now,
// 								Payload:   payloadData,
// 							},
// 						}, "next-token", nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize:  10,
// 					StartTime: timestamppb.New(now.Add(-time.Hour)),
// 					EndTime:   testTimeProto,
// 				},
// 			},
// 			wantVal: func(t require.TestingT, val interface{}, msgAndArgs ...interface{}) {
// 				resp, ok := val.(*dbv1.ListWALEntriesResponse)
// 				require.True(t, ok)
// 				require.Len(t, resp.Entries, 1)
// 				require.Equal(t, testUUID.String(), resp.Entries[0].Id)
// 				require.Equal(t, dbv1.EntityType_ENTITY_TYPE_DATABASE, resp.Entries[0].EntityType)
// 				require.Equal(t, dbv1.MutationType_MUTATION_TYPE_CREATE, resp.Entries[0].OperationType)
// 				require.Equal(t, "next-token", resp.NextPageToken)
// 			},
// 			wantErr: require.NoError,
// 		},

// 		{
// 			name: "success with collection entries",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					collPayload := walmodels.CollectionPayload{
// 						Key: collectionmodels.Key{Database: "db1", Collection: "col1"},
// 					}
// 					payloadData, _ := json.Marshal(collPayload)

// 					m.On("Entries",
// 						mock.Anything,
// 						int32(20),
// 						"token",
// 						mock.AnythingOfType("time.Time"),
// 						mock.AnythingOfType("time.Time"),
// 					).Return([]walmodels.WALEntry{
// 						{
// 							ID:        testUUID,
// 							Entity:    walmodels.EntityTypeCollection,
// 							Mutation:  commonmodels.MutationTypeUpdate,
// 							Timestamp: now,
// 							Payload:   payloadData,
// 						},
// 					}, "", nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize:  20,
// 					PageToken: "token",
// 				},
// 			},
// 			wantVal: func(t require.TestingT, val interface{}, msgAndArgs ...interface{}) {
// 				resp, ok := val.(*dbv1.ListWALEntriesResponse)
// 				require.True(t, ok)
// 				require.Len(t, resp.Entries, 1)
// 				require.Equal(t, "col1", resp.Entries[0].GetCollectionPayload().CollectionId)
// 			},
// 			wantErr: require.NoError,
// 		},

// 		{
// 			name: "success with document entries",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					docPayload := walmodels.DocumentPayload{
// 						Key: documentmodels.Key{Database: "db1", Collection: "col1", Document: "doc1"},
// 					}
// 					payloadData, _ := json.Marshal(docPayload)
// 					m.On("Entries",
// 						mock.Anything,
// 						int32(10),
// 						"",
// 						mock.AnythingOfType("time.Time"),
// 						mock.AnythingOfType("time.Time"),
// 					).Return([]walmodels.WALEntry{
// 						{
// 							ID:        testUUID,
// 							Entity:    walmodels.EntityTypeDocument,
// 							Mutation:  commonmodels.MutationTypeDelete,
// 							Timestamp: now,
// 							Payload:   payloadData,
// 						},
// 					}, "", nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize: 10,
// 				},
// 			},
// 			wantVal: func(t require.TestingT, val interface{}, msgAndArgs ...interface{}) {
// 				resp, ok := val.(*dbv1.ListWALEntriesResponse)
// 				require.True(t, ok)
// 				require.Len(t, resp.Entries, 1)
// 				require.Equal(t, "doc1", resp.Entries[0].GetDocumentPayload().DocumentId)
// 			},
// 			wantErr: require.NoError,
// 		},

// 		{
// 			name: "validation error",
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize: -1, // Invalid value
// 				},
// 			},
// 			wantVal: require.Nil,
// 			wantErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.InvalidArgument, st.Code())
// 			},
// 		},

// 		{
// 			name: "service error",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					m.On("Entries",
// 						mock.Anything,
// 						int32(10),
// 						"",
// 						mock.AnythingOfType("time.Time"),
// 						mock.AnythingOfType("time.Time"),
// 					).Return(nil, "", errors.New("internal error"))
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize: 10,
// 				},
// 			},
// 			wantVal: require.Nil,
// 			wantErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.Internal, st.Code())
// 			},
// 		},

// 		{
// 			name: "payload unmarshal error",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					m.On("Entries",
// 						mock.Anything,
// 						int32(10),
// 						"",
// 						mock.AnythingOfType("time.Time"),
// 						mock.AnythingOfType("time.Time"),
// 					).Return([]walmodels.WALEntry{
// 						{
// 							ID:        testUUID,
// 							Entity:    walmodels.EntityTypeDatabase,
// 							Mutation:  commonmodels.MutationTypeCreate,
// 							Timestamp: now,
// 							Payload:   json.RawMessage("{invalid json}"),
// 						},
// 					}, "", nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.ListWALEntriesRequest{
// 					PageSize: 10,
// 				},
// 			},
// 			wantVal: require.Nil,
// 			wantErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.Internal, st.Code())
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			walServiceMock := mocks.NewWALService(t)
// 			if tt.fields.setupWalServiceMock != nil {
// 				tt.fields.setupWalServiceMock(walServiceMock)
// 			}

// 			serverAPI := walgrpc.New(walServiceMock)
// 			resp, err := serverAPI.ListWALEntries(tt.args.ctx, tt.args.req)

// 			tt.wantVal(t, resp)
// 			tt.wantErr(t, err)

// 			walServiceMock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestServerAPI_TruncateWAL(t *testing.T) {
// 	t.Parallel()

// 	now := time.Now()
// 	testTimeProto := timestamppb.New(now)

// 	type fields struct {
// 		setupWalServiceMock func(m *mocks.WALService)
// 	}
// 	type args struct {
// 		ctx context.Context
// 		req *dbv1.TruncateWALRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantVal require.ValueAssertionFunc
// 		wantErr require.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "successful truncation",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					m.On("Truncate", mock.Anything, mock.AnythingOfType("time.Time")).
// 						Return(nil)
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.TruncateWALRequest{
// 					Before: testTimeProto,
// 				},
// 			},
// 			wantVal: require.Empty,
// 			wantErr: require.NoError,
// 		},
// 		{
// 			name: "service error",
// 			fields: fields{
// 				setupWalServiceMock: func(m *mocks.WALService) {
// 					m.On("Truncate", mock.Anything, mock.AnythingOfType("time.Time")).
// 						Return(errors.New("truncation failed"))
// 				},
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				req: &dbv1.TruncateWALRequest{
// 					Before: testTimeProto,
// 				},
// 			},
// 			wantVal: require.Empty,
// 			wantErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
// 				require.Error(t, err)
// 				st, ok := status.FromError(err)
// 				require.True(t, ok)
// 				require.Equal(t, codes.Internal, st.Code())
// 				require.Contains(t, st.Message(), "failed to truncate WAL")
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			walServiceMock := mocks.NewWALService(t)
// 			if tt.fields.setupWalServiceMock != nil {
// 				tt.fields.setupWalServiceMock(walServiceMock)
// 			}

// 			serverAPI := walgrpc.New(walServiceMock)
// 			got, err := serverAPI.TruncateWAL(tt.args.ctx, tt.args.req)

// 			tt.wantVal(t, got)
// 			tt.wantErr(t, err)
// 			walServiceMock.AssertExpectations(t)
// 		})
// 	}
// }
