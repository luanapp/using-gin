// Code generated by mockery v2.32.4. DO NOT EDIT.

package crud

import mock "github.com/stretchr/testify/mock"

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository[T interface{}] struct {
	mock.Mock
}

type MockRepository_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *MockRepository[T]) EXPECT() *MockRepository_Expecter[T] {
	return &MockRepository_Expecter[T]{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: id
func (_m *MockRepository[T]) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockRepository_Delete_Call[T interface{}] struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id string
func (_e *MockRepository_Expecter[T]) Delete(id interface{}) *MockRepository_Delete_Call[T] {
	return &MockRepository_Delete_Call[T]{Call: _e.mock.On("Delete", id)}
}

func (_c *MockRepository_Delete_Call[T]) Run(run func(id string)) *MockRepository_Delete_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRepository_Delete_Call[T]) Return(_a0 error) *MockRepository_Delete_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Delete_Call[T]) RunAndReturn(run func(string) error) *MockRepository_Delete_Call[T] {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields:
func (_m *MockRepository[T]) GetAll() ([]T, error) {
	ret := _m.Called()

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]T, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []T); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockRepository_GetAll_Call[T interface{}] struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *MockRepository_Expecter[T]) GetAll() *MockRepository_GetAll_Call[T] {
	return &MockRepository_GetAll_Call[T]{Call: _e.mock.On("GetAll")}
}

func (_c *MockRepository_GetAll_Call[T]) Run(run func()) *MockRepository_GetAll_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRepository_GetAll_Call[T]) Return(_a0 []T, _a1 error) *MockRepository_GetAll_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_GetAll_Call[T]) RunAndReturn(run func() ([]T, error)) *MockRepository_GetAll_Call[T] {
	_c.Call.Return(run)
	return _c
}

// GetById provides a mock function with given fields: id
func (_m *MockRepository[T]) GetById(id string) (*T, error) {
	ret := _m.Called(id)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*T, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *T); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type MockRepository_GetById_Call[T interface{}] struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id string
func (_e *MockRepository_Expecter[T]) GetById(id interface{}) *MockRepository_GetById_Call[T] {
	return &MockRepository_GetById_Call[T]{Call: _e.mock.On("GetById", id)}
}

func (_c *MockRepository_GetById_Call[T]) Run(run func(id string)) *MockRepository_GetById_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockRepository_GetById_Call[T]) Return(_a0 *T, _a1 error) *MockRepository_GetById_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_GetById_Call[T]) RunAndReturn(run func(string) (*T, error)) *MockRepository_GetById_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: entity
func (_m *MockRepository[T]) Save(entity *T) (*T, error) {
	ret := _m.Called(entity)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(*T) (*T, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(*T) *T); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(*T) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockRepository_Save_Call[T interface{}] struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - entity *T
func (_e *MockRepository_Expecter[T]) Save(entity interface{}) *MockRepository_Save_Call[T] {
	return &MockRepository_Save_Call[T]{Call: _e.mock.On("Save", entity)}
}

func (_c *MockRepository_Save_Call[T]) Run(run func(entity *T)) *MockRepository_Save_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*T))
	})
	return _c
}

func (_c *MockRepository_Save_Call[T]) Return(_a0 *T, _a1 error) *MockRepository_Save_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRepository_Save_Call[T]) RunAndReturn(run func(*T) (*T, error)) *MockRepository_Save_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: id, entity
func (_m *MockRepository[T]) Update(id string, entity *T) error {
	ret := _m.Called(id, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *T) error); ok {
		r0 = rf(id, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockRepository_Update_Call[T interface{}] struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - id string
//   - entity *T
func (_e *MockRepository_Expecter[T]) Update(id interface{}, entity interface{}) *MockRepository_Update_Call[T] {
	return &MockRepository_Update_Call[T]{Call: _e.mock.On("Update", id, entity)}
}

func (_c *MockRepository_Update_Call[T]) Run(run func(id string, entity *T)) *MockRepository_Update_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*T))
	})
	return _c
}

func (_c *MockRepository_Update_Call[T]) Return(_a0 error) *MockRepository_Update_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockRepository_Update_Call[T]) RunAndReturn(run func(string, *T) error) *MockRepository_Update_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository[T interface{}](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository[T] {
	mock := &MockRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}