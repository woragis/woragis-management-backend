package models

import (
	"time"

	"github.com/google/uuid"
)

type IncomeSource struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string     `gorm:"size:200;not null" json:"name"`
	Type         string     `gorm:"size:32;not null;default:other" json:"type"`
	AmountCents  int64      `gorm:"column:amount_cents;not null;default:0" json:"amountCents"`
	Currency     string     `gorm:"size:8;not null;default:BRL" json:"currency"`
	Frequency    string     `gorm:"size:32;not null;default:monthly" json:"frequency"`
	DayOfMonth   int        `gorm:"column:day_of_month;not null;default:1" json:"dayOfMonth"`
	ProjectID    *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	Active       bool       `gorm:"not null;default:true" json:"active"`
	Notes        string     `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type Expense struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string     `gorm:"size:200;not null" json:"name"`
	Category    string     `gorm:"size:64;not null;default:other" json:"category"`
	AmountCents int64      `gorm:"column:amount_cents;not null;default:0" json:"amountCents"`
	Currency    string     `gorm:"size:8;not null;default:BRL" json:"currency"`
	Frequency   string     `gorm:"size:32;not null;default:monthly" json:"frequency"`
	DayOfMonth  int        `gorm:"column:day_of_month;not null;default:1" json:"dayOfMonth"`
	DueDate     *time.Time `gorm:"column:due_date" json:"dueDate"`
	AutoPay     bool       `gorm:"column:auto_pay;not null;default:false" json:"autoPay"`
	ProjectID   *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	Active      bool       `gorm:"not null;default:true" json:"active"`
	Notes       string     `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type Transaction struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Type           string     `gorm:"size:16;not null" json:"type"`
	AmountCents    int64      `gorm:"column:amount_cents;not null" json:"amountCents"`
	Currency       string     `gorm:"size:8;not null;default:BRL" json:"currency"`
	Description    string     `gorm:"size:500;not null" json:"description"`
	Date           time.Time  `gorm:"type:date;not null;index" json:"date"`
	IncomeSourceID *uuid.UUID `gorm:"column:income_source_id;type:uuid;index" json:"incomeSourceId"`
	ExpenseID      *uuid.UUID `gorm:"column:expense_id;type:uuid;index" json:"expenseId"`
	ProjectID      *uuid.UUID `gorm:"column:project_id;type:uuid;index" json:"projectId"`
	InvoiceID      *uuid.UUID `gorm:"column:invoice_id;type:uuid;index" json:"invoiceId"`
	Notes          string     `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type Invoice struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string        `gorm:"size:200;not null" json:"name"`
	CardLastFour string        `gorm:"column:card_last_four;size:4" json:"cardLastFour"`
	DueDate      time.Time     `gorm:"column:due_date;type:date;not null;index" json:"dueDate"`
	ClosedAt     *time.Time    `gorm:"column:closed_at" json:"closedAt"`
	TotalCents   int64         `gorm:"column:total_cents;not null;default:0" json:"totalCents"`
	PaidCents    int64         `gorm:"column:paid_cents;not null;default:0" json:"paidCents"`
	Status       string        `gorm:"size:16;not null;default:open" json:"status"`
	Notes        string        `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	Items        []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"items,omitempty"`
}

type InvoiceItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	InvoiceID   uuid.UUID `gorm:"column:invoice_id;type:uuid;not null;index" json:"invoiceId"`
	Description string    `gorm:"size:500;not null" json:"description"`
	AmountCents int64     `gorm:"column:amount_cents;not null" json:"amountCents"`
	Date        time.Time `gorm:"type:date;not null" json:"date"`
	Category    string    `gorm:"size:64" json:"category"`
	Installment string    `gorm:"size:32" json:"installment"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type BudgetPlan struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Year         int       `gorm:"not null;index:idx_budget_period" json:"year"`
	Month        int       `gorm:"not null;index:idx_budget_period" json:"month"`
	Category     string    `gorm:"size:64;not null;index:idx_budget_period" json:"category"`
	PlannedCents int64     `gorm:"column:planned_cents;not null;default:0" json:"plannedCents"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
