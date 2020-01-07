package models

import (
	"github.com/kiali/kiali/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AuthorizationPolicies authorizationPolicies
//
// This is used for returning an array of AuthorizationPolicies
//
// swagger:model authorizationRules
// An array of authorizationPolicy
// swagger:allOf
type AuthorizationPolicies []AuthorizationPolicy
type AuthorizationPoliciesFull []AuthorizationPolicyFull

// AuthorizationPolicy authorizationPolicy
//
// This is used for returning an AuthorizationPolicy
//
// swagger:model authorizationPolicy
type AuthorizationPolicy struct {
	meta_v1.TypeMeta
	Metadata meta_v1.ObjectMeta `json:"metadata"`
	Spec     struct {
		Selector interface{} `json:"selector"`
		Rules    interface{} `json:"rules"`
	} `json:"spec"`
}

type AuthorizationPolicyFull struct {
	meta_v1.TypeMeta
	Metadata meta_v1.ObjectMeta `json:"metadata"`
	Spec AuthorizationPolicySpec`json:"spec"`
}

type AuthorizationPolicySpec struct {
	Selector WorkloadSelector `json:"selector"`
	Rules    []Rule `json:"rules"`
}

type WorkloadSelector struct {
	MatchLabels          map[string]string `protobuf:"bytes,1,rep,name=match_labels,json=matchLabels,proto3" json:"match_labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type Rule struct {
	From From
	To To
	When When
}

type From struct {
	Source Source
}

type Source struct {
	Principals []string
	RequestPrincipals []string
	Namespaces []string
	IpBlocks []string
}

type To struct {
	Operation Operation
}

type Operation struct {
	Hosts []string
	Ports []string
	Methods []string
	Paths []string
}

type When struct {
	Condition Condition
}

type Condition struct {
	Key string
	Values []string
}

func (aps *AuthorizationPolicies) Parse(authorizationPolicies []kubernetes.IstioObject) {
	for _, authPol := range authorizationPolicies {
		ap := AuthorizationPolicy{}
		ap.Parse(authPol)
		*aps = append(*aps, ap)
	}
}

func (ap *AuthorizationPolicy) Parse(authorizationPolicy kubernetes.IstioObject) {
	ap.TypeMeta = authorizationPolicy.GetTypeMeta()
	ap.Metadata = authorizationPolicy.GetObjectMeta()
	ap.Spec.Selector = authorizationPolicy.GetSpec()["selector"]
	ap.Spec.Rules = authorizationPolicy.GetSpec()["rules"]
}
