package service

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/content/repository"
	"github.com/woragis/management/backend/server/internal/creativesclient"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
	"github.com/woragis/management/backend/server/internal/models"
)

type Service struct {
	repo             *repository.Repository
	media            *mediasvc.Service
	creatives        *creativesclient.Client
	whatsappWorker   *whatsappworkerclient.Client
	webhookURL       string
	defaultThumbSize string
	msgTemplates     MessagingTemplateLookup
}

// MessagingTemplateLookup resolves unified MessageTemplate bodies (optional).
type MessagingTemplateLookup interface {
	TemplateBodyForProgram(ctx context.Context, programSlug, slug string) (string, bool)
}

func New(repo *repository.Repository, media *mediasvc.Service, creatives *creativesclient.Client, whatsappWorker *whatsappworkerclient.Client, webhookURL, defaultThumbSize string) *Service {
	if defaultThumbSize == "" {
		defaultThumbSize = "1280x720"
	}
	return &Service{
		repo:             repo,
		media:            media,
		creatives:        creatives,
		whatsappWorker:   whatsappWorker,
		webhookURL:       strings.TrimSpace(webhookURL),
		defaultThumbSize: defaultThumbSize,
	}
}

func (s *Service) SetMessagingTemplates(lookup MessagingTemplateLookup) {
	s.msgTemplates = lookup
}

type CreateVideoInput struct {
	Title                 string
	Status                string
	SeriesNumber          *int
	TrackName             *string
	ProblemTitle          *string
	LeetcodeProblemNumber *int
	LeetcodeSlug          *string
	StudyPlanSlug         *string
	Difficulty            *string
	Topics                []string
	ShortDescription      *string
	LeetcodeProblemURL    *string
	LeetcodeSubmissionURL *string
	Notes                 *string
	YoutubeURL            *string
	ProblemDate           *time.Time
	WhatsappEnabled       *bool
}

type UpdateVideoInput struct {
	Title                 *string
	Status                *string
	SeriesNumber          *int
	SeriesNumberSet       bool
	TrackName             *string
	ProblemTitle          *string
	LeetcodeProblemNumber *int
	LeetcodeProblemSet    bool
	LeetcodeSlug          *string
	StudyPlanSlug         *string
	Difficulty            *string
	Topics                []string
	TopicsSet             bool
	ShortDescription      *string
	LeetcodeProblemURL    *string
	LeetcodeSubmissionURL *string
	Notes                 *string
	YoutubeURL            *string
	ProblemDate           *time.Time
	ProblemDateSet        bool
	WhatsappEnabled       *bool
}

func (s *Service) ListVideos(ctx context.Context) ([]models.LeetcodeVideo, error) {
	rows, err := s.repo.ListVideos(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "list videos", err)
	}
	return rows, nil
}

func (s *Service) GetVideo(ctx context.Context, id uuid.UUID) (*models.LeetcodeVideo, error) {
	video, err := s.repo.GetVideo(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := range video.Thumbnails {
		if err := s.syncThumbnail(ctx, &video.Thumbnails[i]); err != nil {
			return nil, err
		}
	}
	return video, nil
}

func (s *Service) CreateVideo(ctx context.Context, in CreateVideoInput) (*models.LeetcodeVideo, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, apperrors.Invalid(apperrors.CodeContentVideoCreateInvalid, apperrors.MsgContentVideoCreateInvalid)
	}
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = models.LeetcodeVideoStatusDraft
	}
	var topicsJSON json.RawMessage
	if len(in.Topics) > 0 {
		topicsJSON, _ = json.Marshal(in.Topics)
	}
	row := &models.LeetcodeVideo{
		Title:                 title,
		Status:                status,
		SeriesNumber:          in.SeriesNumber,
		TrackName:             trimPtr(in.TrackName),
		ProblemTitle:          trimPtr(in.ProblemTitle),
		LeetcodeProblemNumber: in.LeetcodeProblemNumber,
		LeetcodeSlug:          trimPtr(in.LeetcodeSlug),
		StudyPlanSlug:         trimPtr(in.StudyPlanSlug),
		Difficulty:            trimPtr(in.Difficulty),
		Topics:                topicsJSON,
		ShortDescription:      trimPtr(in.ShortDescription),
		LeetcodeProblemURL:    trimPtr(in.LeetcodeProblemURL),
		LeetcodeSubmissionURL: trimPtr(in.LeetcodeSubmissionURL),
		Notes:                 trimPtr(in.Notes),
		YoutubeURL:            trimPtr(in.YoutubeURL),
		ProblemDate:           in.ProblemDate,
		WhatsappEnabled:       true,
	}
	if in.WhatsappEnabled != nil {
		row.WhatsappEnabled = *in.WhatsappEnabled
	}
	if err := s.repo.CreateVideo(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "create video", err)
	}
	return row, nil
}

func (s *Service) UpdateVideo(ctx context.Context, id uuid.UUID, in UpdateVideoInput) (*models.LeetcodeVideo, error) {
	row, err := s.repo.GetVideo(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return nil, apperrors.Invalid(apperrors.CodeContentVideoCreateInvalid, apperrors.MsgContentVideoCreateInvalid)
		}
		row.Title = title
	}
	if in.Status != nil {
		row.Status = strings.TrimSpace(*in.Status)
	}
	if in.SeriesNumberSet {
		row.SeriesNumber = in.SeriesNumber
	}
	if in.TrackName != nil {
		row.TrackName = trimPtr(in.TrackName)
	}
	if in.ProblemTitle != nil {
		row.ProblemTitle = trimPtr(in.ProblemTitle)
	}
	if in.LeetcodeProblemSet {
		row.LeetcodeProblemNumber = in.LeetcodeProblemNumber
	}
	if in.LeetcodeSlug != nil {
		row.LeetcodeSlug = trimPtr(in.LeetcodeSlug)
	}
	if in.StudyPlanSlug != nil {
		row.StudyPlanSlug = trimPtr(in.StudyPlanSlug)
	}
	if in.Difficulty != nil {
		row.Difficulty = trimPtr(in.Difficulty)
	}
	if in.TopicsSet {
		if len(in.Topics) > 0 {
			row.Topics, _ = json.Marshal(in.Topics)
		} else {
			row.Topics = nil
		}
	}
	if in.Notes != nil {
		row.Notes = trimPtr(in.Notes)
	}
	if in.ShortDescription != nil {
		row.ShortDescription = trimPtr(in.ShortDescription)
	}
	if in.LeetcodeProblemURL != nil {
		row.LeetcodeProblemURL = trimPtr(in.LeetcodeProblemURL)
	}
	if in.LeetcodeSubmissionURL != nil {
		row.LeetcodeSubmissionURL = trimPtr(in.LeetcodeSubmissionURL)
	}
	if in.YoutubeURL != nil {
		row.YoutubeURL = trimPtr(in.YoutubeURL)
	}
	if in.ProblemDateSet {
		row.ProblemDate = in.ProblemDate
	}
	if in.WhatsappEnabled != nil {
		row.WhatsappEnabled = *in.WhatsappEnabled
	}
	if err := s.repo.UpdateVideo(ctx, row); err != nil {
		return nil, err
	}
	return s.repo.GetVideo(ctx, id)
}

func (s *Service) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteVideo(ctx, id)
}

type CreateThumbnailInput struct {
	Prompt             string
	NegativePrompt     *string
	Size               string
	Quality            string
	Model              string
	Mode               string
	ReferenceMediaIDs  []uuid.UUID
}

type UpdateThumbnailInput struct {
	Prompt            *string
	NegativePrompt    *string
	Size              *string
	Quality           *string
	Model             *string
	Mode              *string
	ReferenceMediaIDs []uuid.UUID
	ReferenceSet      bool
}

func (s *Service) ListThumbnails(ctx context.Context, videoID uuid.UUID) ([]models.ContentThumbnail, error) {
	if _, err := s.repo.GetVideo(ctx, videoID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListThumbnails(ctx, videoID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "list thumbnails", err)
	}
	for i := range rows {
		if err := s.syncThumbnail(ctx, &rows[i]); err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func (s *Service) GetThumbnail(ctx context.Context, videoID, id uuid.UUID) (*models.ContentThumbnail, error) {
	row, err := s.repo.GetThumbnail(ctx, videoID, id)
	if err != nil {
		return nil, err
	}
	if err := s.syncThumbnail(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) CreateThumbnail(ctx context.Context, videoID uuid.UUID, in CreateThumbnailInput) (*models.ContentThumbnail, error) {
	if _, err := s.repo.GetVideo(ctx, videoID); err != nil {
		return nil, err
	}
	prompt := strings.TrimSpace(in.Prompt)
	if prompt == "" {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailCreateInvalid, apperrors.MsgContentThumbnailCreateInvalid)
	}
	size := strings.TrimSpace(in.Size)
	if size == "" {
		size = s.defaultThumbSize
	}
	quality := strings.TrimSpace(in.Quality)
	if quality == "" {
		quality = "high"
	}
	model := strings.TrimSpace(in.Model)
	if model == "" {
		model = "gpt-image-2"
	}
	mode := strings.TrimSpace(in.Mode)
	if mode == "" {
		mode = "edit"
	}
	refJSON, _ := json.Marshal(uuidStrings(in.ReferenceMediaIDs))
	row := &models.ContentThumbnail{
		VideoID:           videoID,
		Status:            models.ContentThumbnailStatusDraft,
		Prompt:            prompt,
		NegativePrompt:    trimPtr(in.NegativePrompt),
		Size:              size,
		Quality:           quality,
		Model:             model,
		Mode:              mode,
		ReferenceMediaIDs: refJSON,
	}
	if err := s.repo.CreateThumbnail(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "create thumbnail", err)
	}
	return row, nil
}

func (s *Service) UpdateThumbnail(ctx context.Context, videoID, id uuid.UUID, in UpdateThumbnailInput) (*models.ContentThumbnail, error) {
	row, err := s.repo.GetThumbnail(ctx, videoID, id)
	if err != nil {
		return nil, err
	}
	if row.Status == models.ContentThumbnailStatusGenerating {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailGenerateInvalid, "cannot edit thumbnail while generating")
	}
	if in.Prompt != nil {
		prompt := strings.TrimSpace(*in.Prompt)
		if prompt == "" {
			return nil, apperrors.Invalid(apperrors.CodeContentThumbnailCreateInvalid, apperrors.MsgContentThumbnailCreateInvalid)
		}
		row.Prompt = prompt
	}
	if in.NegativePrompt != nil {
		row.NegativePrompt = trimPtr(in.NegativePrompt)
	}
	if in.Size != nil {
		row.Size = strings.TrimSpace(*in.Size)
	}
	if in.Quality != nil {
		row.Quality = strings.TrimSpace(*in.Quality)
	}
	if in.Model != nil {
		row.Model = strings.TrimSpace(*in.Model)
	}
	if in.Mode != nil {
		row.Mode = strings.TrimSpace(*in.Mode)
	}
	if in.ReferenceSet {
		row.ReferenceMediaIDs, _ = json.Marshal(uuidStrings(in.ReferenceMediaIDs))
	}
	if err := s.repo.UpdateThumbnail(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteThumbnail(ctx context.Context, videoID, id uuid.UUID) error {
	return s.repo.DeleteThumbnail(ctx, videoID, id)
}

func (s *Service) GenerateThumbnail(ctx context.Context, videoID, id uuid.UUID) (*models.ContentThumbnail, error) {
	if s.creatives == nil || !s.creatives.Enabled() {
		return nil, apperrors.Invalid(apperrors.CodeCreativesClientUnavailable, apperrors.MsgCreativesClientUnavailable)
	}
	row, err := s.repo.GetThumbnail(ctx, videoID, id)
	if err != nil {
		return nil, err
	}
	if row.Status != models.ContentThumbnailStatusDraft && row.Status != models.ContentThumbnailStatusFailed && row.Status != models.ContentThumbnailStatusReady {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailGenerateInvalid, apperrors.MsgContentThumbnailGenerateInvalid)
	}

	refURLs, err := s.referenceURLs(ctx, row.ReferenceMediaIDs)
	if err != nil {
		return nil, err
	}
	if row.Mode == "edit" && len(refURLs) == 0 {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailGenerateInvalid, "edit mode requires reference images")
	}

	var webhook *string
	if s.webhookURL != "" {
		webhook = &s.webhookURL
	}
	job, err := s.creatives.CreateContentImage(ctx, creativesclient.CreateContentImageRequest{
		ExternalID:    row.ID.String(),
		Mode:          row.Mode,
		Model:         row.Model,
		Prompt:        row.Prompt,
		Size:          row.Size,
		Quality:       row.Quality,
		ReferenceURLs: refURLs,
		WebhookURL:    webhook,
		Metadata: map[string]any{
			"channel": "leetcode",
			"videoId": videoID.String(),
		},
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "creatives create content image", err)
	}

	row.Status = models.ContentThumbnailStatusGenerating
	row.CreativesJobID = &job.ID
	row.ErrorMessage = nil
	row.OutputMediaID = nil
	if err := s.repo.UpdateThumbnail(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ApproveThumbnail(ctx context.Context, videoID, id uuid.UUID) (*models.ContentThumbnail, error) {
	row, err := s.GetThumbnail(ctx, videoID, id)
	if err != nil {
		return nil, err
	}
	if row.Status != models.ContentThumbnailStatusReady {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailGenerateInvalid, "thumbnail must be ready before approval")
	}
	if row.OutputMediaID == nil {
		return nil, apperrors.Invalid(apperrors.CodeContentThumbnailGenerateInvalid, "thumbnail has no output media")
	}
	row.Status = models.ContentThumbnailStatusApproved
	if err := s.repo.UpdateThumbnail(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) HandleCreativesWebhook(ctx context.Context, jobID uuid.UUID) error {
	thumb, err := s.repo.FindThumbnailByCreativesJob(ctx, jobID)
	if err != nil {
		return err
	}
	return s.syncThumbnail(ctx, thumb)
}

func (s *Service) syncThumbnail(ctx context.Context, row *models.ContentThumbnail) error {
	if row.Status != models.ContentThumbnailStatusGenerating || row.CreativesJobID == nil {
		return nil
	}
	if s.creatives == nil || !s.creatives.Enabled() {
		return nil
	}
	job, err := s.creatives.GetContentImage(ctx, *row.CreativesJobID)
	if err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "creatives get content image", err)
	}
	switch job.Status {
	case "queued", "running":
		return nil
	case "failed", "cancelled":
		msg := "generation failed"
		if job.ErrorMessage != nil {
			msg = *job.ErrorMessage
		}
		row.Status = models.ContentThumbnailStatusFailed
		row.ErrorMessage = &msg
		return s.repo.UpdateThumbnail(ctx, row)
	case "completed":
		if job.OutputURL == nil {
			msg := "creatives returned no output url"
			row.Status = models.ContentThumbnailStatusFailed
			row.ErrorMessage = &msg
			return s.repo.UpdateThumbnail(ctx, row)
		}
		data, err := s.creatives.DownloadOutput(ctx, *job.OutputURL)
		if err != nil {
			return apperrors.InternalCause(apperrors.CodeInternal, "download creatives output", err)
		}
		asset, err := s.media.Upload(ctx, mediasvc.UploadInput{
			Filename: "thumbnail.png",
			MimeType: "image/png",
			AltText:  "LeetCode thumbnail",
			Reader:   bytes.NewReader(data),
		})
		if err != nil {
			return err
		}
		row.Status = models.ContentThumbnailStatusReady
		row.OutputMediaID = &asset.ID
		row.ErrorMessage = nil
		return s.repo.UpdateThumbnail(ctx, row)
	default:
		return nil
	}
}

func (s *Service) referenceURLs(ctx context.Context, raw json.RawMessage) ([]string, error) {
	var ids []string
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &ids)
	}
	var urls []string
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		asset, err := s.media.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		urls = append(urls, asset.PublicURL)
	}
	return urls, nil
}

type CreateTemplateInput struct {
	Name           string
	Slug           string
	PromptTemplate string
	IsDefault      bool
}

type UpdateTemplateInput struct {
	Name           *string
	Slug           *string
	PromptTemplate *string
	IsDefault      *bool
}

func (s *Service) ListTemplates(ctx context.Context) ([]models.ContentPromptTemplate, error) {
	rows, err := s.repo.ListTemplates(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "list templates", err)
	}
	return rows, nil
}

func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (*models.ContentPromptTemplate, error) {
	return s.repo.GetTemplate(ctx, id)
}

func (s *Service) CreateTemplate(ctx context.Context, in CreateTemplateInput) (*models.ContentPromptTemplate, error) {
	name := strings.TrimSpace(in.Name)
	slug := strings.TrimSpace(in.Slug)
	prompt := strings.TrimSpace(in.PromptTemplate)
	if name == "" || slug == "" || prompt == "" {
		return nil, apperrors.Invalid(apperrors.CodeContentTemplateCreateInvalid, apperrors.MsgContentTemplateCreateInvalid)
	}
	if in.IsDefault {
		_ = s.repo.ClearDefaultTemplates(ctx)
	}
	row := &models.ContentPromptTemplate{
		ChannelSlug:    "leetcode",
		Name:           name,
		Slug:           slug,
		PromptTemplate: prompt,
		IsDefault:      in.IsDefault,
	}
	if err := s.repo.CreateTemplate(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "create template", err)
	}
	return row, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, id uuid.UUID, in UpdateTemplateInput) (*models.ContentPromptTemplate, error) {
	row, err := s.repo.GetTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		row.Name = strings.TrimSpace(*in.Name)
	}
	if in.Slug != nil {
		row.Slug = strings.TrimSpace(*in.Slug)
	}
	if in.PromptTemplate != nil {
		row.PromptTemplate = strings.TrimSpace(*in.PromptTemplate)
	}
	if in.IsDefault != nil && *in.IsDefault {
		_ = s.repo.ClearDefaultTemplates(ctx)
		row.IsDefault = true
	} else if in.IsDefault != nil {
		row.IsDefault = *in.IsDefault
	}
	if row.Name == "" || row.Slug == "" || row.PromptTemplate == "" {
		return nil, apperrors.Invalid(apperrors.CodeContentTemplateCreateInvalid, apperrors.MsgContentTemplateCreateInvalid)
	}
	if err := s.repo.UpdateTemplate(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTemplate(ctx, id)
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return &v
}

func uuidStrings(ids []uuid.UUID) []string {
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		out = append(out, id.String())
	}
	return out
}
