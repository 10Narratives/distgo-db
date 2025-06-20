// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	mock "github.com/stretchr/testify/mock"
)

// WALStorage is an autogenerated mock type for the WALStorage type
type WALStorage struct {
	mock.Mock
}

// LogEntry provides a mock function with given fields: ctx, entry
func (_m *WALStorage) LogEntry(ctx context.Context, entry walmodels.WALEntry) error {
	ret := _m.Called(ctx, entry)

	if len(ret) == 0 {
		panic("no return value specified for LogEntry")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, walmodels.WALEntry) error); ok {
		r0 = rf(ctx, entry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWALStorage creates a new instance of WALStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWALStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *WALStorage {
	mock := &WALStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
