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

	recordClick := r.URL.Query().Get("recordClick")

	if len(route) == 0 {
		http.Error(w, "missing route param", http.StatusBadRequest)
		return
	}

	routeFound, err := query.LookupRoute(db.DB, route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if routeFound == nil {
		json.NewEncoder(w).Encode(routeFound)
		return
	}

	if recordClick == "true" {
		err := query.RecordClick(db.DB, query.RecordClickRequest{
			LinkID: &routeFound.ID,
		})
		if err != nil {
			http.Error(w, "failed to record click: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(routeFound)
}
