// Code generated by mockery v2.9.4. DO NOT EDIT.

package crud

import (
	"github.com/luanapp/gin-example/pkg/model"

	mock "github.com/stretchr/testify/mock"
)

// repositoryMock is an autogenerated mock type for the repositorier type
type repositoryMock struct {
	mock.Mock
}

// delete provides a mock function with given fields: id
func (_m *repositoryMock) delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// getAll provides a mock function with given fields:
func (_m *repositoryMock) getAll() ([]model.Species, error) {
	ret := _m.Called()

	var r0 []model.Species
	if rf, ok := ret.Get(0).(func() []model.Species); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Species)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// getById provides a mock function with given fields: id
func (_m *repositoryMock) getById(id string) (*model.Species, error) {
	ret := _m.Called(id)

	var r0 *model.Species
	if rf, ok := ret.Get(0).(func(string) *model.Species); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Species)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// save provides a mock function with given fields: sp
func (_m *repositoryMock) save(sp *model.Species) (*model.Species, error) {
	ret := _m.Called(sp)

	var r0 *model.Species
	if rf, ok := ret.Get(0).(func(*model.Species) *model.Species); ok {
		r0 = rf(sp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Species)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Species) error); ok {
		r1 = rf(sp)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// update provides a mock function with given fields: sp
func (_m *repositoryMock) update(sp *model.Species) error {
	ret := _m.Called(sp)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Species) error); ok {
		r0 = rf(sp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
