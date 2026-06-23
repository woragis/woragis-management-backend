package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
)

type devprojectHandler struct {
	svc *devprojectsvc.Service
}

func newDevprojectHandler(svc *devprojectsvc.Service) *devprojectHandler {
	return &devprojectHandler{svc: svc}
}

func (h *devprojectHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var isPublic, featured *bool
	if v := q.Get("isPublic"); v == "true" {
		b := true
		isPublic = &b
	} else if v == "false" {
		b := false
		isPublic = &b
	}
	if v := q.Get("featured"); v == "true" {
		b := true
		featured = &b
	} else if v == "false" {
		b := false
		featured = &b
	}
	filter := devprojectsvc.ListFilter{
		Status:         q.Get("status"),
		Intent:         q.Get("intent"),
		Monetization:   q.Get("monetization"),
		Maturity:       q.Get("maturity"),
		VisibilityGoal: q.Get("visibilityGoal"),
		Distribution:   q.Get("distribution"),
		IsPublic:       isPublic,
		Featured:       featured,
		Query:          q.Get("q"),
	}
	if filter.Status != "" || filter.Intent != "" || filter.Monetization != "" || filter.Maturity != "" ||
		filter.VisibilityGoal != "" || filter.Distribution != "" || filter.IsPublic != nil || filter.Featured != nil || filter.Query != "" {
		items, err := h.svc.ListFiltered(r.Context(), filter)
		if err != nil {
			apperrors.WriteError(w, err)
			return
		}
		apperrors.WriteJSON(w, http.StatusOK, items)
		return
	}
	items, err := h.svc.List(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, items)
}

func (h *devprojectHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, p)
}

func (h *devprojectHandler) create(w http.ResponseWriter, r *http.Request) {
	var body createProjectBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceNameEmpty, "Request body is invalid."))
		return
	}
	p, err := h.svc.Create(r.Context(), body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, p)
}

func (h *devprojectHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body updateProjectBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceNameEmpty, "Request body is invalid."))
		return
	}
	p, err := h.svc.Update(r.Context(), id, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, p)
}

func (h *devprojectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *devprojectHandler) createLink(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body createLinkBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectLinkPostV1ServiceURLInvalid, "Request body is invalid."))
		return
	}
	link, err := h.svc.CreateLink(r.Context(), projectID, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, link)
}

func (h *devprojectHandler) deleteLink(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	linkID, err := parseUUID(r.PathValue("linkId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.DeleteLink(r.Context(), projectID, linkID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *devprojectHandler) createDomain(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body createDomainBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectDomainPostV1ServiceDomainInvalid, "Request body is invalid."))
		return
	}
	d, err := h.svc.CreateDomain(r.Context(), projectID, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, d)
}

func (h *devprojectHandler) deleteDomain(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	domainID, err := parseUUID(r.PathValue("domainId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.DeleteDomain(r.Context(), projectID, domainID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *devprojectHandler) listSecrets(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	items, err := h.svc.ListSecrets(r.Context(), projectID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, items)
}

func (h *devprojectHandler) getSecret(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	secretID, err := parseUUID(r.PathValue("secretId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	item, err := h.svc.GetSecret(r.Context(), projectID, secretID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, item)
}

func (h *devprojectHandler) createSecret(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body createSecretBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectSecretPostV1ServiceValueEmpty, "Request body is invalid."))
		return
	}
	item, err := h.svc.CreateSecret(r.Context(), projectID, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, item)
}

func (h *devprojectHandler) deleteSecret(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	secretID, err := parseUUID(r.PathValue("secretId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.DeleteSecret(r.Context(), projectID, secretID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *devprojectHandler) createGallery(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body createGalleryBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGalleryPostV1ServiceMediaInvalid, "Request body is invalid."))
		return
	}
	item, err := h.svc.CreateGalleryItem(r.Context(), projectID, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, item)
}

func (h *devprojectHandler) deleteGallery(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	itemID, err := parseUUID(r.PathValue("itemId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.DeleteGalleryItem(r.Context(), projectID, itemID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *devprojectHandler) listEnvs(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	items, err := h.svc.ListEnvs(r.Context(), projectID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, items)
}

func (h *devprojectHandler) createEnv(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	var body createEnvBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectEnvPostV1ServiceKeyEmpty, "Request body is invalid."))
		return
	}
	item, err := h.svc.CreateEnv(r.Context(), projectID, body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, item)
}

func (h *devprojectHandler) deleteEnv(w http.ResponseWriter, r *http.Request) {
	projectID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	envID, err := parseUUID(r.PathValue("envId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProjectGetV1HandlerPathIDInvalid, apperrors.MsgProjectGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.DeleteEnv(r.Context(), projectID, envID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createProjectBody struct {
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Description      string     `json:"description"`
	ShortDescription string     `json:"shortDescription"`
	LongDescription  string     `json:"longDescription"`
	Status           string     `json:"status"`
	Intent           string     `json:"intent"`
	Distribution     []string   `json:"distribution"`
	Monetization     string     `json:"monetization"`
	Maturity         string     `json:"maturity"`
	VisibilityGoal   string     `json:"visibilityGoal"`
	Stack            []string   `json:"stack"`
	RepoURL          string     `json:"repoUrl"`
	DemoURL          string     `json:"demoUrl"`
	GithubURL        string     `json:"githubUrl"`
	RepoVisibility   string     `json:"repoVisibility"`
	Notes            string     `json:"notes"`
	IsPublic         bool       `json:"isPublic"`
	Featured         bool       `json:"featured"`
	DisplayOrder     int        `json:"displayOrder"`
	PublicSlug       string     `json:"publicSlug"`
	CoverImageID     *uuid.UUID `json:"coverImageId"`
	ParentProjectID  *uuid.UUID `json:"parentProjectId"`
}

func (b createProjectBody) toInput() devprojectsvc.CreateProjectInput {
	return devprojectsvc.CreateProjectInput(b)
}

type updateProjectBody struct {
	Name             *string    `json:"name"`
	Slug             *string    `json:"slug"`
	Description      *string    `json:"description"`
	ShortDescription *string    `json:"shortDescription"`
	LongDescription  *string    `json:"longDescription"`
	Status           *string    `json:"status"`
	Intent           *string    `json:"intent"`
	Distribution     []string   `json:"distribution"`
	Monetization     *string    `json:"monetization"`
	Maturity         *string    `json:"maturity"`
	VisibilityGoal   *string    `json:"visibilityGoal"`
	Stack            []string   `json:"stack"`
	RepoURL          *string    `json:"repoUrl"`
	DemoURL          *string    `json:"demoUrl"`
	GithubURL        *string    `json:"githubUrl"`
	RepoVisibility   *string    `json:"repoVisibility"`
	Notes            *string    `json:"notes"`
	IsPublic         *bool      `json:"isPublic"`
	Featured         *bool      `json:"featured"`
	DisplayOrder     *int       `json:"displayOrder"`
	PublicSlug       *string    `json:"publicSlug"`
	CoverImageID     *uuid.UUID `json:"coverImageId"`
	ParentProjectID  *uuid.UUID `json:"parentProjectId"`
}

func (b updateProjectBody) toInput() devprojectsvc.UpdateProjectInput {
	in := devprojectsvc.UpdateProjectInput{
		Name:             b.Name,
		Slug:             b.Slug,
		Description:      b.Description,
		ShortDescription: b.ShortDescription,
		LongDescription:  b.LongDescription,
		Status:           b.Status,
		Intent:           b.Intent,
		Monetization:     b.Monetization,
		Maturity:         b.Maturity,
		VisibilityGoal:   b.VisibilityGoal,
		RepoURL:          b.RepoURL,
		DemoURL:          b.DemoURL,
		GithubURL:        b.GithubURL,
		RepoVisibility:   b.RepoVisibility,
		Notes:            b.Notes,
		IsPublic:         b.IsPublic,
		Featured:         b.Featured,
		DisplayOrder:     b.DisplayOrder,
		PublicSlug:       b.PublicSlug,
	}
	if b.Stack != nil {
		in.Stack = b.Stack
		in.StackSet = true
	}
	if b.Distribution != nil {
		in.Distribution = b.Distribution
		in.DistributionSet = true
	}
	if b.CoverImageID != nil {
		in.CoverImageID = b.CoverImageID
		in.CoverImageSet = true
	}
	if b.ParentProjectID != nil {
		in.ParentProjectID = b.ParentProjectID
		in.ParentProjectSet = true
	}
	return in
}

type createLinkBody struct {
	Type        string `json:"type"`
	URL         string `json:"url"`
	Environment string `json:"environment"`
	Label       string `json:"label"`
	IsPublic    bool   `json:"isPublic"`
}

func (b createLinkBody) toInput() devprojectsvc.CreateLinkInput {
	return devprojectsvc.CreateLinkInput(b)
}

type createDomainBody struct {
	Domain    string  `json:"domain"`
	Registrar string  `json:"registrar"`
	Purpose   string  `json:"purpose"`
	ExpiresAt *string `json:"expiresAt"`
	Notes     string  `json:"notes"`
}

func (b createDomainBody) toInput() devprojectsvc.CreateDomainInput {
	in := devprojectsvc.CreateDomainInput{
		Domain:    b.Domain,
		Registrar: b.Registrar,
		Purpose:   b.Purpose,
		Notes:     b.Notes,
	}
	if b.ExpiresAt != nil && *b.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *b.ExpiresAt); err == nil {
			in.ExpiresAt = &t
		}
	}
	return in
}

type createSecretBody struct {
	Name        string  `json:"name"`
	Value       string  `json:"value"`
	Environment string  `json:"environment"`
	Service     string  `json:"service"`
	ExpiresAt   *string `json:"expiresAt"`
	Notes       string  `json:"notes"`
}

func (b createSecretBody) toInput() devprojectsvc.CreateSecretInput {
	in := devprojectsvc.CreateSecretInput{
		Name:        b.Name,
		Value:       b.Value,
		Environment: b.Environment,
		Service:     b.Service,
		Notes:       b.Notes,
	}
	if b.ExpiresAt != nil && *b.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, *b.ExpiresAt); err == nil {
			in.ExpiresAt = &t
		}
	}
	return in
}

type createGalleryBody struct {
	MediaAssetID uuid.UUID `json:"mediaAssetId"`
	DisplayOrder int       `json:"displayOrder"`
	Caption      string    `json:"caption"`
}

func (b createGalleryBody) toInput() devprojectsvc.CreateGalleryInput {
	return devprojectsvc.CreateGalleryInput(b)
}

type createEnvBody struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Environment string `json:"environment"`
	Notes       string `json:"notes"`
}

func (b createEnvBody) toInput() devprojectsvc.CreateEnvInput {
	return devprojectsvc.CreateEnvInput(b)
}

func parseUUID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, err
	}
	return id, nil
}
