package handlers

import (
	"net/http"

	"github.com/kiali/kiali/graph"
)

func RankingLoad(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	// Take default options from graph
	o := graph.NewOptions(r)
	loadSummary, err := business.Rankings.Load(o)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Kiali could not load rankings")
	}
	RespondWithJSON(w, http.StatusOK, loadSummary)
}

func RankingSlowEdges(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	// Take default options from graph
	slowEdges := business.Rankings.GetSlowEdges(5)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Kiali could not fetch slow edges ranking")
	}
	RespondWithJSON(w, http.StatusOK, slowEdges)
}

func RankingMostTimeConsumingEdges(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	// Take default options from graph
	slowEdges := business.Rankings.GetMostTimeConsumingEdges(5)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Kiali could not fetch most time consuming edges ranking")
	}
	RespondWithJSON(w, http.StatusOK, slowEdges)
}
