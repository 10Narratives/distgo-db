package walsrv_test

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
// 	"github.com/10Narratives/distgo-db/internal/services/worker/data/wal/mocks"

// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// 	"github.com/stretchr/testify/require"
// )

// func TestService_WALEntries(t *testing.T) {
// 	t.Parallel()

// 	now := time.Now().UTC() // Use UTC to avoid timezone issues

// 	type fields struct {
// 		setupMock func(m *mocks.WALStorage)
// 	}

// 	type args struct {
// 		ctx   context.Context
// 		size  int32
// 		token string
// 		from  time.Time
// 		to    time.Time
// 	}

// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantEntries []walmodels.WALEntry
// 		wantToken   string
// 		wantErr     error
// 	}{
// 		{
// 			name: "successful listing",
// 			fields: fields{
// 				setupMock: func(m *mocks.WALStorage) {
// 					m.On("Entries", context.Background(), int32(10), "token", now.Add(-time.Hour), now).
// 						Return([]walmodels.WALEntry{
// 							{
// 								ID:        "entry1",
// 								Target:    "table1",
// 								Type:      commonmodels.MutationTypeCreate,
// 								NewValue:  "new value",
// 								OldValue:  "",
// 								Timestamp: now,
// 							},
// 						}, "next_token", nil)
// 				},
// 			},
// 			args: args{
// 				ctx:   context.Background(),
// 				size:  10,
// 				token: "token",
// 				from:  now.Add(-time.Hour),
// 				to:    now,
// 			},
// 			wantEntries: []walmodels.WALEntry{
// 				{
// 					ID:        "entry1",
// 					Target:    "table1",
// 					Type:      commonmodels.MutationTypeCreate,
// 					NewValue:  "new value",
// 					OldValue:  "",
// 					Timestamp: now,
// 				},
// 			},
// 			wantToken: "next_token",
// 			wantErr:   nil,
// 		},
// 		{
// 			name: "empty result",
// 			fields: fields{
// 				setupMock: func(m *mocks.WALStorage) {
// 					m.On("Entries", context.Background(), int32(10), "", now.Add(-time.Hour), now).
// 						Return([]walmodels.WALEntry{}, "", nil)
// 				},
// 			},
// 			args: args{
// 				ctx:   context.Background(),
// 				size:  10,
// 				token: "",
// 				from:  now.Add(-time.Hour),
// 				to:    now,
// 			},
// 			wantEntries: []walmodels.WALEntry{},
// 			wantToken:   "",
// 			wantErr:     nil,
// 		},
// 		{
// 			name: "storage error",
// 			fields: fields{
// 				setupMock: func(m *mocks.WALStorage) {
// 					m.On("Entries", context.Background(), int32(10), "token", now.Add(-time.Hour), now).
// 						Return(nil, "", errors.New("storage error"))
// 				},
// 			},
// 			args: args{
// 				ctx:   context.Background(),
// 				size:  10,
// 				token: "token",
// 				from:  now.Add(-time.Hour),
// 				to:    now,
// 			},
// 			wantEntries: nil,
// 			wantToken:   "",
// 			wantErr:     errors.New("storage error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			mock := mocks.NewWALStorage(t)
// 			tt.fields.setupMock(mock)

// 			service := walsrv.New(mock)
// 			respEntries, respToken, err := service.WALEntries(tt.args.ctx, tt.args.size, tt.args.token, tt.args.from, tt.args.to)

// 			if tt.wantErr != nil {
// 				require.Error(t, err)
// 				require.Equal(t, tt.wantErr.Error(), err.Error())
// 				require.Nil(t, respEntries)
// 				require.Empty(t, respToken)
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, tt.wantEntries, respEntries)
// 				require.Equal(t, tt.wantToken, respToken)
// 			}

// 			mock.AssertExpectations(t)
// 		})
// 	}
// }

// func TestService_TruncateWAL(t *testing.T) {
// 	t.Parallel()

// 	now := time.Now().UTC() // Use UTC to avoid timezone issues

// 	type fields struct {
// 		setupMock func(m *mocks.WALStorage)
// 	}

// 	type args struct {
// 		ctx    context.Context
// 		before time.Time
// 	}

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr error
// 	}{
// 		{
// 			name: "successful truncation",
// 			fields: fields{
// 				setupMock: func(m *mocks.WALStorage) {
// 					m.On("Truncate", context.Background(), now).Return(nil)
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				before: now,
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "truncation error",
// 			fields: fields{
// 				setupMock: func(m *mocks.WALStorage) {
// 					m.On("Truncate", context.Background(), now).Return(errors.New("truncate error"))
// 				},
// 			},
// 			args: args{
// 				ctx:    context.Background(),
// 				before: now,
// 			},
// 			wantErr: errors.New("truncate error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt

// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			mock := mocks.NewWALStorage(t)
// 			tt.fields.setupMock(mock)

// 			service := walsrv.New(mock)
// 			err := service.TruncateWAL(tt.args.ctx, tt.args.before)

// 			if tt.wantErr != nil {
// 				require.Error(t, err)
// 				require.Equal(t, tt.wantErr.Error(), err.Error())
// 			} else {
// 				require.NoError(t, err)
// 			}

// 			mock.AssertExpectations(t)
// 		})
// 	}
// }
