package httpserver

import (
	"net/http"

	"github.com/woragis/management/backend/server/internal/middleware"
)

func Mount(mux *http.ServeMux, app *App) {
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /ready", handleReady(app.DB))

	admin := func(h http.HandlerFunc) http.Handler {
		return middleware.AdminAuth(app.AdminAPIKey, h)
	}

	if app.DevProjects != nil {
		mux.Handle("GET /v1/admin/dashboard", admin(handleDashboard(app.DevProjects, app.MediaRepo)))
	}

	if app.DevProjects != nil {
		dh := newDevprojectHandler(app.DevProjects)
		mux.Handle("GET /v1/admin/projects", admin(dh.list))
		mux.Handle("POST /v1/admin/projects", admin(dh.create))
		mux.Handle("GET /v1/admin/projects/{id}", admin(dh.get))
		mux.Handle("PATCH /v1/admin/projects/{id}", admin(dh.update))
		mux.Handle("DELETE /v1/admin/projects/{id}", admin(dh.delete))
		mux.Handle("POST /v1/admin/projects/{id}/links", admin(dh.createLink))
		mux.Handle("DELETE /v1/admin/projects/{id}/links/{linkId}", admin(dh.deleteLink))
		mux.Handle("POST /v1/admin/projects/{id}/domains", admin(dh.createDomain))
		mux.Handle("DELETE /v1/admin/projects/{id}/domains/{domainId}", admin(dh.deleteDomain))
		mux.Handle("GET /v1/admin/projects/{id}/secrets", admin(dh.listSecrets))
		mux.Handle("GET /v1/admin/projects/{id}/secrets/{secretId}", admin(dh.getSecret))
		mux.Handle("POST /v1/admin/projects/{id}/secrets", admin(dh.createSecret))
		mux.Handle("DELETE /v1/admin/projects/{id}/secrets/{secretId}", admin(dh.deleteSecret))
		mux.Handle("POST /v1/admin/projects/{id}/gallery", admin(dh.createGallery))
		mux.Handle("DELETE /v1/admin/projects/{id}/gallery/{itemId}", admin(dh.deleteGallery))
		mux.Handle("GET /v1/admin/projects/{id}/envs", admin(dh.listEnvs))
		mux.Handle("POST /v1/admin/projects/{id}/envs", admin(dh.createEnv))
		mux.Handle("DELETE /v1/admin/projects/{id}/envs/{envId}", admin(dh.deleteEnv))
	}

	if app.Media != nil {
		mh := newMediaHandler(app.Media)
		admin := func(h http.HandlerFunc) http.Handler {
			return middleware.AdminAuth(app.AdminAPIKey, h)
		}
		mux.Handle("GET /v1/admin/media", admin(mh.list))
		mux.Handle("POST /v1/admin/media", admin(mh.upload))
		mux.Handle("GET /v1/admin/media/{id}", admin(mh.get))
		mux.Handle("DELETE /v1/admin/media/{id}", admin(mh.delete))
		mux.HandleFunc("GET /v1/public/media/{id}/file", mh.serveFile)
		mux.HandleFunc("GET /v1/public/media/{id}", mh.get)
	}

	if app.Profile != nil {
		ph := newProfileHandler(app.Profile)
		admin := func(h http.HandlerFunc) http.Handler {
			return middleware.AdminAuth(app.AdminAPIKey, h)
		}
		mux.Handle("GET /v1/admin/profile", admin(ph.getAdmin))
		mux.Handle("PATCH /v1/admin/profile", admin(ph.update))
		mux.HandleFunc("GET /v1/public/profile", ph.getPublic)
	}

	if app.DevProjects != nil {
		pub := newPublicHandler(app.DevProjects)
		mux.HandleFunc("GET /v1/public/projects", pub.listProjects)
		mux.HandleFunc("GET /v1/public/projects/{slug}", pub.getProject)
	}

	if app.Finance != nil {
		fh := newFinanceHandler(app.Finance)
		mux.Handle("GET /v1/admin/finance/dashboard", admin(fh.dashboard))
		mux.Handle("GET /v1/admin/finance/summary", admin(fh.summary))
		mux.Handle("GET /v1/admin/finance/calendar", admin(fh.calendar))
		mux.Handle("GET /v1/admin/finance/income-sources", admin(fh.listIncomeSources))
		mux.Handle("POST /v1/admin/finance/income-sources", admin(fh.createIncomeSource))
		mux.Handle("GET /v1/admin/finance/income-sources/{id}", admin(fh.getIncomeSource))
		mux.Handle("PATCH /v1/admin/finance/income-sources/{id}", admin(fh.updateIncomeSource))
		mux.Handle("DELETE /v1/admin/finance/income-sources/{id}", admin(fh.deleteIncomeSource))
		mux.Handle("GET /v1/admin/finance/expenses", admin(fh.listExpenses))
		mux.Handle("POST /v1/admin/finance/expenses", admin(fh.createExpense))
		mux.Handle("GET /v1/admin/finance/expenses/{id}", admin(fh.getExpense))
		mux.Handle("PATCH /v1/admin/finance/expenses/{id}", admin(fh.updateExpense))
		mux.Handle("DELETE /v1/admin/finance/expenses/{id}", admin(fh.deleteExpense))
		mux.Handle("GET /v1/admin/finance/transactions", admin(fh.listTransactions))
		mux.Handle("POST /v1/admin/finance/transactions", admin(fh.createTransaction))
		mux.Handle("GET /v1/admin/finance/transactions/{id}", admin(fh.getTransaction))
		mux.Handle("PATCH /v1/admin/finance/transactions/{id}", admin(fh.updateTransaction))
		mux.Handle("DELETE /v1/admin/finance/transactions/{id}", admin(fh.deleteTransaction))
		mux.Handle("GET /v1/admin/finance/invoices", admin(fh.listInvoices))
		mux.Handle("POST /v1/admin/finance/invoices", admin(fh.createInvoice))
		mux.Handle("GET /v1/admin/finance/invoices/{id}", admin(fh.getInvoice))
		mux.Handle("PATCH /v1/admin/finance/invoices/{id}", admin(fh.updateInvoice))
		mux.Handle("DELETE /v1/admin/finance/invoices/{id}", admin(fh.deleteInvoice))
		mux.Handle("POST /v1/admin/finance/invoices/{id}/items", admin(fh.createInvoiceItem))
		mux.Handle("DELETE /v1/admin/finance/invoices/{id}/items/{itemId}", admin(fh.deleteInvoiceItem))
		mux.Handle("GET /v1/admin/finance/budgets", admin(fh.listBudgets))
		mux.Handle("POST /v1/admin/finance/budgets", admin(fh.createBudget))
		mux.Handle("GET /v1/admin/finance/budgets/{id}", admin(fh.getBudget))
		mux.Handle("PATCH /v1/admin/finance/budgets/{id}", admin(fh.updateBudget))
		mux.Handle("DELETE /v1/admin/finance/budgets/{id}", admin(fh.deleteBudget))
	}

	if app.Content != nil {
		ch := newContentHandler(app.Content)
		mux.Handle("GET /v1/admin/content/leetcode/videos", admin(ch.listVideos))
		mux.Handle("POST /v1/admin/content/leetcode/videos", admin(ch.createVideo))
		mux.Handle("GET /v1/admin/content/leetcode/videos/{id}", admin(ch.getVideo))
		mux.Handle("PATCH /v1/admin/content/leetcode/videos/{id}", admin(ch.updateVideo))
		mux.Handle("DELETE /v1/admin/content/leetcode/videos/{id}", admin(ch.deleteVideo))
		mux.Handle("GET /v1/admin/content/leetcode/videos/{id}/thumbnails", admin(ch.listThumbnails))
		mux.Handle("POST /v1/admin/content/leetcode/videos/{id}/thumbnails", admin(ch.createThumbnail))
		mux.Handle("GET /v1/admin/content/leetcode/videos/{id}/thumbnails/{thumbnailId}", admin(ch.getThumbnail))
		mux.Handle("PATCH /v1/admin/content/leetcode/videos/{id}/thumbnails/{thumbnailId}", admin(ch.updateThumbnail))
		mux.Handle("DELETE /v1/admin/content/leetcode/videos/{id}/thumbnails/{thumbnailId}", admin(ch.deleteThumbnail))
		mux.Handle("POST /v1/admin/content/leetcode/videos/{id}/thumbnails/{thumbnailId}/generate", admin(ch.generateThumbnail))
		mux.Handle("POST /v1/admin/content/leetcode/videos/{id}/thumbnails/{thumbnailId}/approve", admin(ch.approveThumbnail))
		mux.Handle("GET /v1/admin/content/leetcode/templates", admin(ch.listTemplates))
		mux.Handle("POST /v1/admin/content/leetcode/templates", admin(ch.createTemplate))
		mux.Handle("GET /v1/admin/content/leetcode/templates/{id}", admin(ch.getTemplate))
		mux.Handle("PATCH /v1/admin/content/leetcode/templates/{id}", admin(ch.updateTemplate))
		mux.Handle("DELETE /v1/admin/content/leetcode/templates/{id}", admin(ch.deleteTemplate))
		mux.HandleFunc("POST /v1/webhooks/creatives", handleCreativesWebhook(app.Content))
	}
}
