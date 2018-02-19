package models

import (
	"k8s.io/api/core/v1"
)

type Namespace struct {
	Name string `json:"name"`
}

func MarshallCollection(nsl *v1.NamespaceList) []Namespace {
	namespaces := make([]Namespace, len(nsl.Items))
	for i, item := range nsl.Items {
		namespaces[i] = Marshall(item)
	}

	return namespaces
}

func Marshall(ns v1.Namespace) Namespace {
	namespace := Namespace{}
	namespace.Name = ns.Name

	return namespace
}
