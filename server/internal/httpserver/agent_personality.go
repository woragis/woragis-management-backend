package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/woragis/management/backend/server/internal/apperrors"
	personalitysvc "github.com/woragis/management/backend/server/internal/agent/personality/service"
)

type agentPersonalityHandler struct {
	svc *personalitysvc.Service
}

func newAgentPersonalityHandler(svc *personalitysvc.Service) *agentPersonalityHandler {
	return &agentPersonalityHandler{svc: svc}
}

func (h *agentPersonalityHandler) get(w http.ResponseWriter, r *http.Request) {
	row, err := h.svc.Get(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *agentPersonalityHandler) update(w http.ResponseWriter, r *http.Request) {
	var body personalityUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.Update(r.Context(), body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *agentPersonalityHandler) reset(w http.ResponseWriter, r *http.Request) {
	row, err := h.svc.Reset(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

type personalityUpdateBody struct {
	AssistantName     *string `json:"assistantName"`
	GreetingMorning   *string `json:"greetingMorning"`
	GreetingAfternoon *string `json:"greetingAfternoon"`
	GreetingEvening   *string `json:"greetingEvening"`
	GreetingEnabled   *bool   `json:"greetingEnabled"`
	SystemPromptExtra *string `json:"systemPromptExtra"`
	VoiceID           *string `json:"voiceId"`
	Language          *string `json:"language"`
	Timezone          *string `json:"timezone"`
}

func (b personalityUpdateBody) toInput() personalitysvc.UpdateInput {
	return personalitysvc.UpdateInput(b)
}
