package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	"github.com/woragis/management/backend/server/internal/finance/repository"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	repo     *repository.Repository
	contacts ContactValidator
}

type ContactValidator interface {
	ValidateActiveContact(ctx context.Context, id uuid.UUID) error
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SetContactValidator(v ContactValidator) {
	s.contacts = v
}

func (s *Service) validateContactID(ctx context.Context, id *uuid.UUID) error {
	if id == nil || *id == uuid.Nil {
		return nil
	}
	if s.contacts == nil {
		return nil
	}
	return s.contacts.ValidateActiveContact(ctx, *id)
}

type CreateIncomeSourceInput struct {
	Name        string
	Type        string
	AmountCents int64
	Currency    string
	Frequency   string
	DayOfMonth  int
	ProjectID   *uuid.UUID
	ContactID   *uuid.UUID
	Active      bool
	Notes       string
}

type UpdateIncomeSourceInput struct {
	Name        *string
	Type        *string
	AmountCents *int64
	Currency    *string
	Frequency   *string
	DayOfMonth  *int
	ProjectID   *uuid.UUID
	ProjectSet  bool
	ContactID   *uuid.UUID
	ContactSet  bool
	Active      *bool
	Notes       *string
}

type IncomeSourceFilter struct {
	ContactID *uuid.UUID
	ProjectID *uuid.UUID
}

func (s *Service) ListIncomeSources(ctx context.Context) ([]models.IncomeSource, error) {
	return s.ListIncomeSourcesFiltered(ctx, IncomeSourceFilter{})
}

func (s *Service) ListIncomeSourcesFiltered(ctx context.Context, f IncomeSourceFilter) ([]models.IncomeSource, error) {
	rows, err := s.repo.ListIncomeSources(ctx, repository.IncomeSourceFilter{
		ActiveOnly: false,
		ContactID:  f.ContactID,
		ProjectID:  f.ProjectID,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load income sources.", err)
	}
	return rows, nil
}

func (s *Service) GetIncomeSource(ctx context.Context, id uuid.UUID) (*models.IncomeSource, error) {
	row, err := s.repo.FindIncomeSource(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Income source not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load income source.", err)
	}
	return row, nil
}

func (s *Service) CreateIncomeSource(ctx context.Context, in CreateIncomeSourceInput) (*models.IncomeSource, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
	}
	if err := s.validateContactID(ctx, in.ContactID); err != nil {
		return nil, err
	}
	row := &models.IncomeSource{
		Name:        name,
		Type:        normalizeIncomeType(in.Type),
		AmountCents: in.AmountCents,
		Currency:    normalizeCurrency(in.Currency),
		Frequency:   normalizeFrequency(in.Frequency),
		DayOfMonth:  clampDay(in.DayOfMonth),
		ProjectID:   in.ProjectID,
		ContactID:   in.ContactID,
		Active:      in.Active,
		Notes:       strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateIncomeSource(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create income source.", err)
	}
	return row, nil
}

func (s *Service) UpdateIncomeSource(ctx context.Context, id uuid.UUID, in UpdateIncomeSourceInput) (*models.IncomeSource, error) {
	row, err := s.GetIncomeSource(ctx, id)
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
	if in.Type != nil {
		row.Type = normalizeIncomeType(*in.Type)
	}
	if in.AmountCents != nil {
		row.AmountCents = *in.AmountCents
	}
	if in.Currency != nil {
		row.Currency = normalizeCurrency(*in.Currency)
	}
	if in.Frequency != nil {
		row.Frequency = normalizeFrequency(*in.Frequency)
	}
	if in.DayOfMonth != nil {
		row.DayOfMonth = clampDay(*in.DayOfMonth)
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.ContactSet {
		if err := s.validateContactID(ctx, in.ContactID); err != nil {
			return nil, err
		}
		row.ContactID = in.ContactID
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SaveIncomeSource(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update income source.", err)
	}
	return row, nil
}

func (s *Service) DeleteIncomeSource(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteIncomeSource(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Income source not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete income source.", err)
	}
	return nil
}

type CreateExpenseInput struct {
	Name        string
	Category    string
	AmountCents int64
	Currency    string
	Frequency   string
	DayOfMonth  int
	DueDate     *time.Time
	AutoPay     bool
	ProjectID   *uuid.UUID
	Active      bool
	Notes       string
}

type UpdateExpenseInput struct {
	Name        *string
	Category    *string
	AmountCents *int64
	Currency    *string
	Frequency   *string
	DayOfMonth  *int
	DueDate     *time.Time
	DueDateSet  bool
	AutoPay     *bool
	ProjectID   *uuid.UUID
	ProjectSet  bool
	Active      *bool
	Notes       *string
}

func (s *Service) ListExpenses(ctx context.Context) ([]models.Expense, error) {
	rows, err := s.repo.ListExpenses(ctx, false)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load expenses.", err)
	}
	return rows, nil
}

func (s *Service) GetExpense(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	row, err := s.repo.FindExpense(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Expense not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load expense.", err)
	}
	return row, nil
}

func (s *Service) CreateExpense(ctx context.Context, in CreateExpenseInput) (*models.Expense, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
	}
	row := &models.Expense{
		Name:        name,
		Category:    normalizeExpenseCategory(in.Category),
		AmountCents: in.AmountCents,
		Currency:    normalizeCurrency(in.Currency),
		Frequency:   normalizeFrequency(in.Frequency),
		DayOfMonth:  clampDay(in.DayOfMonth),
		DueDate:     in.DueDate,
		AutoPay:     in.AutoPay,
		ProjectID:   in.ProjectID,
		Active:      in.Active,
		Notes:       strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateExpense(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create expense.", err)
	}
	return row, nil
}

func (s *Service) UpdateExpense(ctx context.Context, id uuid.UUID, in UpdateExpenseInput) (*models.Expense, error) {
	row, err := s.GetExpense(ctx, id)
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
	if in.Category != nil {
		row.Category = normalizeExpenseCategory(*in.Category)
	}
	if in.AmountCents != nil {
		row.AmountCents = *in.AmountCents
	}
	if in.Currency != nil {
		row.Currency = normalizeCurrency(*in.Currency)
	}
	if in.Frequency != nil {
		row.Frequency = normalizeFrequency(*in.Frequency)
	}
	if in.DayOfMonth != nil {
		row.DayOfMonth = clampDay(*in.DayOfMonth)
	}
	if in.DueDateSet {
		row.DueDate = in.DueDate
	}
	if in.AutoPay != nil {
		row.AutoPay = *in.AutoPay
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.Active != nil {
		row.Active = *in.Active
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SaveExpense(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update expense.", err)
	}
	return row, nil
}

func (s *Service) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteExpense(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Expense not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete expense.", err)
	}
	return nil
}

func normalizeIncomeType(t string) string {
	switch strings.TrimSpace(strings.ToLower(t)) {
	case "salary", "freelance", "saas", "business":
		return strings.TrimSpace(strings.ToLower(t))
	default:
		return "other"
	}
}

func normalizeExpenseCategory(c string) string {
	switch strings.TrimSpace(strings.ToLower(c)) {
	case "subscription", "utilities", "rent", "food", "investment", "transport", "health":
		return strings.TrimSpace(strings.ToLower(c))
	default:
		return "other"
	}
}

func normalizeFrequency(f string) string {
	switch strings.TrimSpace(strings.ToLower(f)) {
	case "weekly", "yearly", "one_time":
		return strings.TrimSpace(strings.ToLower(f))
	default:
		return "monthly"
	}
}

func normalizeCurrency(c string) string {
	c = strings.TrimSpace(strings.ToUpper(c))
	if c == "" {
		return "BRL"
	}
	return c
}

func clampDay(d int) int {
	if d < 1 {
		return 1
	}
	if d > 31 {
		return 31
	}
	return d
}

func parseYearMonth(year, month int) (int, int) {
	now := time.Now().UTC()
	if year <= 0 {
		year = now.Year()
	}
	if month <= 0 || month > 12 {
		month = int(now.Month())
	}
	return year, month
}
