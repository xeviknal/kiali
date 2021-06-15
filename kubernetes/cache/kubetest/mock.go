package kubetest

import "github.com/stretchr/testify/mock"

type CacheMock struct {
	mock.Mock
}

// Constructor

func NewCacheMock() *CacheMock {
	cache := new(CacheMock)
	cache.On("CheckNamespace", mock.AnythingOfType("string")).Return(true)
	cache.On("RefreshNamespace", mock.AnythingOfType("string"))
	cache.On("Stop", mock.AnythingOfType("string"))
	return cache
}

func (m *CacheMock) CheckNamespace(namespace string) bool {
	args := m.Called(namespace)
	return args.Get(0).(bool)
}

func (m *CacheMock) RefreshNamespace(namespace string) {
	m.Called(namespace)
}

func (m *CacheMock) Stop(namespace string) {
	m.Called(namespace)
}
