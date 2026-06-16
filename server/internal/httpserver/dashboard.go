package httpserver

import (
	"net/http"

	"github.com/woragis/management/backend/server/internal/apperrors"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	mediarepo "github.com/woragis/management/backend/server/internal/media/repository"
)

func handleDashboard(projects *devprojectsvc.Service, media *mediarepo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var counter devprojectsvc.MediaCounter
		if media != nil {
			counter = media
		}
		d, err := projects.Dashboard(r.Context(), counter)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, d)
	}
}
