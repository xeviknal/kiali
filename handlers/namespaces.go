package handlers

import (
	"net/http"

	"github.com/swift-sunshine/swscore/kubernetes"
	"github.com/swift-sunshine/swscore/log"
)

func NamespaceList(w http.ResponseWriter, r *http.Request) {
	istioClient, err := kubernetes.NewClient()
	if err != nil {
		log.Error(err)
		RespondWithError(w, 500, err.Error())
		return
	}

	namespaces, err := istioClient.GetNamespaces()
	if err != nil {
		log.Error(err)
		RespondWithError(w, 500, err.Error())
		return
	}

	RespondWithJSON(w, 200, namespaces)
}
