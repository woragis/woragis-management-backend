package httpserver

import (
	"net/http"
	"time"

	"github.com/woragis/management/backend/server/internal/apperrors"
	contactssvc "github.com/woragis/management/backend/server/internal/contacts/service"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
)

type agentToolsHandler struct {
	contacts    *contactssvc.Service
	personality *agentPersonalityHandler
	contactsH   *contactsHandler
	financeH    *financeHandler
	devH        *devprojectHandler
	devProjects *devprojectsvc.Service
	presenceH   *presenceHandler
	presence    *presencesvc.Service
}

func newAgentToolsHandler(app *App) *agentToolsHandler {
	h := &agentToolsHandler{
		contacts: app.Contacts,
	}
	if app.Personality != nil {
		h.personality = newAgentPersonalityHandler(app.Personality)
	}
	if app.Contacts != nil {
		h.contactsH = newContactsHandler(app.Contacts, app.Finance)
	}
	if app.Finance != nil {
		h.financeH = newFinanceHandler(app.Finance)
	}
	if app.DevProjects != nil {
		h.devH = newDevprojectHandler(app.DevProjects)
		h.devProjects = app.DevProjects
	}
	if app.Presence != nil {
		h.presence = app.Presence
		h.presenceH = newPresenceHandler(app.Presence)
	}
	return h
}

func (h *agentToolsHandler) getPersonality(w http.ResponseWriter, r *http.Request) {
	if h.personality == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Personality service unavailable."))
		return
	}
	h.personality.get(w, r)
}

func (h *agentToolsHandler) updatePersonality(w http.ResponseWriter, r *http.Request) {
	if h.personality == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Personality service unavailable."))
		return
	}
	h.personality.update(w, r)
}

func (h *agentToolsHandler) resetPersonality(w http.ResponseWriter, r *http.Request) {
	if h.personality == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Personality service unavailable."))
		return
	}
	h.personality.reset(w, r)
}

func (h *agentToolsHandler) searchContacts(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.list(w, r)
}

func (h *agentToolsHandler) getContact(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.get(w, r)
}

func (h *agentToolsHandler) createContact(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.create(w, r)
}

func (h *agentToolsHandler) updateContact(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.update(w, r)
}

func (h *agentToolsHandler) logInteraction(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.createInteraction(w, r)
}

func (h *agentToolsHandler) listContactsDueFollowUp(w http.ResponseWriter, r *http.Request) {
	if h.contacts == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	before := time.Now().UTC()
	if v := r.URL.Query().Get("before"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			before = t.UTC()
		}
	}
	rows, err := h.contacts.ListDueFollowUp(r.Context(), before)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *agentToolsHandler) getContactFinance(w http.ResponseWriter, r *http.Request) {
	if h.contactsH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Contacts service unavailable."))
		return
	}
	h.contactsH.contactFinance(w, r)
}

func (h *agentToolsHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	if h.devH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Projects service unavailable."))
		return
	}
	h.devH.list(w, r)
}

func (h *agentToolsHandler) getProject(w http.ResponseWriter, r *http.Request) {
	if h.devH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Projects service unavailable."))
		return
	}
	h.devH.get(w, r)
}

func (h *agentToolsHandler) createProject(w http.ResponseWriter, r *http.Request) {
	if h.devH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Projects service unavailable."))
		return
	}
	h.devH.create(w, r)
}

func (h *agentToolsHandler) financeDashboard(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.dashboard(w, r)
}

func (h *agentToolsHandler) financeSummary(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.summary(w, r)
}

func (h *agentToolsHandler) financeCalendar(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.calendar(w, r)
}

func (h *agentToolsHandler) listIncomeSources(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.listIncomeSources(w, r)
}

func (h *agentToolsHandler) listTransactions(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.listTransactions(w, r)
}

func (h *agentToolsHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.createTransaction(w, r)
}

func (h *agentToolsHandler) createIncomeSource(w http.ResponseWriter, r *http.Request) {
	if h.financeH == nil {
		apperrors.WriteError(w, apperrors.InternalErr(apperrors.CodeInternal, "Finance service unavailable."))
		return
	}
	h.financeH.createIncomeSource(w, r)
}
