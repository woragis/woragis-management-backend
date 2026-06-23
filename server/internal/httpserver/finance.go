package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/management/backend/server/internal/apperrors"
	financesvc "github.com/woragis/management/backend/server/internal/finance/service"
)

type financeHandler struct {
	svc *financesvc.Service
}

func newFinanceHandler(svc *financesvc.Service) *financeHandler {
	return &financeHandler{svc: svc}
}

func (h *financeHandler) dashboard(w http.ResponseWriter, r *http.Request) {
	d, err := h.svc.Dashboard(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, d)
}

func (h *financeHandler) summary(w http.ResponseWriter, r *http.Request) {
	year, month := parseYearMonthQuery(r)
	s, err := h.svc.MonthlySummary(r.Context(), year, month)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, s)
}

func (h *financeHandler) calendar(w http.ResponseWriter, r *http.Request) {
	year, month := parseYearMonthQuery(r)
	events, err := h.svc.Calendar(r.Context(), year, month)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, events)
}

func (h *financeHandler) listIncomeSources(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := financesvc.IncomeSourceFilter{}
	if cid := q.Get("contactId"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			f.ContactID = &id
		}
	}
	if pid := q.Get("projectId"); pid != "" {
		if id, err := uuid.Parse(pid); err == nil {
			f.ProjectID = &id
		}
	}
	rows, err := h.svc.ListIncomeSourcesFiltered(r.Context(), f)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *financeHandler) getIncomeSource(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetIncomeSource(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) createIncomeSource(w http.ResponseWriter, r *http.Request) {
	var body incomeSourceBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateIncomeSource(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) updateIncomeSource(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body incomeSourceUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateIncomeSource(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) deleteIncomeSource(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteIncomeSource(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *financeHandler) listExpenses(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListExpenses(r.Context())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *financeHandler) getExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetExpense(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) createExpense(w http.ResponseWriter, r *http.Request) {
	var body expenseBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateExpense(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) updateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body expenseUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateExpense(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) deleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteExpense(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *financeHandler) listTransactions(w http.ResponseWriter, r *http.Request) {
	f := financesvc.TransactionFilter{
		Type: r.URL.Query().Get("type"),
	}
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			f.Year = v
		}
	}
	if m := r.URL.Query().Get("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			f.Month = v
		}
	}
	if pid := r.URL.Query().Get("projectId"); pid != "" {
		if id, err := uuid.Parse(pid); err == nil {
			f.ProjectID = &id
		}
	}
	if cid := r.URL.Query().Get("contactId"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			f.ContactID = &id
		}
	}
	rows, err := h.svc.ListTransactions(r.Context(), f)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *financeHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetTransaction(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	var body transactionBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateTransaction(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body transactionUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateTransaction(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteTransaction(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *financeHandler) listInvoices(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListInvoices(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *financeHandler) getInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetInvoice(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) createInvoice(w http.ResponseWriter, r *http.Request) {
	var body invoiceBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateInvoice(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) updateInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body invoiceUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateInvoice(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) deleteInvoice(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteInvoice(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *financeHandler) createInvoiceItem(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body invoiceItemBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateInvoiceItem(r.Context(), invoiceID, body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) deleteInvoiceItem(w http.ResponseWriter, r *http.Request) {
	invoiceID, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	itemID, err := parseUUID(r.PathValue("itemId"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid item id."))
		return
	}
	if err := h.svc.DeleteInvoiceItem(r.Context(), invoiceID, itemID); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *financeHandler) listBudgets(w http.ResponseWriter, r *http.Request) {
	year, month := parseYearMonthQuery(r)
	rows, err := h.svc.ListBudgets(r.Context(), year, month)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, rows)
}

func (h *financeHandler) getBudget(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	row, err := h.svc.GetBudget(r.Context(), id)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) createBudget(w http.ResponseWriter, r *http.Request) {
	var body budgetBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.CreateBudget(r.Context(), body.toCreate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusCreated, row)
}

func (h *financeHandler) updateBudget(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	var body budgetUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Request body is invalid."))
		return
	}
	row, err := h.svc.UpdateBudget(r.Context(), id, body.toUpdate())
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	apperrors.WriteJSON(w, http.StatusOK, row)
}

func (h *financeHandler) deleteBudget(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		apperrors.WriteError(w, apperrors.Invalid(apperrors.CodeInternal, "Invalid id."))
		return
	}
	if err := h.svc.DeleteBudget(r.Context(), id); err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseYearMonthQuery(r *http.Request) (int, int) {
	year, month := 0, 0
	if y := r.URL.Query().Get("year"); y != "" {
		if v, err := strconv.Atoi(y); err == nil {
			year = v
		}
	}
	if m := r.URL.Query().Get("month"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			month = v
		}
	}
	return year, month
}

type incomeSourceBody struct {
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	AmountCents int64      `json:"amountCents"`
	Currency    string     `json:"currency"`
	Frequency   string     `json:"frequency"`
	DayOfMonth  int        `json:"dayOfMonth"`
	ProjectID   *uuid.UUID `json:"projectId"`
	ContactID   *uuid.UUID `json:"contactId"`
	Active      bool       `json:"active"`
	Notes       string     `json:"notes"`
}

func (b incomeSourceBody) toCreate() financesvc.CreateIncomeSourceInput {
	return financesvc.CreateIncomeSourceInput(b)
}

type incomeSourceUpdateBody struct {
	Name        *string    `json:"name"`
	Type        *string    `json:"type"`
	AmountCents *int64     `json:"amountCents"`
	Currency    *string    `json:"currency"`
	Frequency   *string    `json:"frequency"`
	DayOfMonth  *int       `json:"dayOfMonth"`
	ProjectID   *uuid.UUID `json:"projectId"`
	ContactID   *uuid.UUID `json:"contactId"`
	Active      *bool      `json:"active"`
	Notes       *string    `json:"notes"`
}

func (b incomeSourceUpdateBody) toUpdate() financesvc.UpdateIncomeSourceInput {
	in := financesvc.UpdateIncomeSourceInput{
		Name:        b.Name,
		Type:        b.Type,
		AmountCents: b.AmountCents,
		Currency:    b.Currency,
		Frequency:   b.Frequency,
		DayOfMonth:  b.DayOfMonth,
		Active:      b.Active,
		Notes:       b.Notes,
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	if b.ContactID != nil {
		in.ContactID = b.ContactID
		in.ContactSet = true
	}
	return in
}

type expenseBody struct {
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	AmountCents int64      `json:"amountCents"`
	Currency    string     `json:"currency"`
	Frequency   string     `json:"frequency"`
	DayOfMonth  int        `json:"dayOfMonth"`
	DueDate     *time.Time `json:"dueDate"`
	AutoPay     bool       `json:"autoPay"`
	ProjectID   *uuid.UUID `json:"projectId"`
	Active      bool       `json:"active"`
	Notes       string     `json:"notes"`
}

func (b expenseBody) toCreate() financesvc.CreateExpenseInput {
	return financesvc.CreateExpenseInput(b)
}

type expenseUpdateBody struct {
	Name        *string    `json:"name"`
	Category    *string    `json:"category"`
	AmountCents *int64     `json:"amountCents"`
	Currency    *string    `json:"currency"`
	Frequency   *string    `json:"frequency"`
	DayOfMonth  *int       `json:"dayOfMonth"`
	DueDate     *time.Time `json:"dueDate"`
	AutoPay     *bool      `json:"autoPay"`
	ProjectID   *uuid.UUID `json:"projectId"`
	Active      *bool      `json:"active"`
	Notes       *string    `json:"notes"`
}

func (b expenseUpdateBody) toUpdate() financesvc.UpdateExpenseInput {
	in := financesvc.UpdateExpenseInput{
		Name:        b.Name,
		Category:    b.Category,
		AmountCents: b.AmountCents,
		Currency:    b.Currency,
		Frequency:   b.Frequency,
		DayOfMonth:  b.DayOfMonth,
		AutoPay:     b.AutoPay,
		Active:      b.Active,
		Notes:       b.Notes,
	}
	if b.DueDate != nil {
		in.DueDate = b.DueDate
		in.DueDateSet = true
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	return in
}

type transactionBody struct {
	Type           string     `json:"type"`
	AmountCents    int64      `json:"amountCents"`
	Currency       string     `json:"currency"`
	Description    string     `json:"description"`
	Date           time.Time  `json:"date"`
	IncomeSourceID *uuid.UUID `json:"incomeSourceId"`
	ExpenseID      *uuid.UUID `json:"expenseId"`
	ProjectID      *uuid.UUID `json:"projectId"`
	ContactID      *uuid.UUID `json:"contactId"`
	InvoiceID      *uuid.UUID `json:"invoiceId"`
	Notes          string     `json:"notes"`
}

func (b transactionBody) toCreate() financesvc.CreateTransactionInput {
	return financesvc.CreateTransactionInput(b)
}

type transactionUpdateBody struct {
	Type            *string    `json:"type"`
	AmountCents     *int64     `json:"amountCents"`
	Currency        *string    `json:"currency"`
	Description     *string    `json:"description"`
	Date            *time.Time `json:"date"`
	IncomeSourceID  *uuid.UUID `json:"incomeSourceId"`
	ExpenseID       *uuid.UUID `json:"expenseId"`
	ProjectID       *uuid.UUID `json:"projectId"`
	ContactID       *uuid.UUID `json:"contactId"`
	InvoiceID       *uuid.UUID `json:"invoiceId"`
	Notes           *string    `json:"notes"`
}

func (b transactionUpdateBody) toUpdate() financesvc.UpdateTransactionInput {
	in := financesvc.UpdateTransactionInput{
		Type:        b.Type,
		AmountCents: b.AmountCents,
		Currency:    b.Currency,
		Description: b.Description,
		Date:        b.Date,
		Notes:       b.Notes,
	}
	if b.IncomeSourceID != nil {
		in.IncomeSourceID = b.IncomeSourceID
		in.IncomeSourceSet = true
	}
	if b.ExpenseID != nil {
		in.ExpenseID = b.ExpenseID
		in.ExpenseSet = true
	}
	if b.ProjectID != nil {
		in.ProjectID = b.ProjectID
		in.ProjectSet = true
	}
	if b.ContactID != nil {
		in.ContactID = b.ContactID
		in.ContactSet = true
	}
	if b.InvoiceID != nil {
		in.InvoiceID = b.InvoiceID
		in.InvoiceSet = true
	}
	return in
}

type invoiceBody struct {
	Name         string     `json:"name"`
	CardLastFour string     `json:"cardLastFour"`
	DueDate      time.Time  `json:"dueDate"`
	ClosedAt     *time.Time `json:"closedAt"`
	TotalCents   int64      `json:"totalCents"`
	PaidCents    int64      `json:"paidCents"`
	Status       string     `json:"status"`
	Notes        string     `json:"notes"`
}

func (b invoiceBody) toCreate() financesvc.CreateInvoiceInput {
	return financesvc.CreateInvoiceInput(b)
}

type invoiceUpdateBody struct {
	Name         *string    `json:"name"`
	CardLastFour *string    `json:"cardLastFour"`
	DueDate      *time.Time `json:"dueDate"`
	ClosedAt     *time.Time `json:"closedAt"`
	TotalCents   *int64     `json:"totalCents"`
	PaidCents    *int64     `json:"paidCents"`
	Status       *string    `json:"status"`
	Notes        *string    `json:"notes"`
}

func (b invoiceUpdateBody) toUpdate() financesvc.UpdateInvoiceInput {
	in := financesvc.UpdateInvoiceInput{
		Name:         b.Name,
		CardLastFour: b.CardLastFour,
		DueDate:      b.DueDate,
		TotalCents:   b.TotalCents,
		PaidCents:    b.PaidCents,
		Status:       b.Status,
		Notes:        b.Notes,
	}
	if b.ClosedAt != nil {
		in.ClosedAt = b.ClosedAt
		in.ClosedAtSet = true
	}
	return in
}

type invoiceItemBody struct {
	Description string    `json:"description"`
	AmountCents int64     `json:"amountCents"`
	Date        time.Time `json:"date"`
	Category    string    `json:"category"`
	Installment string    `json:"installment"`
	Notes       string    `json:"notes"`
}

func (b invoiceItemBody) toCreate() financesvc.CreateInvoiceItemInput {
	return financesvc.CreateInvoiceItemInput(b)
}

type budgetBody struct {
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Category     string `json:"category"`
	PlannedCents int64  `json:"plannedCents"`
	Notes        string `json:"notes"`
}

func (b budgetBody) toCreate() financesvc.CreateBudgetInput {
	return financesvc.CreateBudgetInput(b)
}

type budgetUpdateBody struct {
	Year         *int    `json:"year"`
	Month        *int    `json:"month"`
	Category     *string `json:"category"`
	PlannedCents *int64  `json:"plannedCents"`
	Notes        *string `json:"notes"`
}

func (b budgetUpdateBody) toUpdate() financesvc.UpdateBudgetInput {
	return financesvc.UpdateBudgetInput(b)
}
