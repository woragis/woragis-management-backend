package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/contacts/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	validRelationships = map[string]bool{
		"lead": true, "prospect": true, "client": true,
		"investor": true, "partner": true, "other": true,
	}
	validStages = map[string]bool{
		"cold": true, "warm": true, "active": true, "paused": true, "churned": true,
	}
	validInteractionTypes = map[string]bool{
		"call": true, "meeting": true, "message": true, "email": true, "note": true,
	}
	validChannels = map[string]bool{
		"telegram": true, "whatsapp": true, "phone": true, "in_person": true, "other": true,
	}
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

type ListFilter struct {
	Query        string
	Relationship string
	Organization string
	Stage        string
	ProjectID    *uuid.UUID
	ActiveOnly   bool
}

type CreateContactInput struct {
	Name           string
	DisplayName    string
	Email          string
	Phone          string
	Telegram       string
	Whatsapp       string
	Organization   string
	RoleTitle      string
	Relationship   string
	Stage          string
	Source         string
	Notes          string
	Tags           []string
	ProjectID      *uuid.UUID
	NextFollowUpAt *time.Time
	Active         bool
}

type UpdateContactInput struct {
	Name           *string
	DisplayName    *string
	Email          *string
	Phone          *string
	Telegram       *string
	Whatsapp       *string
	Organization   *string
	RoleTitle      *string
	Relationship   *string
	Stage          *string
	Source         *string
	Notes          *string
	Tags           []string
	TagsSet        bool
	ProjectID      *uuid.UUID
	ProjectSet     bool
	NextFollowUpAt *time.Time
	NextFollowUpSet bool
	Active         *bool
}

type CreateInteractionInput struct {
	Type       string
	Channel    string
	Summary    string
	HappenedAt time.Time
}

func (s *Service) List(ctx context.Context, f ListFilter) ([]models.Contact, error) {
	rows, err := s.repo.ListContacts(ctx, repository.ListFilter(f))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load contacts.", err)
	}
	return rows, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Contact, error) {
	row, err := s.repo.FindContact(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Contact not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load contact.", err)
	}
	return row, nil
}

func (s *Service) Create(ctx context.Context, in CreateContactInput) (*models.Contact, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
	}
	displayName := strings.TrimSpace(in.DisplayName)
	if displayName == "" {
		displayName = buildDisplayName(name, in.RoleTitle, in.Organization)
	}
	row := &models.Contact{
		Name:           name,
		DisplayName:    displayName,
		Email:          strings.TrimSpace(in.Email),
		Phone:          strings.TrimSpace(in.Phone),
		Telegram:       strings.TrimSpace(in.Telegram),
		Whatsapp:       strings.TrimSpace(in.Whatsapp),
		Organization:   strings.TrimSpace(in.Organization),
		RoleTitle:      strings.TrimSpace(in.RoleTitle),
		Relationship:   normalizeRelationship(in.Relationship),
		Stage:          normalizeStage(in.Stage),
		Source:         strings.TrimSpace(in.Source),
		Notes:          strings.TrimSpace(in.Notes),
		Tags:           tagsJSON(in.Tags),
		ProjectID:      in.ProjectID,
		NextFollowUpAt: in.NextFollowUpAt,
		Active:         in.Active,
	}
	if err := s.repo.CreateContact(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create contact.", err)
	}
	return row, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateContactInput) (*models.Contact, error) {
	row, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		name := strings.TrimSpace(*in.Name)
		if name == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
		}
		row.Name = name
	}
	if in.DisplayName != nil {
		row.DisplayName = strings.TrimSpace(*in.DisplayName)
	}
	if in.Email != nil {
		row.Email = strings.TrimSpace(*in.Email)
	}
	if in.Phone != nil {
		row.Phone = strings.TrimSpace(*in.Phone)
	}
	if in.Telegram != nil {
		row.Telegram = strings.TrimSpace(*in.Telegram)
	}
	if in.Whatsapp != nil {
		row.Whatsapp = strings.TrimSpace(*in.Whatsapp)
	}
	if in.Organization != nil {
		row.Organization = strings.TrimSpace(*in.Organization)
	}
	if in.RoleTitle != nil {
		row.RoleTitle = strings.TrimSpace(*in.RoleTitle)
	}
	if in.Relationship != nil {
		row.Relationship = normalizeRelationship(*in.Relationship)
	}
	if in.Stage != nil {
		row.Stage = normalizeStage(*in.Stage)
	}
	if in.Source != nil {
		row.Source = strings.TrimSpace(*in.Source)
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if in.TagsSet {
		row.Tags = tagsJSON(in.Tags)
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.NextFollowUpSet {
		row.NextFollowUpAt = in.NextFollowUpAt
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if row.DisplayName == "" {
		row.DisplayName = buildDisplayName(row.Name, row.RoleTitle, row.Organization)
	}
	if err := s.repo.SaveContact(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update contact.", err)
	}
	return row, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	row, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	row.Active = false
	if err := s.repo.SaveContact(ctx, row); err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to deactivate contact.", err)
	}
	return nil
}

func (s *Service) ListInteractions(ctx context.Context, contactID uuid.UUID) ([]models.ContactInteraction, error) {
	if _, err := s.GetByID(ctx, contactID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListInteractions(ctx, contactID)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load interactions.", err)
	}
	return rows, nil
}

func (s *Service) CreateInteraction(ctx context.Context, contactID uuid.UUID, in CreateInteractionInput) (*models.ContactInteraction, error) {
	contact, err := s.GetByID(ctx, contactID)
	if err != nil {
		return nil, err
	}
	txType := normalizeInteractionType(in.Type)
	if txType == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Interaction type is invalid.")
	}
	happenedAt := in.HappenedAt
	if happenedAt.IsZero() {
		happenedAt = time.Now().UTC()
	}
	row := &models.ContactInteraction{
		ContactID:  contactID,
		Type:       txType,
		Channel:    normalizeChannel(in.Channel),
		Summary:    strings.TrimSpace(in.Summary),
		HappenedAt: happenedAt.UTC(),
	}
	if err := s.repo.CreateInteraction(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create interaction.", err)
	}
	contact.LastContactedAt = &row.HappenedAt
	if err := s.repo.SaveContact(ctx, contact); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update last contacted.", err)
	}
	return row, nil
}

func (s *Service) ListDueFollowUp(ctx context.Context, before time.Time) ([]models.Contact, error) {
	if before.IsZero() {
		before = time.Now().UTC()
	}
	rows, err := s.repo.ListContactsDueFollowUp(ctx, before)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load follow-up contacts.", err)
	}
	return rows, nil
}

func buildDisplayName(name, role, org string) string {
	name = strings.TrimSpace(name)
	role = strings.TrimSpace(role)
	org = strings.TrimSpace(org)
	switch {
	case role != "" && org != "":
		return name + " — " + role + ", " + org
	case role != "":
		return name + " — " + role
	case org != "":
		return name + " — " + org
	default:
		return name
	}
}

func normalizeRelationship(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if validRelationships[v] {
		return v
	}
	return "other"
}

func normalizeStage(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if validStages[v] {
		return v
	}
	return "cold"
}

func normalizeInteractionType(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if validInteractionTypes[v] {
		return v
	}
	return ""
}

func normalizeChannel(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if validChannels[v] {
		return v
	}
	return "other"
}

func tagsJSON(tags []string) datatypes.JSON {
	if len(tags) == 0 {
		return datatypes.JSON([]byte("[]"))
	}
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	b, _ := json.Marshal(out)
	return datatypes.JSON(b)
}

func ParseTags(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return []string{}
	}
	return out
}
