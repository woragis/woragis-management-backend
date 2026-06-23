package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/content/templaterender"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
)

func defaultWhatsappTemplates() []models.WhatsappMessageTemplate {
	return []models.WhatsappMessageTemplate{
		{
			ChannelSlug: "leetcode",
			Slug:        models.WhatsappTplProblemDaily,
			Name:        "Problema do dia",
			Body: `🚀 *PROBLEMA DO DIA #{{number}}*

📚 {{track}}
🎯 {{problemTitle}}
⭐ {{difficulty}}

🔗 {{leetcodeUrl}}

📝 *Desafio*
{{shortDescription}}

⏰ Solução em vídeo será enviada hoje à noite.
Tentem resolver antes de assistir.

Compartilhem ideias, complexidade e dúvidas.
Bons estudos 🚀`,
		},
		{
			ChannelSlug: "leetcode",
			Slug:        models.WhatsappTplDiscussionNudge,
			Name:        "Estímulo à discussão",
			Body: `🤔 *Como vocês resolveram o problema de hoje?*

• Qual a complexidade da sua solução?
• Dá para resolver em uma passagem?
• Como explicaria em uma entrevista?

Compartilhem antes da solução oficial 👇`,
		},
		{
			ChannelSlug: "leetcode",
			Slug:        models.WhatsappTplSolutionVideo,
			Name:        "Solução em vídeo",
			Body: `🎥 *SOLUÇÃO DISPONÍVEL*

{{track}} — Problema #{{number}}
*{{problemTitle}}*

▶️ {{youtubeUrl}}

Neste vídeo:
✅ Entendimento do problema
✅ Lógica passo a passo
✅ Implementação
✅ Complexidade
✅ Dicas para entrevista

Ficou claro? Resolveram sozinhos?`,
		},
		{
			ChannelSlug: "leetcode",
			Slug:        models.WhatsappTplWeeklySummary,
			Name:        "Resumo semanal",
			Body: `📊 *RESUMO DA SEMANA*

Problemas:
{{problemList}}

Próxima semana: {{nextTheme}}

Consistência > intensidade 💪`,
		},
		{
			ChannelSlug: "leetcode",
			Slug:        models.WhatsappTplGroupInvite,
			Name:        "Convite ao grupo",
			Body: `🚀 *LeetCode & Programação Avançada*

Grupo focado em estruturas de dados, algoritmos, LeetCode e entrevistas técnicas.

*Como funciona:*
• 1 problema por dia
• Discussão livre durante o dia
• Solução em vídeo à noite

Entre aqui: {{inviteLink}}`,
		},
	}
}

func (s *Service) EnsureWhatsappDefaults(ctx context.Context) error {
	if _, err := s.repo.EnsureSettings(ctx); err != nil {
		return err
	}
	for _, tpl := range defaultWhatsappTemplates() {
		exists, err := s.repo.WhatsappTemplateExists(ctx, tpl.Slug)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		row := tpl
		if err := s.repo.CreateWhatsappTemplate(ctx, &row); err != nil {
			return err
		}
	}
	return nil
}

type UpdateSettingsInput struct {
	Timezone             *string
	ProblemPostTime      *string
	DiscussionPostTime   *string
	SolutionPostTime     *string
	WeeklySummaryDay     *string
	WeeklySummaryTime    *string
	DiscussionEnabled    *bool
	InviteLink           *string
	DefaultStudyPlanSlug *string
	NextTheme            *string
}

func (s *Service) GetSettings(ctx context.Context) (*models.LeetcodeChannelSettings, error) {
	return s.repo.EnsureSettings(ctx)
}

func (s *Service) UpdateSettings(ctx context.Context, in UpdateSettingsInput) (*models.LeetcodeChannelSettings, error) {
	row, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, err
	}
	if in.Timezone != nil {
		row.Timezone = strings.TrimSpace(*in.Timezone)
	}
	if in.ProblemPostTime != nil {
		row.ProblemPostTime = strings.TrimSpace(*in.ProblemPostTime)
	}
	if in.DiscussionPostTime != nil {
		row.DiscussionPostTime = strings.TrimSpace(*in.DiscussionPostTime)
	}
	if in.SolutionPostTime != nil {
		row.SolutionPostTime = strings.TrimSpace(*in.SolutionPostTime)
	}
	if in.WeeklySummaryDay != nil {
		row.WeeklySummaryDay = strings.TrimSpace(*in.WeeklySummaryDay)
	}
	if in.WeeklySummaryTime != nil {
		row.WeeklySummaryTime = strings.TrimSpace(*in.WeeklySummaryTime)
	}
	if in.DiscussionEnabled != nil {
		row.DiscussionEnabled = *in.DiscussionEnabled
	}
	if in.InviteLink != nil {
		row.InviteLink = trimPtr(in.InviteLink)
	}
	if in.DefaultStudyPlanSlug != nil {
		row.DefaultStudyPlanSlug = trimPtr(in.DefaultStudyPlanSlug)
	}
	if in.NextTheme != nil {
		row.NextTheme = trimPtr(in.NextTheme)
	}
	if err := s.repo.UpdateSettings(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "update settings", err)
	}
	return row, nil
}

type WhatsappTemplateInput struct {
	Name string
	Slug string
	Body string
}

func (s *Service) ListWhatsappTemplates(ctx context.Context) ([]models.WhatsappMessageTemplate, error) {
	rows, err := s.repo.ListWhatsappTemplates(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "list whatsapp templates", err)
	}
	return rows, nil
}

func (s *Service) UpdateWhatsappTemplate(ctx context.Context, id uuid.UUID, in WhatsappTemplateInput) (*models.WhatsappMessageTemplate, error) {
	row, err := s.repo.GetWhatsappTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != "" {
		row.Name = strings.TrimSpace(in.Name)
	}
	if in.Body != "" {
		row.Body = in.Body
	}
	if err := s.repo.UpdateWhatsappTemplate(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

type DispatchResult struct {
	VideoID      string `json:"videoId,omitempty"`
	TemplateSlug string `json:"templateSlug"`
	Message      string `json:"message"`
	Skip         bool   `json:"skip"`
	SkipReason   string `json:"skipReason,omitempty"`
}

func (s *Service) Dispatch(ctx context.Context, dispatchType, dateStr string) (*DispatchResult, error) {
	settings, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, err
	}
	day := templaterender.TodayInTZ(settings.Timezone)
	if dateStr != "" {
		day, err = templaterender.ParseDateInTZ(dateStr, settings.Timezone)
		if err != nil {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid date.")
		}
	}

	switch dispatchType {
	case "problem":
		return s.dispatchProblem(ctx, settings, day)
	case "discussion":
		return s.dispatchDiscussion(ctx, settings, day)
	case "solution":
		return s.dispatchSolution(ctx, settings, day)
	case "weekly":
		return s.dispatchWeekly(ctx, settings, day)
	default:
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid dispatch type.")
	}
}

func (s *Service) dispatchProblem(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchResult, error) {
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplProblemDaily}, nil
	}
	if video.WhatsappProblemSentAt != nil {
		return &DispatchResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplProblemDaily, VideoID: video.ID.String()}, nil
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, models.WhatsappTplProblemDaily)
	if err != nil {
		return nil, err
	}
	msg := templaterender.Render(tpl.Body, templaterender.FromVideo(video, settings))
	return &DispatchResult{VideoID: video.ID.String(), TemplateSlug: tpl.Slug, Message: msg}, nil
}

func (s *Service) dispatchDiscussion(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchResult, error) {
	if !settings.DiscussionEnabled {
		return &DispatchResult{Skip: true, SkipReason: "discussion disabled", TemplateSlug: models.WhatsappTplDiscussionNudge}, nil
	}
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplDiscussionNudge}, nil
	}
	if video.WhatsappDiscussionSentAt != nil {
		return &DispatchResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplDiscussionNudge, VideoID: video.ID.String()}, nil
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, models.WhatsappTplDiscussionNudge)
	if err != nil {
		return nil, err
	}
	msg := templaterender.Render(tpl.Body, templaterender.FromVideo(video, settings))
	return &DispatchResult{VideoID: video.ID.String(), TemplateSlug: tpl.Slug, Message: msg}, nil
}

func (s *Service) dispatchSolution(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchResult, error) {
	video, err := s.repo.GetVideoByProblemDate(ctx, day)
	if err != nil {
		return &DispatchResult{Skip: true, SkipReason: "no video for date", TemplateSlug: models.WhatsappTplSolutionVideo}, nil
	}
	if video.WhatsappSolutionSentAt != nil {
		return &DispatchResult{Skip: true, SkipReason: "already sent", TemplateSlug: models.WhatsappTplSolutionVideo, VideoID: video.ID.String()}, nil
	}
	if video.YoutubeURL == nil || strings.TrimSpace(*video.YoutubeURL) == "" {
		return &DispatchResult{Skip: true, SkipReason: "youtube url missing", TemplateSlug: models.WhatsappTplSolutionVideo, VideoID: video.ID.String()}, nil
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, models.WhatsappTplSolutionVideo)
	if err != nil {
		return nil, err
	}
	msg := templaterender.Render(tpl.Body, templaterender.FromVideo(video, settings))
	return &DispatchResult{VideoID: video.ID.String(), TemplateSlug: tpl.Slug, Message: msg}, nil
}

func (s *Service) dispatchWeekly(ctx context.Context, settings *models.LeetcodeChannelSettings, day time.Time) (*DispatchResult, error) {
	loc, _ := time.LoadLocation(settings.Timezone)
	if loc == nil {
		loc = time.UTC
	}
	day = day.In(loc)
	weekday := int(day.Weekday())
	start := day.AddDate(0, 0, -weekday)
	end := start.AddDate(0, 0, 7)
	videos, err := s.repo.ListVideosWithProblemSentBetween(ctx, start, end)
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return &DispatchResult{Skip: true, SkipReason: "no problems this week", TemplateSlug: models.WhatsappTplWeeklySummary}, nil
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, models.WhatsappTplWeeklySummary)
	if err != nil {
		return nil, err
	}
	vars := templaterender.Vars{ProblemList: templaterender.FormatProblemList(videos)}
	if settings.NextTheme != nil {
		vars.NextTheme = strings.TrimSpace(*settings.NextTheme)
	}
	msg := templaterender.Render(tpl.Body, vars)
	return &DispatchResult{TemplateSlug: tpl.Slug, Message: msg}, nil
}

func (s *Service) PreviewWhatsapp(ctx context.Context, videoID uuid.UUID, templateSlug string) (string, error) {
	video, err := s.repo.GetVideo(ctx, videoID)
	if err != nil {
		return "", err
	}
	settings, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return "", err
	}
	if templateSlug == "" {
		templateSlug = models.WhatsappTplProblemDaily
	}
	tpl, err := s.repo.GetWhatsappTemplateBySlug(ctx, templateSlug)
	if err != nil {
		return "", err
	}
	vars := templaterender.FromVideo(video, settings)
	if templateSlug == models.WhatsappTplWeeklySummary {
		loc, _ := time.LoadLocation(settings.Timezone)
		if loc == nil {
			loc = time.UTC
		}
		now := time.Now().In(loc)
		weekday := int(now.Weekday())
		start := now.AddDate(0, 0, -weekday)
		end := start.AddDate(0, 0, 7)
		videos, _ := s.repo.ListVideosWithProblemSentBetween(ctx, start, end)
		vars.ProblemList = templaterender.FormatProblemList(videos)
	}
	return templaterender.Render(tpl.Body, vars), nil
}

type WhatsappStatusPatch struct {
	ProblemSent     bool
	DiscussionSent  bool
	SolutionSent    bool
}

func (s *Service) PatchWhatsappStatus(ctx context.Context, videoID uuid.UUID, patch WhatsappStatusPatch) error {
	video, err := s.repo.GetVideo(ctx, videoID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if patch.ProblemSent && video.WhatsappProblemSentAt == nil {
		video.WhatsappProblemSentAt = &now
	}
	if patch.DiscussionSent && video.WhatsappDiscussionSentAt == nil {
		video.WhatsappDiscussionSentAt = &now
	}
	if patch.SolutionSent && video.WhatsappSolutionSentAt == nil {
		video.WhatsappSolutionSentAt = &now
	}
	return s.repo.UpdateVideo(ctx, video)
}

func (s *Service) SendWhatsappNow(ctx context.Context, videoID uuid.UUID, dispatchType, externalID, destinationID string) (*DispatchResult, error) {
	if strings.TrimSpace(externalID) == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "externalId is required.")
	}
	settings, err := s.repo.EnsureSettings(ctx)
	if err != nil {
		return nil, err
	}
	video, err := s.repo.GetVideo(ctx, videoID)
	if err != nil {
		return nil, err
	}
	day := templaterender.TodayInTZ(settings.Timezone)
	if video.ProblemDate != nil {
		day = *video.ProblemDate
	}
	var res *DispatchResult
	switch dispatchType {
	case "problem":
		res, err = s.dispatchProblem(ctx, settings, day)
	case "discussion":
		res, err = s.dispatchDiscussion(ctx, settings, day)
	case "solution":
		res, err = s.dispatchSolution(ctx, settings, day)
	default:
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Invalid dispatch type.")
	}
	if err != nil {
		return nil, err
	}
	if res.Skip {
		return res, nil
	}
	if s.whatsappWorker == nil || !s.whatsappWorker.Enabled() {
		return nil, apperrors.Invalid(apperrors.CodeWhatsappWorkerUnavailable, apperrors.MsgWhatsappWorkerUnavailable)
	}
	if err := s.whatsappWorker.Send(ctx, whatsappworkerclient.SendRequest{
		Message:       res.Message,
		ExternalID:    externalID,
		DestinationID: destinationID,
		Type:          dispatchType,
		VideoID:       res.VideoID,
		TemplateSlug:  res.TemplateSlug,
	}); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "send whatsapp", err)
	}
	if res.VideoID != "" {
		patch := WhatsappStatusPatch{}
		switch dispatchType {
		case "problem":
			patch.ProblemSent = true
		case "discussion":
			patch.DiscussionSent = true
		case "solution":
			patch.SolutionSent = true
		}
		if err := s.PatchWhatsappStatus(ctx, video.ID, patch); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (s *Service) WhatsappWorkerStatus(ctx context.Context) (map[string]any, error) {
	if s.whatsappWorker == nil || !s.whatsappWorker.Enabled() {
		return map[string]any{"configured": false, "connected": false}, nil
	}
	st, err := s.whatsappWorker.Status(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "whatsapp worker status", err)
	}
	return map[string]any{"configured": true, "connected": st.Connected}, nil
}

func (s *Service) WhatsappWorkerQR(ctx context.Context) (map[string]any, error) {
	if s.whatsappWorker == nil || !s.whatsappWorker.Enabled() {
		return nil, apperrors.Invalid(apperrors.CodeWhatsappWorkerUnavailable, apperrors.MsgWhatsappWorkerUnavailable)
	}
	qr, err := s.whatsappWorker.QR(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "whatsapp worker qr", err)
	}
	if qr == nil || qr.QR == "" {
		return map[string]any{"qr": nil}, nil
	}
	return map[string]any{"qr": qr.QR}, nil
}
