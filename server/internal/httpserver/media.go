package httpserver

import (
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
)

type mediaHandler struct {
	svc *mediasvc.Service
}

func newMediaHandler(svc *mediasvc.Service) *mediaHandler {
	return &mediaHandler{svc: svc}
}

func (h *mediaHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.List(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, items)
}

func (h *mediaHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeMediaGetV1HandlerPathIDInvalid, apperrors.MsgMediaGetV1HandlerPathIDInvalid))
		return
	}
	item, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, item)
}

func (h *mediaHandler) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(105 << 20); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, "Invalid multipart form."))
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeMediaPostV1ServiceFileMissing, apperrors.MsgMediaPostV1ServiceFileMissing))
		return
	}
	defer func() { _ = file.Close() }()

	item, err := h.svc.Upload(r.Context(), mediasvc.UploadInput{
		Filename: header.Filename,
		MimeType: header.Header.Get("Content-Type"),
		AltText:  r.FormValue("altText"),
		Reader:   file,
	})
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, item)
}

func (h *mediaHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeMediaGetV1HandlerPathIDInvalid, apperrors.MsgMediaGetV1HandlerPathIDInvalid))
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *mediaHandler) serveFile(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeMediaGetV1HandlerPathIDInvalid, apperrors.MsgMediaGetV1HandlerPathIDInvalid))
		return
	}
	asset, f, err := h.svc.Open(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	defer func() { _ = f.Close() }()
	w.Header().Set("Content-Type", asset.MimeType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	_, _ = io.Copy(w, f)
}
