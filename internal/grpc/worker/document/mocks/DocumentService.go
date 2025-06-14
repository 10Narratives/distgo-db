// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"

	mock "github.com/stretchr/testify/mock"
)

// DocumentService is an autogenerated mock type for the DocumentService type
type DocumentService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, collection, content
func (_m *DocumentService) Create(ctx context.Context, collection string, content map[string]interface{}) (documentmodels.Document, error) {
	ret := _m.Called(ctx, collection, content)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 documentmodels.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) (documentmodels.Document, error)); ok {
		return rf(ctx, collection, content)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) documentmodels.Document); ok {
		r0 = rf(ctx, collection, content)
	} else {
		r0 = ret.Get(0).(documentmodels.Document)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, collection, content)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, collection, documentID
func (_m *DocumentService) Delete(ctx context.Context, collection string, documentID string) error {
	ret := _m.Called(ctx, collection, documentID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, collection, documentID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, collection, documentID
func (_m *DocumentService) Get(ctx context.Context, collection string, documentID string) (documentmodels.Document, error) {
	ret := _m.Called(ctx, collection, documentID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 documentmodels.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (documentmodels.Document, error)); ok {
		return rf(ctx, collection, documentID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) documentmodels.Document); ok {
		r0 = rf(ctx, collection, documentID)
	} else {
		r0 = ret.Get(0).(documentmodels.Document)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, collection, documentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, collection
func (_m *DocumentService) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
	ret := _m.Called(ctx, collection)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []documentmodels.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]documentmodels.Document, error)); ok {
		return rf(ctx, collection)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []documentmodels.Document); ok {
		r0 = rf(ctx, collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]documentmodels.Document)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, collection)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, collection, documentId, changes
func (_m *DocumentService) Update(ctx context.Context, collection string, documentId string, changes map[string]interface{}) (documentmodels.Document, error) {
	ret := _m.Called(ctx, collection, documentId, changes)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 documentmodels.Document
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]interface{}) (documentmodels.Document, error)); ok {
		return rf(ctx, collection, documentId, changes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, map[string]interface{}) documentmodels.Document); ok {
		r0 = rf(ctx, collection, documentId, changes)
	} else {
		r0 = ret.Get(0).(documentmodels.Document)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, collection, documentId, changes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDocumentService creates a new instance of DocumentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDocumentService(t interface {
	mock.TestingT
	Cleanup(func())
}) *DocumentService {
	mock := &DocumentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
