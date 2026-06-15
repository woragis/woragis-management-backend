package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/models"
	profilesvc "github.com/woragis/management/backend/server/internal/profile/service"
)

type profileHandler struct {
	svc *profilesvc.Service
}

func newProfileHandler(svc *profilesvc.Service) *profileHandler {
	return &profileHandler{svc: svc}
}

func (h *profileHandler) getAdmin(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.Get(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, p)
}

func (h *profileHandler) update(w http.ResponseWriter, r *http.Request) {
	var body updateProfileBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeProfilePatchV1ServiceDisplayNameEmpty, "Request body is invalid."))
		return
	}
	p, err := h.svc.Update(r.Context(), body.toInput())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, p)
}

func (h *profileHandler) getPublic(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetPublic(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=300")
	apperrors.WriteJSON(w, http.StatusOK, p)
}

type updateProfileBody struct {
	DisplayName   *string             `json:"displayName"`
	Headline      *string             `json:"headline"`
	Bio           *string             `json:"bio"`
	AvatarID      *uuid.UUID          `json:"avatarId"`
	Location      *string             `json:"location"`
	Availability  *string             `json:"availability"`
	ResumeAssetID *uuid.UUID          `json:"resumeAssetId"`
	SocialLinks   []models.SocialLink `json:"socialLinks"`
}

func (b updateProfileBody) toInput() profilesvc.UpdateInput {
	in := profilesvc.UpdateInput{
		DisplayName:  b.DisplayName,
		Headline:     b.Headline,
		Bio:          b.Bio,
		Location:     b.Location,
		Availability: b.Availability,
	}
	if b.AvatarID != nil {
		in.AvatarID = b.AvatarID
		in.AvatarSet = true
	}
	if b.ResumeAssetID != nil {
		in.ResumeAssetID = b.ResumeAssetID
		in.ResumeSet = true
	}
	if b.SocialLinks != nil {
		in.SocialLinks = b.SocialLinks
		in.SocialSet = true
	}
	return in
}
