package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListIncomeSources(ctx context.Context, f IncomeSourceFilter) ([]models.IncomeSource, error) {
	var out []models.IncomeSource
	q := r.db.WithContext(ctx).Order("name ASC")
	if f.ActiveOnly {
		q = q.Where("active = ?", true)
	}
	if f.ContactID != nil {
		q = q.Where("contact_id = ?", *f.ContactID)
	}
	if f.ProjectID != nil {
		q = q.Where("project_id = ?", *f.ProjectID)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list income sources: %w", err)
	}
	return out, nil
}

type IncomeSourceFilter struct {
	ActiveOnly bool
	ContactID  *uuid.UUID
	ProjectID  *uuid.UUID
}

func (r *Repository) FindIncomeSource(ctx context.Context, id uuid.UUID) (*models.IncomeSource, error) {
	var row models.IncomeSource
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find income source: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateIncomeSource(ctx context.Context, row *models.IncomeSource) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create income source: %w", err)
	}
	return nil
}

func (r *Repository) SaveIncomeSource(ctx context.Context, row *models.IncomeSource) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save income source: %w", err)
	}
	return nil
}

func (r *Repository) DeleteIncomeSource(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.IncomeSource{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete income source: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) ListExpenses(ctx context.Context, activeOnly bool) ([]models.Expense, error) {
	var out []models.Expense
	q := r.db.WithContext(ctx).Order("name ASC")
	if activeOnly {
		q = q.Where("active = ?", true)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list expenses: %w", err)
	}
	return out, nil
}

func (r *Repository) FindExpense(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	var row models.Expense
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find expense: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateExpense(ctx context.Context, row *models.Expense) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create expense: %w", err)
	}
	return nil
}

func (r *Repository) SaveExpense(ctx context.Context, row *models.Expense) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save expense: %w", err)
	}
	return nil
}

func (r *Repository) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.Expense{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete expense: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type TransactionFilter struct {
	Type      string
	Year      int
	Month     int
	ProjectID *uuid.UUID
	ContactID *uuid.UUID
}

func (r *Repository) ListTransactions(ctx context.Context, f TransactionFilter) ([]models.Transaction, error) {
	var out []models.Transaction
	q := r.db.WithContext(ctx).Order("date DESC, created_at DESC")
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.Year > 0 && f.Month > 0 {
		start := time.Date(f.Year, time.Month(f.Month), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		q = q.Where("date >= ? AND date < ?", start, end)
	}
	if f.ProjectID != nil {
		q = q.Where("project_id = ?", *f.ProjectID)
	}
	if f.ContactID != nil {
		q = q.Where("contact_id = ?", *f.ContactID)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	return out, nil
}

func (r *Repository) FindTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var row models.Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find transaction: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateTransaction(ctx context.Context, row *models.Transaction) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func (r *Repository) SaveTransaction(ctx context.Context, row *models.Transaction) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save transaction: %w", err)
	}
	return nil
}

func (r *Repository) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.Transaction{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete transaction: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type MonthTotals struct {
	IncomeCents  int64
	ExpenseCents int64
}

func (r *Repository) SumTransactionsByMonth(ctx context.Context, year, month int) (MonthTotals, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	var totals MonthTotals
	type row struct {
		Type        string
		AmountCents int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Select("type, COALESCE(SUM(amount_cents), 0) as amount_cents").
		Where("date >= ? AND date < ?", start, end).
		Group("type").
		Scan(&rows).Error
	if err != nil {
		return totals, fmt.Errorf("sum transactions: %w", err)
	}
	for _, rw := range rows {
		switch rw.Type {
		case "income":
			totals.IncomeCents = rw.AmountCents
		case "expense":
			totals.ExpenseCents = rw.AmountCents
		}
	}
	return totals, nil
}

func (r *Repository) ListInvoices(ctx context.Context, status string) ([]models.Invoice, error) {
	var out []models.Invoice
	q := r.db.WithContext(ctx).Order("due_date DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	return out, nil
}

func (r *Repository) FindInvoice(ctx context.Context, id uuid.UUID) (*models.Invoice, error) {
	var row models.Invoice
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find invoice: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateInvoice(ctx context.Context, row *models.Invoice) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create invoice: %w", err)
	}
	return nil
}

func (r *Repository) SaveInvoice(ctx context.Context, row *models.Invoice) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save invoice: %w", err)
	}
	return nil
}

func (r *Repository) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.InvoiceItem{}, "invoice_id = ?", id).Error; err != nil {
			return fmt.Errorf("delete invoice items: %w", err)
		}
		res := tx.Delete(&models.Invoice{}, "id = ?", id)
		if res.Error != nil {
			return fmt.Errorf("delete invoice: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (r *Repository) CreateInvoiceItem(ctx context.Context, row *models.InvoiceItem) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create invoice item: %w", err)
	}
	return nil
}

func (r *Repository) DeleteInvoiceItem(ctx context.Context, invoiceID, itemID uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.InvoiceItem{}, "id = ? AND invoice_id = ?", itemID, invoiceID)
	if res.Error != nil {
		return fmt.Errorf("delete invoice item: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) RecalcInvoiceTotal(ctx context.Context, invoiceID uuid.UUID) error {
	var sum int64
	if err := r.db.WithContext(ctx).Model(&models.InvoiceItem{}).
		Where("invoice_id = ?", invoiceID).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&sum).Error; err != nil {
		return fmt.Errorf("sum invoice items: %w", err)
	}
	return r.db.WithContext(ctx).Model(&models.Invoice{}).
		Where("id = ?", invoiceID).
		Update("total_cents", sum).Error
}

func (r *Repository) ListInvoicesDueBetween(ctx context.Context, from, to time.Time) ([]models.Invoice, error) {
	var out []models.Invoice
	err := r.db.WithContext(ctx).
		Where("due_date >= ? AND due_date <= ?", from, to).
		Order("due_date ASC").
		Find(&out).Error
	if err != nil {
		return nil, fmt.Errorf("list invoices due: %w", err)
	}
	return out, nil
}

func (r *Repository) CountOpenInvoices(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Invoice{}).Where("status = ?", "open").Count(&n).Error
	return n, err
}

func (r *Repository) ListBudgets(ctx context.Context, year, month int) ([]models.BudgetPlan, error) {
	var out []models.BudgetPlan
	q := r.db.WithContext(ctx).Order("category ASC")
	if year > 0 {
		q = q.Where("year = ?", year)
	}
	if month > 0 {
		q = q.Where("month = ?", month)
	}
	if err := q.Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list budgets: %w", err)
	}
	return out, nil
}

func (r *Repository) FindBudget(ctx context.Context, id uuid.UUID) (*models.BudgetPlan, error) {
	var row models.BudgetPlan
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find budget: %w", err)
	}
	return &row, nil
}

func (r *Repository) CreateBudget(ctx context.Context, row *models.BudgetPlan) error {
	if row.ID == uuid.Nil {
		row.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return fmt.Errorf("create budget: %w", err)
	}
	return nil
}

func (r *Repository) SaveBudget(ctx context.Context, row *models.BudgetPlan) error {
	if err := r.db.WithContext(ctx).Save(row).Error; err != nil {
		return fmt.Errorf("save budget: %w", err)
	}
	return nil
}

func (r *Repository) DeleteBudget(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.BudgetPlan{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete budget: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) SumExpensesByCategory(ctx context.Context, year, month int) (map[string]int64, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	type row struct {
		Category    string
		AmountCents int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Select("COALESCE(e.category, 'other') as category, COALESCE(SUM(transactions.amount_cents), 0) as amount_cents").
		Joins("LEFT JOIN expenses e ON e.id = transactions.expense_id").
		Where("transactions.type = ? AND transactions.date >= ? AND transactions.date < ?", "expense", start, end).
		Group("e.category").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sum expenses by category: %w", err)
	}
	out := make(map[string]int64, len(rows))
	for _, rw := range rows {
		cat := rw.Category
		if cat == "" {
			cat = "other"
		}
		out[cat] = rw.AmountCents
	}
	return out, nil
}
