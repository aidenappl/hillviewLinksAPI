package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hillview.tv/linksAPI/db"
	"github.com/hillview.tv/linksAPI/query"
	"github.com/hillview.tv/linksAPI/responder"
)

func CheckLinkRouteHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	route := params["route"]

	recordClick := r.URL.Query().Get("recordClick")

	if len(route) == 0 {
		responder.ErrMissingBodyRequirement(w, "route")
		return
	}

	routeFound, err := query.LookupRoute(db.DB, route)
	if err != nil {
		responder.ErrInternal(w, err, "failed to lookup route")
		return
	}

	if routeFound == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if recordClick == "true" {
		err := query.RecordClick(db.DB, query.RecordClickRequest{
			LinkID: &routeFound.ID,
		})
		if err != nil {
			responder.ErrInternal(w, err, "failed to record click")
			return
		}
	}

	responder.New(w, routeFound, "successfully found route")
}
