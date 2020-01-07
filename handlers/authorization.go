package handlers

import (
	"net/http"

	"github.com/kiali/kiali/autothorization"
	"github.com/kiali/kiali/graph"
	"github.com/kiali/kiali/graph/api"
)

type AuthorizationResponse struct {
	Payload interface{}
	Policies interface{}
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	o := graph.NewOptions(r)
	o.ConfigOptions.GraphType = graph.GraphTypeWorkload
	o.TelemetryOptions.InjectServiceNodes = false

	business, err := getBusiness(r)
	graph.CheckError(err)

	response := AuthorizationResponse{}
	code, payload := api.GraphNamespaces(business, o)
	policies := autothorization.BuildAuthorizationGraph(o.Namespace, payload)
	response.Payload = payload
	response.Policies = policies

	respond(w, code, response)
}
