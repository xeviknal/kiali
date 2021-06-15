package business

import (
	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/kubernetes/cache/kubetest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCached_GetPodProxyStatus(t *testing.T) {
	assert := assert.New(t)

	kialiCache := kubetest.NewCacheMock()
	kialiCache.On("CheckProxyStatus").Return(true)
	kialiCache.On("GetPodProxyStatus", "bookinfo", "details-v1-25042120-jcs").Return(&kubernetes.ProxyStatus{
		SyncStatus: kubernetes.SyncStatus {
			ProxyID:       "istiod-250421-jcs",
			ProxyVersion:  "manual-version",
			IstioVersion:  "v1.10-rc",
			ClusterSent:   "abc",
			ClusterAcked:  "abc",
			ListenerSent:  "def",
			ListenerAcked: "def",
			RouteSent:     "ghi",
			RouteAcked:    "ghi",
			EndpointSent:  "jkl",
			EndpointAcked: "jkl",
		},
	})

	conf := config.NewConfig()
	config.Set(conf)

	proxyStatusService := ProxyStatusService{k8s: nil, businessLayer: NewWithBackends(nil, nil, nil)}
	ps, err := proxyStatusService.GetPodProxyStatus("bookinfo", "details-v1-25042120-jcs")

	assert.NotNil(ps)
	assert.Equal("istiod-250421-jcs", ps.ProxyID)
	assert.NoError(err)
}
