package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/linksAPI/db"
	"github.com/hillview.tv/linksAPI/query"
)

func CheckLinkRouteHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	route := params["route"]

	if len(route) == 0 {
		http.Error(w, "missing route param", http.StatusBadRequest)
		return
	}

	routeFound, err := query.LookupRoute(db.DB, route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routeFound)
}
