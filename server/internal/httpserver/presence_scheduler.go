package httpserver

import (
	"net/http"
	"time"

	"github.com/woragis/management/backend/server/internal/apperrors"
	presencereminder "github.com/woragis/management/backend/server/internal/presence/reminder"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
)

func handlePresenceDueReminders(svc *presencesvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := svc.ListDueReminders(r.Context(), time.Now().UTC(), 50)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, rows)
	}
}

func handlePresenceSendReminder(exec *presencereminder.Executor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseUUID(r.PathValue("id"))
		if err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
			return
		}
		result, err := exec.SendReminder(r.Context(), id)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, result)
	}
}
