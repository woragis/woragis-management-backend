package httpserver

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
)

type contentHandler struct {
	svc       *contentsvc.Service
	messaging *messagingsvc.Service
}

func newContentHandler(svc *contentsvc.Service, messaging *messagingsvc.Service) *contentHandler {
	return &contentHandler{svc: svc, messaging: messaging}
}

func (h *contentHandler) listVideos(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListVideos(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contentHandler) getVideo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetVideo(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) createVideo(w http.ResponseWriter, r *http.Request) {
	var body leetcodeVideoBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateVideo(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *contentHandler) updateVideo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body leetcodeVideoUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateVideo(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) deleteVideo(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteVideo(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *contentHandler) listThumbnails(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	rows, err := h.svc.ListThumbnails(r.Context(), videoID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contentHandler) getThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	thumbID, err := parseUUID(r.PathValue("thumbnailId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid thumbnail id."))
		return
	}
	row, err := h.svc.GetThumbnail(r.Context(), videoID, thumbID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) createThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	var body thumbnailBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateThumbnail(r.Context(), videoID, body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *contentHandler) updateThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	thumbID, err := parseUUID(r.PathValue("thumbnailId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid thumbnail id."))
		return
	}
	var body thumbnailUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateThumbnail(r.Context(), videoID, thumbID, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) deleteThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	thumbID, err := parseUUID(r.PathValue("thumbnailId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid thumbnail id."))
		return
	}
	if err := h.svc.DeleteThumbnail(r.Context(), videoID, thumbID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *contentHandler) generateThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	thumbID, err := parseUUID(r.PathValue("thumbnailId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid thumbnail id."))
		return
	}
	row, err := h.svc.GenerateThumbnail(r.Context(), videoID, thumbID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) approveThumbnail(w http.ResponseWriter, r *http.Request) {
	videoID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid video id."))
		return
	}
	thumbID, err := parseUUID(r.PathValue("thumbnailId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid thumbnail id."))
		return
	}
	row, err := h.svc.ApproveThumbnail(r.Context(), videoID, thumbID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) listTemplates(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListTemplates(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *contentHandler) getTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetTemplate(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) createTemplate(w http.ResponseWriter, r *http.Request) {
	var body templateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateTemplate(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *contentHandler) updateTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body templateUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateTemplate(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *contentHandler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteTemplate(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type leetcodeVideoBody struct {
	Title                 string   `json:"title"`
	Status                string   `json:"status"`
	SeriesNumber          *int     `json:"seriesNumber"`
	TrackName             *string  `json:"trackName"`
	ProblemTitle          *string  `json:"problemTitle"`
	LeetcodeProblemNumber *int     `json:"leetcodeProblemNumber"`
	LeetcodeSlug          *string  `json:"leetcodeSlug"`
	StudyPlanSlug         *string  `json:"studyPlanSlug"`
	Difficulty            *string  `json:"difficulty"`
	Topics                []string `json:"topics"`
	ShortDescription      *string  `json:"shortDescription"`
	LeetcodeProblemURL    *string  `json:"leetcodeProblemUrl"`
	LeetcodeSubmissionURL *string  `json:"leetcodeSubmissionUrl"`
	Notes                 *string  `json:"notes"`
	YoutubeURL            *string  `json:"youtubeUrl"`
	ProblemDate           *string  `json:"problemDate"`
	WhatsappEnabled       *bool    `json:"whatsappEnabled"`
}

func parseProblemDate(s *string) *time.Time {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*s))
	if err != nil {
		return nil
	}
	return &t
}

func (b leetcodeVideoBody) toCreate() contentsvc.CreateVideoInput {
	return contentsvc.CreateVideoInput{
		Title:                 b.Title,
		Status:                b.Status,
		SeriesNumber:          b.SeriesNumber,
		TrackName:             b.TrackName,
		ProblemTitle:          b.ProblemTitle,
		LeetcodeProblemNumber: b.LeetcodeProblemNumber,
		LeetcodeSlug:          b.LeetcodeSlug,
		StudyPlanSlug:         b.StudyPlanSlug,
		Difficulty:            b.Difficulty,
		Topics:                b.Topics,
		ShortDescription:      b.ShortDescription,
		LeetcodeProblemURL:    b.LeetcodeProblemURL,
		LeetcodeSubmissionURL: b.LeetcodeSubmissionURL,
		Notes:                 b.Notes,
		YoutubeURL:            b.YoutubeURL,
		ProblemDate:           parseProblemDate(b.ProblemDate),
		WhatsappEnabled:       b.WhatsappEnabled,
	}
}

type leetcodeVideoUpdateBody struct {
	Title                 *string  `json:"title"`
	Status                *string  `json:"status"`
	SeriesNumber          *int     `json:"seriesNumber"`
	SeriesNumberSet       bool     `json:"seriesNumberSet"`
	TrackName             *string  `json:"trackName"`
	ProblemTitle          *string  `json:"problemTitle"`
	LeetcodeProblemNumber *int     `json:"leetcodeProblemNumber"`
	LeetcodeProblemSet    bool     `json:"leetcodeProblemSet"`
	LeetcodeSlug          *string  `json:"leetcodeSlug"`
	StudyPlanSlug         *string  `json:"studyPlanSlug"`
	Difficulty            *string  `json:"difficulty"`
	Topics                []string `json:"topics"`
	TopicsSet             bool     `json:"topicsSet"`
	ShortDescription      *string  `json:"shortDescription"`
	LeetcodeProblemURL    *string  `json:"leetcodeProblemUrl"`
	LeetcodeSubmissionURL *string  `json:"leetcodeSubmissionUrl"`
	Notes                 *string  `json:"notes"`
	YoutubeURL            *string  `json:"youtubeUrl"`
	ProblemDate           *string  `json:"problemDate"`
	ProblemDateSet        bool     `json:"problemDateSet"`
	WhatsappEnabled       *bool    `json:"whatsappEnabled"`
}

func (b leetcodeVideoUpdateBody) toUpdate() contentsvc.UpdateVideoInput {
	var problemDate *time.Time
	if b.ProblemDateSet {
		problemDate = parseProblemDate(b.ProblemDate)
	}
	return contentsvc.UpdateVideoInput{
		Title:                 b.Title,
		Status:                b.Status,
		SeriesNumber:          b.SeriesNumber,
		SeriesNumberSet:       b.SeriesNumberSet,
		TrackName:             b.TrackName,
		ProblemTitle:          b.ProblemTitle,
		LeetcodeProblemNumber: b.LeetcodeProblemNumber,
		LeetcodeProblemSet:    b.LeetcodeProblemSet,
		LeetcodeSlug:          b.LeetcodeSlug,
		StudyPlanSlug:         b.StudyPlanSlug,
		Difficulty:            b.Difficulty,
		Topics:                b.Topics,
		TopicsSet:             b.TopicsSet,
		ShortDescription:      b.ShortDescription,
		LeetcodeProblemURL:    b.LeetcodeProblemURL,
		LeetcodeSubmissionURL: b.LeetcodeSubmissionURL,
		Notes:                 b.Notes,
		YoutubeURL:            b.YoutubeURL,
		ProblemDate:           problemDate,
		ProblemDateSet:        b.ProblemDateSet,
		WhatsappEnabled:       b.WhatsappEnabled,
	}
}

type thumbnailBody struct {
	Prompt            string   `json:"prompt"`
	NegativePrompt    *string  `json:"negativePrompt"`
	Size              string   `json:"size"`
	Quality           string   `json:"quality"`
	Model             string   `json:"model"`
	Mode              string   `json:"mode"`
	ReferenceMediaIDs []string `json:"referenceMediaIds"`
}

func (b thumbnailBody) toCreate() contentsvc.CreateThumbnailInput {
	return contentsvc.CreateThumbnailInput{
		Prompt:            b.Prompt,
		NegativePrompt:    b.NegativePrompt,
		Size:              b.Size,
		Quality:           b.Quality,
		Model:             b.Model,
		Mode:              b.Mode,
		ReferenceMediaIDs: parseUUIDList(b.ReferenceMediaIDs),
	}
}

type thumbnailUpdateBody struct {
	Prompt            *string  `json:"prompt"`
	NegativePrompt    *string  `json:"negativePrompt"`
	Size              *string  `json:"size"`
	Quality           *string  `json:"quality"`
	Model             *string  `json:"model"`
	Mode              *string  `json:"mode"`
	ReferenceMediaIDs []string `json:"referenceMediaIds"`
	ReferenceSet      bool     `json:"referenceSet"`
}

func (b thumbnailUpdateBody) toUpdate() contentsvc.UpdateThumbnailInput {
	return contentsvc.UpdateThumbnailInput{
		Prompt:            b.Prompt,
		NegativePrompt:    b.NegativePrompt,
		Size:              b.Size,
		Quality:           b.Quality,
		Model:             b.Model,
		Mode:              b.Mode,
		ReferenceMediaIDs: parseUUIDList(b.ReferenceMediaIDs),
		ReferenceSet:      b.ReferenceSet,
	}
}

type templateBody struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	PromptTemplate string `json:"promptTemplate"`
	IsDefault      bool   `json:"isDefault"`
}

func (b templateBody) toCreate() contentsvc.CreateTemplateInput {
	return contentsvc.CreateTemplateInput{
		Name:           b.Name,
		Slug:           b.Slug,
		PromptTemplate: b.PromptTemplate,
		IsDefault:      b.IsDefault,
	}
}

type templateUpdateBody struct {
	Name           *string `json:"name"`
	Slug           *string `json:"slug"`
	PromptTemplate *string `json:"promptTemplate"`
	IsDefault      *bool   `json:"isDefault"`
}

func (b templateUpdateBody) toUpdate() contentsvc.UpdateTemplateInput {
	return contentsvc.UpdateTemplateInput{
		Name:           b.Name,
		Slug:           b.Slug,
		PromptTemplate: b.PromptTemplate,
		IsDefault:      b.IsDefault,
	}
}

func parseUUIDList(items []string) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(items))
	for _, s := range items {
		if id, err := uuid.Parse(s); err == nil {
			out = append(out, id)
		}
	}
	return out
}

func handleCreativesWebhook(svc *contentsvc.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			JobID string `json:"jobId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
			return
		}
		jobID, err := uuid.Parse(body.JobID)
		if err != nil {
			apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid job id."))
			return
		}
		if err := svc.HandleCreativesWebhook(r.Context(), jobID); err != nil {
			apperrors.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
