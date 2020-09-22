package business

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/kubernetes/kubetest"
	"github.com/kiali/kiali/models"
)

func TestGetWorkloadProxyStatusWithoutError(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	k8s := new(kubetest.K8SClientMock)
	conf := config.NewConfig()
	config.Set(conf)

	workload := "reviews-v1-982hashydas-212"
	namespace := "bookinfo"

	k8s.On("GetProxyStatus").Return([]*kubernetes.ProxyStatus{
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("%s.%s", workload, namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ClusterAcked:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteAcked:    "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
	}, nil)

	layer := ProxyStatus{k8s: k8s}
	ws := layer.GetWorkloadProxyStatus(workload, namespace, models.WorkloadStatus{})

	assert.Len(ws.ProxyStatus, 0)
}

func TestGetWorkloadProxyStatusWithError(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	k8s := new(kubetest.K8SClientMock)
	conf := config.NewConfig()
	config.Set(conf)

	workload := "reviews-v1-982hashydas-212"
	namespace := "bookinfo"

	k8s.On("GetProxyStatus").Return([]*kubernetes.ProxyStatus{
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("%s.%s", workload, namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
	}, nil)

	layer := ProxyStatus{k8s: k8s}
	ws := layer.GetWorkloadProxyStatus(workload, namespace, models.WorkloadStatus{})

	assert.Len(ws.ProxyStatus, 2)
	assert.Contains(ws.ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(ws.ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})
}

func TestGetWorkloadsProxyStatus(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	k8s := new(kubetest.K8SClientMock)
	conf := config.NewConfig()
	config.Set(conf)

	namespace := "bookinfo"

	k8s.On("GetProxyStatus").Return([]*kubernetes.ProxyStatus{
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v1-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v2-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
	}, nil)

	layer := ProxyStatus{k8s: k8s}
	wss := layer.GetWorkloadsProxyStatuses(namespace, fakeMultipleWorkloadStatus())

	assert.Len(wss, 2)
	assert.Equal("reviews-v1-982hashydas-212", wss[0].Name)
	assert.Contains(wss[0].ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(wss[0].ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})

	assert.Equal("reviews-v2-982hashydas-212", wss[1].Name)
	assert.Contains(wss[1].ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(wss[1].ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})
}

func fakeMultipleWorkloadStatus() []models.WorkloadStatus {
	return []models.WorkloadStatus{
		{
			Name:              "reviews-v1-982hashydas-212",
			DesiredReplicas:   1,
			CurrentReplicas:   1,
			AvailableReplicas: 1,
			ProxyStatus:       []models.ProxyStatus{},
		},
		{
			Name:              "reviews-v2-982hashydas-212",
			DesiredReplicas:   1,
			CurrentReplicas:   1,
			AvailableReplicas: 1,
			ProxyStatus:       []models.ProxyStatus{},
		},
	}
}

func TestGetNamespaceAppProxyStatus(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	k8s := new(kubetest.K8SClientMock)
	conf := config.NewConfig()
	config.Set(conf)

	namespace := "bookinfo"

	k8s.On("GetProxyStatus").Return([]*kubernetes.ProxyStatus{
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v1-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v2-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
	}, nil)

	layer := ProxyStatus{k8s: k8s}
	nsProxyStatuses, err := layer.GetNamespaceAppProxyStatus(namespace, fakeNamespaceAppHealth())

	assert.NoError(err)
	assert.Len(nsProxyStatuses, 1)
	assert.NotEmpty(nsProxyStatuses["reviews"])

	wss := nsProxyStatuses["reviews"].WorkloadStatuses
	assert.Equal("reviews-v1-982hashydas-212", wss[0].Name)
	assert.Contains(wss[0].ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(wss[0].ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})

	assert.Equal("reviews-v2-982hashydas-212", wss[1].Name)
	assert.Contains(wss[1].ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(wss[1].ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})
}

func fakeNamespaceAppHealth() models.NamespaceAppHealth {
	return models.NamespaceAppHealth{
		"reviews": {
			WorkloadStatuses: fakeMultipleWorkloadStatus(),
		},
	}
}

func TestGetNamespaceWorkloadProxyStatus(t *testing.T) {
	assert := assert.New(t)

	// Setup mocks
	k8s := new(kubetest.K8SClientMock)
	conf := config.NewConfig()
	config.Set(conf)

	namespace := "bookinfo"

	k8s.On("GetProxyStatus").Return([]*kubernetes.ProxyStatus{
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v1-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
		{SyncStatus: kubernetes.SyncStatus{
			ProxyID:       fmt.Sprintf("reviews-v2-982hashydas-212.%s", namespace),
			ProxyVersion:  "1.7.1",
			IstioVersion:  "1.7.1",
			ClusterSent:   "clusterdifferent",
			ClusterAcked:  "clusterequals",
			ListenerSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			ListenerAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			RouteSent:     "routedifferent",
			RouteAcked:    "",
			EndpointSent:  "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
			EndpointAcked: "zI1yscSI0RY=e9dbc143-4e20-44ec-aa5e-3c0b8097f21a",
		}},
	}, nil)

	layer := ProxyStatus{k8s: k8s}
	nsProxyStatuses, err := layer.GetNamespaceWorkloadProxyStatus(namespace, fakeNamespaceWorkloadHealth())

	assert.NoError(err)
	assert.Len(nsProxyStatuses, 1)
	assert.NotEmpty(nsProxyStatuses["reviews-v1-982hashydas-212"])

	wss := nsProxyStatuses["reviews-v1-982hashydas-212"].WorkloadStatus
	assert.Equal("reviews-v1-982hashydas-212", wss.Name)
	assert.Contains(wss.ProxyStatus, models.ProxyStatus{Component: "CDS", Status: models.Stale})
	assert.Contains(wss.ProxyStatus, models.ProxyStatus{Component: "RDS", Status: models.StaleNa})
}

func fakeNamespaceWorkloadHealth() models.NamespaceWorkloadHealth {
	return models.NamespaceWorkloadHealth{
		"reviews-v1-982hashydas-212": {
			WorkloadStatus: models.WorkloadStatus{
				Name:              "reviews-v1-982hashydas-212",
				DesiredReplicas:   1,
				CurrentReplicas:   1,
				AvailableReplicas: 1,
				ProxyStatus:       nil,
			},
		},
	}
}
