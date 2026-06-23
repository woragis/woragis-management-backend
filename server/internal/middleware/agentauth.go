package middleware

import (
	"net/http"
	"strings"

	"github.com/woragis/management/backend/server/internal/apperrors"
)

const headerAgentKey = "X-Agent-Key"

func AgentAuth(expectedKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(expectedKey) == "" {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeAgentAuthInvalid, apperrors.MsgAgentAuthInvalid))
			return
		}
		key := strings.TrimSpace(r.Header.Get(headerAgentKey))
		if key == "" {
			auth := strings.TrimSpace(r.Header.Get("Authorization"))
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				key = strings.TrimSpace(auth[7:])
			}
		}
		if key != expectedKey {
			apperrors.WriteError(w, apperrors.Unauthorized(apperrors.CodeAgentAuthInvalid, apperrors.MsgAgentAuthInvalid))
			return
		}
		next.ServeHTTP(w, r)
	})
}
