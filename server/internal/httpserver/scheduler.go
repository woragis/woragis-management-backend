package httpserver

import (
	"net/http"
	"time"

	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/messaging/executor"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
)

func handleSchedulerDue(svc *messagingsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := svc.ListDueJobs(r.Context(), time.Now().UTC())
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, rows)
	}
}

func handleSchedulerExecute(exec *executor.Executor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseUUID(r.PathValue("id"))
		if err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
			return
		}
		result, err := exec.ExecuteJob(r.Context(), id)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, result)
	}
}
