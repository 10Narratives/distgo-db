package transactiongrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	transactiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/transaction"
	"github.com/10Narratives/distgo-db/internal/grpc/worker/data/transaction/mocks"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	startedAt     = time.Now().UTC()
	transactionID = "txn-1234"
	description   = "description"
)

type fields struct {
	setupTxServiceMock func(m *mocks.TransactionService)
}

func TestServerAPI_Begin(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *dbv1.BeginRequest
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
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Begin", mock.Anything, description).
						Return(transactionID, startedAt, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.BeginRequest{
					Description: description,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				resp, ok := got.(*dbv1.BeginResponse)
				require.True(t, ok)

				assert.Equal(t, transactionID, resp.TransactionId)
				assert.Equal(t, startedAt, resp.StartedAt.AsTime())
			},
			wantErr: require.NoError,
		},
		{
			name: "internal error",
			fields: fields{
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Begin", mock.Anything, description).
						Return("transactionID", time.Time{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.BeginRequest{
					Description: description,
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

			txServiceMock := mocks.NewTransactionService(t)
			tt.fields.setupTxServiceMock(txServiceMock)

			serverAPI := transactiongrpc.New(txServiceMock)
			resp, err := serverAPI.Begin(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			txServiceMock.AssertExpectations(t)
		})
	}
}

func TestServerAPI_Commit(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *dbv1.CommitRequest
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
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Commit", mock.Anything, transactionID).
						Return(startedAt, nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CommitRequest{
					TransactionId: transactionID,
				},
			},
			wantVal: func(tt require.TestingT, got interface{}, i2 ...interface{}) {
				resp, ok := got.(*dbv1.CommitResponse)
				require.True(t, ok)

				assert.Equal(t, startedAt, resp.CommittedAt.AsTime())
			},
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupTxServiceMock: func(m *mocks.TransactionService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CommitRequest{},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Commit", mock.Anything, transactionID).
						Return(time.Time{}, errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.CommitRequest{
					TransactionId: transactionID,
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

			txServiceMock := mocks.NewTransactionService(t)
			tt.fields.setupTxServiceMock(txServiceMock)

			serverAPI := transactiongrpc.New(txServiceMock)
			resp, err := serverAPI.Commit(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			txServiceMock.AssertExpectations(t)
		})
	}
}

func TestServerAPI_Rollback(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *dbv1.RollbackRequest
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
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Rollback", mock.Anything, transactionID).Return(nil)
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.RollbackRequest{
					TransactionId: transactionID,
				},
			},
			wantVal: require.Empty,
			wantErr: require.NoError,
		},
		{
			name: "validation error",
			fields: fields{
				setupTxServiceMock: func(m *mocks.TransactionService) {
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.RollbackRequest{},
			},
			wantVal: require.Empty,
			wantErr: require.Error,
		},
		{
			name: "internal error",
			fields: fields{
				setupTxServiceMock: func(m *mocks.TransactionService) {
					m.On("Rollback", mock.Anything, transactionID).Return(errors.New("internal error"))
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dbv1.RollbackRequest{
					TransactionId: transactionID,
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

			txServiceMock := mocks.NewTransactionService(t)
			tt.fields.setupTxServiceMock(txServiceMock)

			serverAPI := transactiongrpc.New(txServiceMock)
			resp, err := serverAPI.Rollback(tt.args.ctx, tt.args.req)

			tt.wantVal(t, resp)
			tt.wantErr(t, err)

			txServiceMock.AssertExpectations(t)
		})
	}
}
