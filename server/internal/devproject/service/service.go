package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	secretcrypto "github.com/woragis/management/backend/server/internal/crypto"
	"github.com/woragis/management/backend/server/internal/devproject/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Service struct {
	repo       *repository.Repository
	secretsKey []byte
}

func New(repo *repository.Repository, secretsKey []byte) *Service {
	return &Service{repo: repo, secretsKey: secretsKey}
}

type CreateProjectInput struct {
	Name             string
	Slug             string
	Description      string
	ShortDescription string
	LongDescription  string
	Status           string
	Intent           string
	Distribution     []string
	Monetization     string
	Maturity         string
	VisibilityGoal   string
	Stack            []string
	RepoURL          string
	DemoURL          string
	GithubURL        string
	RepoVisibility   string
	Notes            string
	IsPublic         bool
	Featured         bool
	DisplayOrder     int
	PublicSlug       string
	CoverImageID     *uuid.UUID
	ParentProjectID  *uuid.UUID
}

type UpdateProjectInput struct {
	Name             *string
	Slug             *string
	Description      *string
	ShortDescription *string
	LongDescription  *string
	Status           *string
	Intent           *string
	Distribution     []string
	DistributionSet  bool
	Monetization     *string
	Maturity         *string
	VisibilityGoal   *string
	Stack            []string
	StackSet         bool
	RepoURL          *string
	DemoURL          *string
	GithubURL        *string
	RepoVisibility   *string
	Notes            *string
	IsPublic         *bool
	Featured         *bool
	DisplayOrder     *int
	PublicSlug       *string
	CoverImageID     *uuid.UUID
	CoverImageSet    bool
	ParentProjectID  *uuid.UUID
	ParentProjectSet bool
}

type CreateLinkInput struct {
	Type        string
	URL         string
	Environment string
	Label       string
	IsPublic    bool
}

type CreateDomainInput struct {
	Domain    string
	Registrar string
	Purpose   string
	ExpiresAt *time.Time
	Notes     string
}

type CreateSecretInput struct {
	Name        string
	Value       string
	Environment string
	Service     string
	ExpiresAt   *time.Time
	Notes       string
}

type CreateGalleryInput struct {
	MediaAssetID uuid.UUID
	DisplayOrder int
	Caption      string
}

type PublicProject struct {
	ID               uuid.UUID               `json:"id"`
	Name             string                  `json:"name"`
	PublicSlug       string                  `json:"publicSlug"`
	ShortDescription string                  `json:"shortDescription"`
	LongDescription  string                  `json:"longDescription"`
	Status           string                  `json:"status"`
	Stack            []string                `json:"stack"`
	DemoURL          string                  `json:"demoUrl"`
	GithubURL        string                  `json:"githubUrl,omitempty"`
	RepoURL          string                  `json:"repoUrl,omitempty"`
	Featured         bool                    `json:"featured"`
	DisplayOrder     int                     `json:"displayOrder"`
	CoverImageID     *uuid.UUID              `json:"coverImageId"`
	Links            []models.ProjectLink    `json:"links,omitempty"`
	Gallery          []models.ProjectGallery `json:"gallery,omitempty"`
}

func (s *Service) List(ctx context.Context) ([]models.Project, error) {
	out, err := s.repo.ListProjects(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	p, err := s.repo.FindProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeProjectGetV1ServiceNotFound, apperrors.MsgProjectGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeProjectGetV1ServiceLoadFailed, apperrors.MsgProjectGetV1ServiceLoadFailed, err)
	}
	return p, nil
}

func (s *Service) Create(ctx context.Context, in CreateProjectInput) (*models.Project, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceNameEmpty, apperrors.MsgProjectPostV1ServiceNameEmpty)
	}
	slug := normalizeSlug(in.Slug, name)
	if slug == "" || !slugPattern.MatchString(slug) {
		return nil, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceSlugInvalid, apperrors.MsgProjectPostV1ServiceSlugInvalid)
	}
	p := &models.Project{
		Name:             name,
		Slug:             slug,
		Description:      strings.TrimSpace(in.Description),
		ShortDescription: strings.TrimSpace(in.ShortDescription),
		LongDescription:  strings.TrimSpace(in.LongDescription),
		Status:           normalizeStatus(in.Status),
		Intent:           normalizeIntent(in.Intent),
		Distribution:     distributionJSON(in.Distribution),
		Monetization:     normalizeMonetization(in.Monetization),
		Maturity:         normalizeMaturity(in.Maturity),
		VisibilityGoal:   normalizeVisibilityGoal(in.VisibilityGoal),
		Stack:            stackJSON(in.Stack),
		RepoURL:          strings.TrimSpace(in.RepoURL),
		DemoURL:          strings.TrimSpace(in.DemoURL),
		GithubURL:        strings.TrimSpace(in.GithubURL),
		RepoVisibility:   normalizeRepoVisibility(in.RepoVisibility),
		Notes:            strings.TrimSpace(in.Notes),
		IsPublic:         in.IsPublic,
		Featured:         in.Featured,
		DisplayOrder:     in.DisplayOrder,
		PublicSlug:       strings.TrimSpace(in.PublicSlug),
		CoverImageID:     in.CoverImageID,
		ParentProjectID:  in.ParentProjectID,
	}
	if err := s.repo.CreateProject(ctx, p); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectPostV1ServiceCreateFailed, apperrors.MsgProjectPostV1ServiceCreateFailed, err)
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateProjectInput) (*models.Project, error) {
	p, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceNameEmpty, apperrors.MsgProjectPostV1ServiceNameEmpty)
		}
		p.Name = name
	}
	if in.Slug != nil {
		slug := normalizeSlug(*in.Slug, p.Name)
		if slug == "" || !slugPattern.MatchString(slug) {
			return nil, apperrors.Invalid(apperrors.CodeProjectPostV1ServiceSlugInvalid, apperrors.MsgProjectPostV1ServiceSlugInvalid)
		}
		p.Slug = slug
	}
	if in.Description != nil {
		p.Description = strings.TrimSpace(*in.Description)
	}
	if in.ShortDescription != nil {
		p.ShortDescription = strings.TrimSpace(*in.ShortDescription)
	}
	if in.LongDescription != nil {
		p.LongDescription = strings.TrimSpace(*in.LongDescription)
	}
	if in.Status != nil {
		p.Status = normalizeStatus(*in.Status)
	}
	if in.Intent != nil {
		p.Intent = normalizeIntent(*in.Intent)
	}
	if in.DistributionSet {
		p.Distribution = distributionJSON(in.Distribution)
	}
	if in.Monetization != nil {
		p.Monetization = normalizeMonetization(*in.Monetization)
	}
	if in.Maturity != nil {
		p.Maturity = normalizeMaturity(*in.Maturity)
	}
	if in.VisibilityGoal != nil {
		p.VisibilityGoal = normalizeVisibilityGoal(*in.VisibilityGoal)
	}
	if in.StackSet {
		p.Stack = stackJSON(in.Stack)
	}
	if in.RepoURL != nil {
		p.RepoURL = strings.TrimSpace(*in.RepoURL)
	}
	if in.DemoURL != nil {
		p.DemoURL = strings.TrimSpace(*in.DemoURL)
	}
	if in.GithubURL != nil {
		p.GithubURL = strings.TrimSpace(*in.GithubURL)
	}
	if in.RepoVisibility != nil {
		p.RepoVisibility = normalizeRepoVisibility(*in.RepoVisibility)
	}
	if in.Notes != nil {
		p.Notes = strings.TrimSpace(*in.Notes)
	}
	if in.IsPublic != nil {
		p.IsPublic = *in.IsPublic
	}
	if in.Featured != nil {
		p.Featured = *in.Featured
	}
	if in.DisplayOrder != nil {
		p.DisplayOrder = *in.DisplayOrder
	}
	if in.PublicSlug != nil {
		p.PublicSlug = strings.TrimSpace(*in.PublicSlug)
	}
	if in.CoverImageSet {
		p.CoverImageID = in.CoverImageID
	}
	if in.ParentProjectSet {
		p.ParentProjectID = in.ParentProjectID
	}
	if err := s.repo.SaveProject(ctx, p); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectPatchV1ServiceUpdateFailed, apperrors.MsgProjectPatchV1ServiceUpdateFailed, err)
	}
	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteProject(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectDeleteV1ServiceNotFound, apperrors.MsgProjectDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectDeleteV1ServiceDeleteFailed, apperrors.MsgProjectDeleteV1ServiceDeleteFailed, err)
	}
	return nil
}

func (s *Service) CreateLink(ctx context.Context, projectID uuid.UUID, in CreateLinkInput) (*models.ProjectLink, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	url := strings.TrimSpace(in.URL)
	if url == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectLinkPostV1ServiceURLInvalid, apperrors.MsgProjectLinkPostV1ServiceURLInvalid)
	}
	link := &models.ProjectLink{
		ProjectID:   projectID,
		Type:        strings.TrimSpace(in.Type),
		URL:         url,
		Environment: normalizeEnvironment(in.Environment),
		Label:       strings.TrimSpace(in.Label),
		IsPublic:    in.IsPublic,
	}
	if link.Type == "" {
		link.Type = "other"
	}
	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return link, nil
}

func (s *Service) DeleteLink(ctx context.Context, projectID, linkID uuid.UUID) error {
	if err := s.repo.DeleteLink(ctx, projectID, linkID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectLinkDeleteV1ServiceNotFound, apperrors.MsgProjectLinkDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return nil
}

func (s *Service) CreateDomain(ctx context.Context, projectID uuid.UUID, in CreateDomainInput) (*models.ProjectDomain, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	domain := strings.TrimSpace(in.Domain)
	if domain == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectDomainPostV1ServiceDomainInvalid, apperrors.MsgProjectDomainPostV1ServiceDomainInvalid)
	}
	d := &models.ProjectDomain{
		ProjectID: projectID,
		Domain:    domain,
		Registrar: strings.TrimSpace(in.Registrar),
		Purpose:   strings.TrimSpace(in.Purpose),
		ExpiresAt: in.ExpiresAt,
		Notes:     strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateDomain(ctx, d); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return d, nil
}

func (s *Service) DeleteDomain(ctx context.Context, projectID, domainID uuid.UUID) error {
	if err := s.repo.DeleteDomain(ctx, projectID, domainID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectDomainDeleteV1ServiceNotFound, apperrors.MsgProjectDomainDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return nil
}

func (s *Service) ListSecrets(ctx context.Context, projectID uuid.UUID) ([]models.ProjectSecretView, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListSecrets(ctx, projectID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	out := make([]models.ProjectSecretView, 0, len(rows))
	for _, row := range rows {
		out = append(out, secretView(row, ""))
	}
	return out, nil
}

func (s *Service) GetSecret(ctx context.Context, projectID, secretID uuid.UUID) (*models.ProjectSecretView, error) {
	row, err := s.repo.FindSecret(ctx, projectID, secretID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeProjectSecretGetV1ServiceNotFound, apperrors.MsgProjectSecretGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeProjectGetV1ServiceLoadFailed, apperrors.MsgProjectGetV1ServiceLoadFailed, err)
	}
	plain, err := secretcrypto.Decrypt(s.secretsKey, row.EncryptedValue)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectSecretGetV1ServiceDecryptFailed, apperrors.MsgProjectSecretGetV1ServiceDecryptFailed, err)
	}
	v := secretView(*row, plain)
	return &v, nil
}

func (s *Service) CreateSecret(ctx context.Context, projectID uuid.UUID, in CreateSecretInput) (*models.ProjectSecretView, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectSecretPostV1ServiceNameEmpty, apperrors.MsgProjectSecretPostV1ServiceNameEmpty)
	}
	value := in.Value
	if strings.TrimSpace(value) == "" {
		return nil, apperrors.Invalid(apperrors.CodeProjectSecretPostV1ServiceValueEmpty, apperrors.MsgProjectSecretPostV1ServiceValueEmpty)
	}
	enc, err := secretcrypto.Encrypt(s.secretsKey, value)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectSecretPostV1ServiceEncryptFailed, apperrors.MsgProjectSecretPostV1ServiceEncryptFailed, err)
	}
	row := &models.ProjectSecret{
		ProjectID:      projectID,
		Name:           name,
		EncryptedValue: enc,
		Environment:    normalizeEnvironment(in.Environment),
		Service:        strings.TrimSpace(in.Service),
		ExpiresAt:      in.ExpiresAt,
		Notes:          strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateSecret(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	v := secretView(*row, "")
	return &v, nil
}

func (s *Service) DeleteSecret(ctx context.Context, projectID, secretID uuid.UUID) error {
	if err := s.repo.DeleteSecret(ctx, projectID, secretID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectSecretGetV1ServiceNotFound, apperrors.MsgProjectSecretGetV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return nil
}

func (s *Service) CreateGalleryItem(ctx context.Context, projectID uuid.UUID, in CreateGalleryInput) (*models.ProjectGallery, error) {
	if _, err := s.GetByID(ctx, projectID); err != nil {
		return nil, err
	}
	if in.MediaAssetID == uuid.Nil {
		return nil, apperrors.Invalid(apperrors.CodeProjectGalleryPostV1ServiceMediaInvalid, apperrors.MsgProjectGalleryPostV1ServiceMediaInvalid)
	}
	item := &models.ProjectGallery{
		ProjectID:    projectID,
		MediaAssetID: in.MediaAssetID,
		DisplayOrder: in.DisplayOrder,
		Caption:      strings.TrimSpace(in.Caption),
	}
	if err := s.repo.CreateGalleryItem(ctx, item); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return item, nil
}

func (s *Service) DeleteGalleryItem(ctx context.Context, projectID, itemID uuid.UUID) error {
	if err := s.repo.DeleteGalleryItem(ctx, projectID, itemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeProjectGalleryDeleteV1ServiceNotFound, apperrors.MsgProjectGalleryDeleteV1ServiceNotFound)
		}
		return apperrors.InternalCause(apperrors.CodeProjectLinkPostV1ServiceCreateFailed, apperrors.MsgProjectLinkPostV1ServiceCreateFailed, err)
	}
	return nil
}

func (s *Service) ListPublic(ctx context.Context, featuredOnly bool) ([]PublicProject, error) {
	rows, err := s.repo.ListPublicProjects(ctx, featuredOnly)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeProjectListV1ServiceLoadFailed, apperrors.MsgProjectListV1ServiceLoadFailed, err)
	}
	out := make([]PublicProject, 0, len(rows))
	for _, p := range rows {
		out = append(out, toPublicProject(p))
	}
	return out, nil
}

func (s *Service) GetPublicBySlug(ctx context.Context, slug string) (*PublicProject, error) {
	p, err := s.repo.FindProjectByPublicSlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeProjectGetV1ServiceNotFound, apperrors.MsgProjectGetV1ServiceNotFound)
		}
		return nil, apperrors.InternalCause(apperrors.CodeProjectGetV1ServiceLoadFailed, apperrors.MsgProjectGetV1ServiceLoadFailed, err)
	}
	pub := toPublicProject(*p)
	return &pub, nil
}

func toPublicProject(p models.Project) PublicProject {
	publicLinks := make([]models.ProjectLink, 0)
	for _, l := range p.Links {
		if l.IsPublic {
			publicLinks = append(publicLinks, l)
		}
	}
	pub := PublicProject{
		ID:               p.ID,
		Name:             p.Name,
		PublicSlug:       p.PublicSlug,
		ShortDescription: p.ShortDescription,
		LongDescription:  p.LongDescription,
		Status:           p.Status,
		Stack:            parseStack(p.Stack),
		DemoURL:          p.DemoURL,
		Featured:         p.Featured,
		DisplayOrder:     p.DisplayOrder,
		CoverImageID:     p.CoverImageID,
		Links:            publicLinks,
		Gallery:          p.Gallery,
	}
	if normalizeRepoVisibility(p.RepoVisibility) == "public" {
		pub.GithubURL = p.GithubURL
		if pub.GithubURL == "" {
			pub.GithubURL = p.RepoURL
		}
		pub.RepoURL = p.RepoURL
	}
	return pub
}

func secretView(row models.ProjectSecret, value string) models.ProjectSecretView {
	return models.ProjectSecretView{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		Name:        row.Name,
		Value:       value,
		Environment: row.Environment,
		Service:     row.Service,
		ExpiresAt:   row.ExpiresAt,
		Notes:       row.Notes,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func normalizeSlug(slug, name string) string {
	s := strings.TrimSpace(strings.ToLower(slug))
	if s == "" {
		s = strings.TrimSpace(strings.ToLower(name))
	}
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func normalizeStatus(status string) string {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case "paused", "archived":
		return strings.TrimSpace(strings.ToLower(status))
	default:
		return "active"
	}
}

func normalizeEnvironment(env string) string {
	switch strings.TrimSpace(strings.ToLower(env)) {
	case "staging", "local":
		return strings.TrimSpace(strings.ToLower(env))
	default:
		return "production"
	}
}

func normalizeRepoVisibility(v string) string {
	if strings.TrimSpace(strings.ToLower(v)) == "public" {
		return "public"
	}
	return "private"
}

func stackJSON(stack []string) datatypes.JSON {
	if stack == nil {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(stack)
	return datatypes.JSON(b)
}

func parseStack(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return []string{}
	}
	return out
}
