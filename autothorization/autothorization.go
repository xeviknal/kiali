package autothorization

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kiali/kiali/graph/config/cytoscape"
	"github.com/kiali/kiali/models"
)

type AppName string
type TargetId string
type IncommingEdges []*cytoscape.EdgeData
type IncomingTrafficMap map[AppName]IncommingEdges
type WorkloadAppMap map[TargetId]AppName

func BuildAuthorizationGraph(ns string, g interface{}) models.AuthorizationPoliciesFull {
	graphConfig := g.(cytoscape.Config)
	itm := IncomingTrafficMap{}
	wam := WorkloadAppMap{}

	for _, n := range graphConfig.Elements.Nodes {
		wam[TargetId(n.Data.Id)] = AppName(n.Data.App)
	}

	for _, e := range graphConfig.Elements.Edges {
		// Find the app value of a Target
		an, af := wam[TargetId(e.Data.Target)]
		if !af {
			continue
		}

		// Append edges to specific app

		if _, found := itm[an]; !found {
			itm[an] = make(IncommingEdges, 0)
		}

		itm[an] = append(itm[an], e.Data)
	}

	fmt.Println(itm)

	policies := make(models.AuthorizationPoliciesFull, 0)
	for app := range itm {
		policies = append(policies, buildPolicy(string(app), itm, ns))
	}

	return policies
}

func buildPolicy(appName string, itm IncomingTrafficMap, ns string) models.AuthorizationPolicyFull {
	ap := models.AuthorizationPolicyFull {
		Metadata: v1.ObjectMeta {
			Name: fmt.Sprintf("authz-policy-%s", appName),
			Namespace: ns,
		},
		Spec: models.AuthorizationPolicySpec {
			Selector: models.WorkloadSelector {
				MatchLabels: map[string]string {
					"app": appName,
				},
			},
		},
	}
	return ap
}
