package models

import (
	"github.com/swift-sunshine/swscore/kubernetes"
	"k8s.io/api/core/v1"
)

type ServiceOverview struct {
	Name string `json:"name"`
}

type ServiceList struct {
	Namespace Namespace         `json:"namespace"`
	Service   []ServiceOverview `json:"services"`
}

type Service struct {
	Name      string            `json:"name"`
	Namespace Namespace         `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Type      string            `json:"type"`
	Ip        string            `json:"ip"`
	Ports     []Port            `json:"ports"`
	Endpoints Endpoints         `json:"endpoints"`
	Pods      Pods              `json:"pods"`
}

func GetServicesByNamespace(namespaceName string) ([]ServiceOverview, error) {
	istioClient, err := kubernetes.NewClient()
	if err != nil {
		return nil, err
	}

	services, err := istioClient.GetServices(namespaceName)
	if err != nil {
		return nil, err
	}

	return CastServiceOverviewCollection(services), nil
}

func GetServiceDetails(namespaceName, serviceName string) (*Service, error) {
	istioClient, err := kubernetes.NewClient()
	if err != nil {
		return nil, err
	}

	serviceDetails, err := istioClient.GetServiceDetails(namespaceName, serviceName)
	if err != nil {
		return nil, err
	}

	return CastService(serviceDetails), nil
}

func CastServiceOverviewCollection(sl *v1.ServiceList) []ServiceOverview {
	services := make([]ServiceOverview, len(sl.Items))
	for i, item := range sl.Items {
		services[i] = CastServiceOverview(item)
	}

	return services
}

func CastServiceOverview(s v1.Service) ServiceOverview {
	service := ServiceOverview{}
	service.Name = s.Name

	return service
}

func CastService(s *kubernetes.ServiceDetails) *Service {
	service := &Service{}
	service.Name = s.Service.Name
	service.Namespace = Namespace{s.Service.Namespace}
	service.Labels = s.Service.Labels
	service.Type = string(s.Service.Spec.Type)
	service.Ip = s.Service.Spec.ClusterIP
	service.Ports = CastServicePortCollection(s.Service.Spec.Ports)
	(&service.Endpoints).Parse(s.Endpoints)
	(&service.Pods).Parse(s.Pods)

	return service
}
