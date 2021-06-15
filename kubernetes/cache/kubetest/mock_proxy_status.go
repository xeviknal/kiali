package kubetest

import "github.com/kiali/kiali/kubernetes"

func (m *CacheMock) CheckProxyStatus() bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *CacheMock) GetPodProxyStatus(namespace, pod string) (*kubernetes.ProxyStatus, error) {
	args := m.Called(namespace, pod)
	return args.Get(0).(*kubernetes.ProxyStatus), args.Get(0).(error)
}