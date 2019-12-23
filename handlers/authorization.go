package handlers

import (
	"net/http"

	"github.com/kiali/kiali/autothorization"
	"github.com/kiali/kiali/graph"
	"github.com/kiali/kiali/graph/api"
)

func Authorization(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)

	o := graph.NewOptions(r)
	o.ConfigOptions.GraphType = graph.GraphTypeWorkload
	o.TelemetryOptions.InjectServiceNodes = false

	business, err := getBusiness(r)
	graph.CheckError(err)

	code, payload := api.GraphNamespaces(business, o)
	autothorization.BuildAuthorizationGraph(payload)

	respond(w, code, payload)
}
