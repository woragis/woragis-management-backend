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

type CreateTransactionInput struct {
	Type           string
	AmountCents    int64
	Currency       string
	Description    string
	Date           time.Time
	IncomeSourceID *uuid.UUID
	ExpenseID      *uuid.UUID
	ProjectID      *uuid.UUID
	InvoiceID      *uuid.UUID
	Notes          string
}

type UpdateTransactionInput struct {
	Type           *string
	AmountCents    *int64
	Currency       *string
	Description    *string
	Date           *time.Time
	IncomeSourceID *uuid.UUID
	IncomeSourceSet bool
	ExpenseID      *uuid.UUID
	ExpenseSet     bool
	ProjectID      *uuid.UUID
	ProjectSet     bool
	InvoiceID      *uuid.UUID
	InvoiceSet     bool
	Notes          *string
}

type TransactionFilter struct {
	Type      string
	Year      int
	Month     int
	ProjectID *uuid.UUID
}

func (s *Service) ListTransactions(ctx context.Context, f TransactionFilter) ([]models.Transaction, error) {
	year, month := parseYearMonth(f.Year, f.Month)
	rows, err := s.repo.ListTransactions(ctx, repository.TransactionFilter{
		Type:      f.Type,
		Year:      year,
		Month:     month,
		ProjectID: f.ProjectID,
	})
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load transactions.", err)
	}
	return rows, nil
}

func (s *Service) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	row, err := s.repo.FindTransaction(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Transaction not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load transaction.", err)
	}
	return row, nil
}

func (s *Service) CreateTransaction(ctx context.Context, in CreateTransactionInput) (*models.Transaction, error) {
	txType := normalizeTransactionType(in.Type)
	if txType == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Type must be income or expense.")
	}
	desc := strings.TrimSpace(in.Description)
	if desc == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Description is required.")
	}
	if in.AmountCents <= 0 {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Amount must be positive.")
	}
	date := in.Date
	if date.IsZero() {
		date = time.Now().UTC()
	}
	row := &models.Transaction{
		Type:           txType,
		AmountCents:    in.AmountCents,
		Currency:       normalizeCurrency(in.Currency),
		Description:    desc,
		Date:           date.UTC().Truncate(24 * time.Hour),
		IncomeSourceID: in.IncomeSourceID,
		ExpenseID:      in.ExpenseID,
		ProjectID:      in.ProjectID,
		InvoiceID:      in.InvoiceID,
		Notes:          strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateTransaction(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create transaction.", err)
	}
	return row, nil
}

func (s *Service) UpdateTransaction(ctx context.Context, id uuid.UUID, in UpdateTransactionInput) (*models.Transaction, error) {
	row, err := s.GetTransaction(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Type != nil {
		t := normalizeTransactionType(*in.Type)
		if t == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Type must be income or expense.")
		}
		row.Type = t
	}
	if in.AmountCents != nil {
		if *in.AmountCents <= 0 {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Amount must be positive.")
		}
		row.AmountCents = *in.AmountCents
	}
	if in.Currency != nil {
		row.Currency = normalizeCurrency(*in.Currency)
	}
	if in.Description != nil {
		desc := strings.TrimSpace(*in.Description)
		if desc == "" {
			return nil, apperrors.Invalid(apperrors.CodeInternal, "Description is required.")
		}
		row.Description = desc
	}
	if in.Date != nil {
		row.Date = in.Date.UTC().Truncate(24 * time.Hour)
	}
	if in.IncomeSourceSet {
		row.IncomeSourceID = in.IncomeSourceID
	}
	if in.ExpenseSet {
		row.ExpenseID = in.ExpenseID
	}
	if in.ProjectSet {
		row.ProjectID = in.ProjectID
	}
	if in.InvoiceSet {
		row.InvoiceID = in.InvoiceID
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SaveTransaction(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update transaction.", err)
	}
	return row, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteTransaction(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Transaction not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete transaction.", err)
	}
	return nil
}

type MonthlySummary struct {
	Year          int            `json:"year"`
	Month         int            `json:"month"`
	IncomeCents   int64          `json:"incomeCents"`
	ExpenseCents  int64          `json:"expenseCents"`
	NetCents      int64          `json:"netCents"`
	ByCategory    map[string]int64 `json:"byCategory"`
	Budgets       []models.BudgetPlan `json:"budgets,omitempty"`
}

func (s *Service) MonthlySummary(ctx context.Context, year, month int) (*MonthlySummary, error) {
	year, month = parseYearMonth(year, month)
	totals, err := s.repo.SumTransactionsByMonth(ctx, year, month)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to compute summary.", err)
	}
	byCat, err := s.repo.SumExpensesByCategory(ctx, year, month)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to compute category totals.", err)
	}
	budgets, err := s.repo.ListBudgets(ctx, year, month)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load budgets.", err)
	}
	return &MonthlySummary{
		Year:         year,
		Month:        month,
		IncomeCents:  totals.IncomeCents,
		ExpenseCents: totals.ExpenseCents,
		NetCents:     totals.IncomeCents - totals.ExpenseCents,
		ByCategory:   byCat,
		Budgets:      budgets,
	}, nil
}

func normalizeTransactionType(t string) string {
	switch strings.TrimSpace(strings.ToLower(t)) {
	case "income", "expense":
		return strings.TrimSpace(strings.ToLower(t))
	default:
		return ""
	}
}

type CreateInvoiceInput struct {
	Name         string
	CardLastFour string
	DueDate      time.Time
	ClosedAt     *time.Time
	TotalCents   int64
	PaidCents    int64
	Status       string
	Notes        string
}

type UpdateInvoiceInput struct {
	Name         *string
	CardLastFour *string
	DueDate      *time.Time
	ClosedAt     *time.Time
	ClosedAtSet  bool
	TotalCents   *int64
	PaidCents    *int64
	Status       *string
	Notes        *string
}

func (s *Service) ListInvoices(ctx context.Context, status string) ([]models.Invoice, error) {
	rows, err := s.repo.ListInvoices(ctx, status)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load invoices.", err)
	}
	return rows, nil
}

func (s *Service) GetInvoice(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	row, err := s.repo.FindInvoice(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Invoice not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load invoice.", err)
	}
	return row, nil
}

func (s *Service) CreateInvoice(ctx context.Context, in CreateInvoiceInput) (*models.Invoice, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Name is required.")
	}
	if in.DueDate.IsZero() {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Due date is required.")
	}
	row := &models.Invoice{
		Name:         name,
		CardLastFour: strings.TrimSpace(in.CardLastFour),
		DueDate:      in.DueDate.UTC().Truncate(24 * time.Hour),
		ClosedAt:     in.ClosedAt,
		TotalCents:   in.TotalCents,
		PaidCents:    in.PaidCents,
		Status:       normalizeInvoiceStatus(in.Status),
		Notes:        strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateInvoice(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create invoice.", err)
	}
	return row, nil
}

func (s *Service) UpdateInvoice(ctx context.Context, id uuid.UUID, in UpdateInvoiceInput) (*models.Invoice, error) {
	row, err := s.GetInvoice(ctx, id)
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
	if in.CardLastFour != nil {
		row.CardLastFour = strings.TrimSpace(*in.CardLastFour)
	}
	if in.DueDate != nil {
		row.DueDate = in.DueDate.UTC().Truncate(24 * time.Hour)
	}
	if in.ClosedAtSet {
		row.ClosedAt = in.ClosedAt
	}
	if in.TotalCents != nil {
		row.TotalCents = *in.TotalCents
	}
	if in.PaidCents != nil {
		row.PaidCents = *in.PaidCents
	}
	if in.Status != nil {
		row.Status = normalizeInvoiceStatus(*in.Status)
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SaveInvoice(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update invoice.", err)
	}
	return s.GetInvoice(ctx, id)
}

func (s *Service) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteInvoice(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Invoice not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete invoice.", err)
	}
	return nil
}

type CreateInvoiceItemInput struct {
	Description string
	AmountCents int64
	Date        time.Time
	Category    string
	Installment string
	Notes       string
}

func (s *Service) CreateInvoiceItem(ctx context.Context, invoiceID uuid.UUID, in CreateInvoiceItemInput) (*models.InvoiceItem, error) {
	if _, err := s.GetInvoice(ctx, invoiceID); err != nil {
		return nil, err
	}
	desc := strings.TrimSpace(in.Description)
	if desc == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Description is required.")
	}
	if in.AmountCents <= 0 {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Amount must be positive.")
	}
	date := in.Date
	if date.IsZero() {
		date = time.Now().UTC()
	}
	row := &models.InvoiceItem{
		InvoiceID:   invoiceID,
		Description: desc,
		AmountCents: in.AmountCents,
		Date:        date.UTC().Truncate(24 * time.Hour),
		Category:    normalizeExpenseCategory(in.Category),
		Installment: strings.TrimSpace(in.Installment),
		Notes:       strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateInvoiceItem(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create invoice item.", err)
	}
	if err := s.repo.RecalcInvoiceTotal(ctx, invoiceID); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update invoice total.", err)
	}
	return row, nil
}

func (s *Service) DeleteInvoiceItem(ctx context.Context, invoiceID, itemID uuid.UUID) error {
	if err := s.repo.DeleteInvoiceItem(ctx, invoiceID, itemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Invoice item not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete invoice item.", err)
	}
	if err := s.repo.RecalcInvoiceTotal(ctx, invoiceID); err != nil {
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to update invoice total.", err)
	}
	return nil
}

func normalizeInvoiceStatus(st string) string {
	switch strings.TrimSpace(strings.ToLower(st)) {
	case "paid", "overdue":
		return strings.TrimSpace(strings.ToLower(st))
	default:
		return "open"
	}
}

type CreateBudgetInput struct {
	Year         int
	Month        int
	Category     string
	PlannedCents int64
	Notes        string
}

type UpdateBudgetInput struct {
	Year         *int
	Month        *int
	Category     *string
	PlannedCents *int64
	Notes        *string
}

func (s *Service) ListBudgets(ctx context.Context, year, month int) ([]models.BudgetPlan, error) {
	year, month = parseYearMonth(year, month)
	rows, err := s.repo.ListBudgets(ctx, year, month)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load budgets.", err)
	}
	return rows, nil
}

func (s *Service) GetBudget(ctx context.Context, id uuid.UUID) (*models.BudgetPlan, error) {
	row, err := s.repo.FindBudget(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound(apperrors.CodeInternal, "Budget not found.")
		}
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load budget.", err)
	}
	return row, nil
}

func (s *Service) CreateBudget(ctx context.Context, in CreateBudgetInput) (*models.BudgetPlan, error) {
	year, month := parseYearMonth(in.Year, in.Month)
	cat := normalizeExpenseCategory(in.Category)
	if cat == "" {
		return nil, apperrors.Invalid(apperrors.CodeInternal, "Category is required.")
	}
	row := &models.BudgetPlan{
		Year:         year,
		Month:        month,
		Category:     cat,
		PlannedCents: in.PlannedCents,
		Notes:        strings.TrimSpace(in.Notes),
	}
	if err := s.repo.CreateBudget(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to create budget.", err)
	}
	return row, nil
}

func (s *Service) UpdateBudget(ctx context.Context, id uuid.UUID, in UpdateBudgetInput) (*models.BudgetPlan, error) {
	row, err := s.GetBudget(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Year != nil {
		row.Year = *in.Year
	}
	if in.Month != nil {
		row.Month = *in.Month
	}
	if in.Category != nil {
		row.Category = normalizeExpenseCategory(*in.Category)
	}
	if in.PlannedCents != nil {
		row.PlannedCents = *in.PlannedCents
	}
	if in.Notes != nil {
		row.Notes = strings.TrimSpace(*in.Notes)
	}
	if err := s.repo.SaveBudget(ctx, row); err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to update budget.", err)
	}
	return row, nil
}

func (s *Service) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteBudget(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NotFound(apperrors.CodeInternal, "Budget not found.")
		}
		return apperrors.InternalCause(apperrors.CodeInternal, "Failed to delete budget.", err)
	}
	return nil
}

type CalendarEvent struct {
	Date        string `json:"date"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	AmountCents int64  `json:"amountCents"`
	RefID       string `json:"refId"`
}

type FinanceDashboard struct {
	MonthIncomeCents   int64             `json:"monthIncomeCents"`
	MonthExpenseCents  int64             `json:"monthExpenseCents"`
	MonthNetCents      int64             `json:"monthNetCents"`
	OpenInvoiceCount   int64             `json:"openInvoiceCount"`
	ActiveIncomeCount  int               `json:"activeIncomeCount"`
	ActiveExpenseCount int               `json:"activeExpenseCount"`
	UpcomingInvoices   []models.Invoice  `json:"upcomingInvoices"`
}

func (s *Service) Dashboard(ctx context.Context) (*FinanceDashboard, error) {
	now := time.Now().UTC()
	totals, err := s.repo.SumTransactionsByMonth(ctx, now.Year(), int(now.Month()))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to compute dashboard.", err)
	}
	openCount, err := s.repo.CountOpenInvoices(ctx)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to count invoices.", err)
	}
	incomes, err := s.repo.ListIncomeSources(ctx, true)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load income sources.", err)
	}
	expenses, err := s.repo.ListExpenses(ctx, true)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load expenses.", err)
	}
	upcoming, err := s.repo.ListInvoicesDueBetween(ctx, now, now.AddDate(0, 0, 30))
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load upcoming invoices.", err)
	}
	return &FinanceDashboard{
		MonthIncomeCents:   totals.IncomeCents,
		MonthExpenseCents:  totals.ExpenseCents,
		MonthNetCents:      totals.IncomeCents - totals.ExpenseCents,
		OpenInvoiceCount:   openCount,
		ActiveIncomeCount:  len(incomes),
		ActiveExpenseCount: len(expenses),
		UpcomingInvoices:   upcoming,
	}, nil
}

func (s *Service) Calendar(ctx context.Context, year, month int) ([]CalendarEvent, error) {
	year, month = parseYearMonth(year, month)
	daysInMonth := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	var events []CalendarEvent

	incomes, err := s.repo.ListIncomeSources(ctx, true)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load income sources.", err)
	}
	for _, inc := range incomes {
		if inc.Frequency != "monthly" {
			continue
		}
		day := inc.DayOfMonth
		if day > daysInMonth {
			day = daysInMonth
		}
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		events = append(events, CalendarEvent{
			Date:        d.Format("2006-01-02"),
			Type:        "income",
			Title:       inc.Name,
			AmountCents: inc.AmountCents,
			RefID:       inc.ID.String(),
		})
	}

	expenses, err := s.repo.ListExpenses(ctx, true)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load expenses.", err)
	}
	for _, exp := range expenses {
		if exp.Frequency == "one_time" && exp.DueDate != nil {
			if exp.DueDate.Year() == year && int(exp.DueDate.Month()) == month {
				events = append(events, CalendarEvent{
					Date:        exp.DueDate.Format("2006-01-02"),
					Type:        "expense",
					Title:       exp.Name,
					AmountCents: exp.AmountCents,
					RefID:       exp.ID.String(),
				})
			}
			continue
		}
		if exp.Frequency != "monthly" {
			continue
		}
		day := exp.DayOfMonth
		if day > daysInMonth {
			day = daysInMonth
		}
		d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		events = append(events, CalendarEvent{
			Date:        d.Format("2006-01-02"),
			Type:        "expense",
			Title:       exp.Name,
			AmountCents: exp.AmountCents,
			RefID:       exp.ID.String(),
		})
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	invoices, err := s.repo.ListInvoicesDueBetween(ctx, start, end)
	if err != nil {
		return nil, apperrors.InternalCause(apperrors.CodeInternal, "Failed to load invoices.", err)
	}
	for _, inv := range invoices {
		events = append(events, CalendarEvent{
			Date:        inv.DueDate.Format("2006-01-02"),
			Type:        "invoice",
			Title:       inv.Name,
			AmountCents: inv.TotalCents - inv.PaidCents,
			RefID:       inv.ID.String(),
		})
	}
	return events, nil
}
