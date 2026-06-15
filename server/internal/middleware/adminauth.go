package middleware

import (
	"net/http"
	"strings"

	"github.com/woragis/management/backend/server/internal/apperrors"
)

const headerAdminKey = "X-Admin-Key"

func AdminAuth(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimSpace(r.Header.Get(headerAdminKey))
		if key == "" {
			auth := strings.TrimSpace(r.Header.Get("Authorization"))
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				key = strings.TrimSpace(auth[7:])
			}
		}
		if key == "" {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeAdminAuthV1HandlerKeyMissing, apperrors.MsgAdminAuthV1HandlerKeyMissing))
			return
		}
		if key != expectedKey {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeAdminAuthV1HandlerKeyInvalid, apperrors.MsgAdminAuthV1HandlerKeyInvalid))
			return
		}
		next.ServeHTTP(w, r)
	})
}
