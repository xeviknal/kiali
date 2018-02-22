package models

import "k8s.io/api/core/v1"

type Ports []Port
type Port struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     int32  `json:"port"`
}

func CastServicePortCollection(ps []v1.ServicePort) []Port {
	ports := make([]Port, len(ps))
	for i, servicePort := range ps {
		ports[i] = CastServicePort(servicePort)
	}

	return ports
}

func CastServicePort(p v1.ServicePort) Port {
	port := Port{}
	port.Name = p.Name
	port.Protocol = string(p.Protocol)
	port.Port = p.Port

	return port
}

func (ports *Ports) ParseEndpointPorts(ps []v1.EndpointPort) {
	for _, endpointPort := range ps {
		port := Port{}
		port.ParseEndpointPort(endpointPort)
		*ports = append(*ports, port)
	}
}

func (port *Port) ParseEndpointPort(p v1.EndpointPort) {
	port.Name = p.Name
	port.Protocol = string(p.Protocol)
	port.Port = p.Port
}
