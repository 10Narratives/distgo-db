package walgrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal/mocks"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServerAPI_ListWALEntries(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC() // Use UTC to avoid timezone issues
	testTime := timestamppb.New(now)

	type fields struct {
		setupMock func(m *mocks.WALService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.ListWALEntriesRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal *dbv1.ListWALEntriesResponse
		wantErr error
	}{
		{
			name: "successful listing",
			fields: fields{
				setupMock: func(m *mocks.WALService) {
					m.On("WALEntries", context.Background(), int32(10), "token", now.Add(-time.Hour), now).
						Return([]walmodels.WALEntry{
							{
								ID:        "entry1",
								Target:    "table1",
								Type:      commonmodels.MutationTypeCreate,
								NewValue:  "new value",
								OldValue:  "",
								Timestamp: now,
							},
						}, "", nil) // Include empty token and nil error
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListWALEntriesRequest{
					PageSize:  10,
					PageToken: "token",
					StartTime: timestamppb.New(now.Add(-time.Hour)),
					EndTime:   testTime,
				},
			},
			wantVal: &dbv1.ListWALEntriesResponse{
				Entries: []*dbv1.WALEntry{
					{
						Id:            "entry1",
						Target:        "table1",
						OperationType: dbv1.MutationType_MUTATION_TYPE_CREATE,
						NewValue:      []byte("new value"),
						OldValue:      []byte(""),
						Timestamp:     testTime,
					},
				},
				NextPageToken: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid request",
			fields: fields{
				setupMock: func(m *mocks.WALService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListWALEntriesRequest{
					PageSize: -1, // invalid page size
				},
			},
			wantVal: nil,
			wantErr: status.Error(codes.InvalidArgument, "invalid ListWALEntriesRequest.PageSize: value must be inside range [0, 1000]"),
		},
		{
			name: "empty result",
			fields: fields{
				setupMock: func(m *mocks.WALService) {
					m.On("WALEntries", context.Background(), int32(10), "", now.Add(-time.Hour), now).
						Return([]walmodels.WALEntry{}, "", nil) // Include empty token and nil error
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.ListWALEntriesRequest{
					PageSize:  10,
					StartTime: timestamppb.New(now.Add(-time.Hour)),
					EndTime:   testTime,
				},
			},
			wantVal: &dbv1.ListWALEntriesResponse{
				Entries:       []*dbv1.WALEntry{},
				NextPageToken: "",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewWALService(t)
			tt.fields.setupMock(mock)

			serverAPI := walgrpc.New(mock)
			resp, err := serverAPI.ListWALEntries(tt.args.ctx, tt.args.req)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantVal, resp)
			}

			mock.AssertExpectations(t)
		})
	}
}

func TestServerAPI_TruncateWAL(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC() // Use UTC to avoid timezone issues
	testTime := timestamppb.New(now)

	type fields struct {
		setupMock func(m *mocks.WALService)
	}

	type args struct {
		ctx context.Context
		req *dbv1.TruncateWALRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal *emptypb.Empty
		wantErr error
	}{
		{
			name: "successful truncation",
			fields: fields{
				setupMock: func(m *mocks.WALService) {
					m.On("TruncateWAL", context.Background(), now).
						Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.TruncateWALRequest{
					Before: testTime,
				},
			},
			wantVal: &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "missing before timestamp",
			fields: fields{
				setupMock: func(m *mocks.WALService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.TruncateWALRequest{
					Before: nil,
				},
			},
			wantVal: nil,
			wantErr: status.Error(codes.InvalidArgument, "missing 'before' timestamp"),
		},
		{
			name: "zero before timestamp",
			fields: fields{
				setupMock: func(m *mocks.WALService) {},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.TruncateWALRequest{
					Before: timestamppb.New(time.Time{}),
				},
			},
			wantVal: nil,
			wantErr: status.Error(codes.InvalidArgument, "invalid 'before' timestamp"),
		},
		{
			name: "truncation error",
			fields: fields{
				setupMock: func(m *mocks.WALService) {
					m.On("TruncateWAL", context.Background(), now).
						Return(errors.New("truncation failed"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.TruncateWALRequest{
					Before: testTime,
				},
			},
			wantVal: nil,
			wantErr: status.Error(codes.Internal, "truncation failed"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewWALService(t)
			tt.fields.setupMock(mock)

			serverAPI := walgrpc.New(mock)
			resp, err := serverAPI.TruncateWAL(tt.args.ctx, tt.args.req)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Equal(t, tt.wantErr.Error(), err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantVal, resp)
			}

			mock.AssertExpectations(t)
		})
	}
}
