package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/models"
)

type ContactFinance struct {
	ContactID      uuid.UUID              `json:"contactId"`
	IncomeSources  []models.IncomeSource  `json:"incomeSources"`
	Transactions   []models.Transaction   `json:"transactions"`
	IncomeCents    int64                  `json:"incomeCents"`
	ExpenseCents   int64                  `json:"expenseCents"`
}

func (s *Service) ContactFinance(ctx context.Context, contactID uuid.UUID) (*ContactFinance, error) {
	if err := s.validateContactID(ctx, &contactID); err != nil {
		return nil, err
	}
	incomes, err := s.ListIncomeSourcesFiltered(ctx, IncomeSourceFilter{ContactID: &contactID})
	if err != nil {
		return nil, err
	}
	txs, err := s.ListTransactions(ctx, TransactionFilter{ContactID: &contactID})
	if err != nil {
		return nil, err
	}
	var incomeCents, expenseCents int64
	for _, tx := range txs {
		switch tx.Type {
		case "income":
			incomeCents += tx.AmountCents
		case "expense":
			expenseCents += tx.AmountCents
		}
	}
	return &ContactFinance{
		ContactID:     contactID,
		IncomeSources: incomes,
		Transactions:  txs,
		IncomeCents:   incomeCents,
		ExpenseCents:  expenseCents,
	}, nil
}
