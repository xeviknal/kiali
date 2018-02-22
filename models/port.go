package models

import "k8s.io/api/core/v1"

type Port struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     int32  `json:"port"`
}

func CastPortCollection(ps []v1.ServicePort) []Port {
	ports := make([]Port, len(ps))
	for i, servicePort := range ps {
		ports[i] = CastPort(servicePort)
	}

	return ports
}

func CastPort(p v1.ServicePort) Port {
	port := Port{}
	port.Name = p.Name
	port.Protocol = string(p.Protocol)
	port.Port = p.Port

	return port
}
