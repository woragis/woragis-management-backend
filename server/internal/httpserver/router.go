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

	if app.Contacts != nil {
		ch := newContactsHandler(app.Contacts, app.Finance)
		mux.Handle("GET /v1/admin/contacts", admin(ch.list))
		mux.Handle("POST /v1/admin/contacts", admin(ch.create))
		mux.Handle("GET /v1/admin/contacts/{id}", admin(ch.get))
		mux.Handle("PATCH /v1/admin/contacts/{id}", admin(ch.update))
		mux.Handle("DELETE /v1/admin/contacts/{id}", admin(ch.delete))
		mux.Handle("GET /v1/admin/contacts/{id}/interactions", admin(ch.listInteractions))
		mux.Handle("POST /v1/admin/contacts/{id}/interactions", admin(ch.createInteraction))
		mux.Handle("GET /v1/admin/contacts/{id}/finance", admin(ch.contactFinance))
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
		ch := newContentHandler(app.Content, app.Messaging)
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
		mux.Handle("GET /v1/admin/content/leetcode/settings", admin(ch.getSettings))
		mux.Handle("PATCH /v1/admin/content/leetcode/settings", admin(ch.updateSettings))
		mux.Handle("GET /v1/admin/content/leetcode/whatsapp-templates", admin(ch.listWhatsappTemplates))
		mux.Handle("PATCH /v1/admin/content/leetcode/whatsapp-templates/{id}", admin(ch.updateWhatsappTemplate))
		mux.Handle("GET /v1/admin/content/leetcode/videos/{id}/whatsapp-preview", admin(ch.whatsappPreview))
		mux.Handle("POST /v1/admin/content/leetcode/videos/{id}/whatsapp-send", admin(ch.whatsappSendNow))
		mux.Handle("GET /v1/admin/content/leetcode/whatsapp/status", admin(ch.whatsappWorkerStatus))
		mux.Handle("GET /v1/admin/content/leetcode/whatsapp/qr", admin(ch.whatsappWorkerQR))
		mux.HandleFunc("POST /v1/webhooks/creatives", handleCreativesWebhook(app.Content))
	}

	if app.Content != nil && app.WorkerAPIKey != "" {
		worker := func(h http.HandlerFunc) http.Handler {
			return middleware.WorkerAuth(app.WorkerAPIKey, h)
		}
		mux.Handle("GET /v1/internal/content/leetcode/dispatch", worker(handleInternalDispatch(app.Content)))
		mux.Handle("GET /v1/internal/content/leetcode/settings", worker(handleInternalSettings(app.Content)))
		mux.Handle("PATCH /v1/internal/content/leetcode/videos/{id}/whatsapp-status", worker(handleInternalWhatsappStatus(app.Content)))
	}

	if app.Messaging != nil {
		mh := newMessagingHandler(app.Messaging, app.Scheduler)
		mux.Handle("GET /v1/admin/messaging/destinations", admin(mh.listDestinations))
		mux.Handle("POST /v1/admin/messaging/destinations", admin(mh.createDestination))
		mux.Handle("GET /v1/admin/messaging/destinations/{id}", admin(mh.getDestination))
		mux.Handle("PATCH /v1/admin/messaging/destinations/{id}", admin(mh.updateDestination))
		mux.Handle("DELETE /v1/admin/messaging/destinations/{id}", admin(mh.deleteDestination))
		mux.Handle("GET /v1/admin/messaging/templates", admin(mh.listTemplates))
		mux.Handle("POST /v1/admin/messaging/templates", admin(mh.createTemplate))
		mux.Handle("GET /v1/admin/messaging/templates/{id}", admin(mh.getTemplate))
		mux.Handle("PATCH /v1/admin/messaging/templates/{id}", admin(mh.updateTemplate))
		mux.Handle("DELETE /v1/admin/messaging/templates/{id}", admin(mh.deleteTemplate))
		mux.Handle("GET /v1/admin/messaging/jobs", admin(mh.listJobs))
		mux.Handle("POST /v1/admin/messaging/jobs", admin(mh.createJob))
		mux.Handle("GET /v1/admin/messaging/jobs/{id}", admin(mh.getJob))
		mux.Handle("PATCH /v1/admin/messaging/jobs/{id}", admin(mh.updateJob))
		mux.Handle("DELETE /v1/admin/messaging/jobs/{id}", admin(mh.deleteJob))
		mux.Handle("GET /v1/admin/messaging/deliveries", admin(mh.listDeliveries))
		mux.Handle("GET /v1/admin/messaging/catalog", admin(mh.catalogFields))
		mux.Handle("POST /v1/admin/messaging/templates/preview", admin(mh.previewTemplate))
	}

	if app.Presence != nil {
		ph := newPresenceHandler(app.Presence)
		mux.Handle("GET /v1/admin/presence/campaigns", admin(ph.listCampaigns))
		mux.Handle("POST /v1/admin/presence/campaigns", admin(ph.createCampaign))
		mux.Handle("GET /v1/admin/presence/campaigns/{id}", admin(ph.getCampaign))
		mux.Handle("PATCH /v1/admin/presence/campaigns/{id}", admin(ph.updateCampaign))
		mux.Handle("DELETE /v1/admin/presence/campaigns/{id}", admin(ph.deleteCampaign))
		mux.Handle("GET /v1/admin/presence/templates", admin(ph.listTemplates))
		mux.Handle("POST /v1/admin/presence/templates", admin(ph.createTemplate))
		mux.Handle("GET /v1/admin/presence/templates/{id}", admin(ph.getTemplate))
		mux.Handle("PATCH /v1/admin/presence/templates/{id}", admin(ph.updateTemplate))
		mux.Handle("DELETE /v1/admin/presence/templates/{id}", admin(ph.deleteTemplate))
		mux.Handle("GET /v1/admin/presence/posts", admin(ph.listPosts))
		mux.Handle("POST /v1/admin/presence/posts", admin(ph.createPost))
		mux.Handle("GET /v1/admin/presence/posts/{id}", admin(ph.getPost))
		mux.Handle("PATCH /v1/admin/presence/posts/{id}", admin(ph.updatePost))
		mux.Handle("DELETE /v1/admin/presence/posts/{id}", admin(ph.deletePost))
		mux.Handle("GET /v1/admin/presence/settings", admin(ph.getSettings))
		mux.Handle("PATCH /v1/admin/presence/settings", admin(ph.updateSettings))
	}

	if app.Messaging != nil && app.Scheduler != nil && app.WorkerAPIKey != "" {
		worker := func(h http.HandlerFunc) http.Handler {
			return middleware.WorkerAuth(app.WorkerAPIKey, h)
		}
		mux.Handle("GET /v1/internal/scheduler/due", worker(handleSchedulerDue(app.Messaging)))
		mux.Handle("POST /v1/internal/scheduler/jobs/{id}/execute", worker(handleSchedulerExecute(app.Scheduler)))
		if app.Presence != nil {
			mux.Handle("GET /v1/internal/presence/reminders/due", worker(handlePresenceDueReminders(app.Presence)))
		}
		if app.PresenceReminder != nil {
			mux.Handle("POST /v1/internal/presence/reminders/{id}/send", worker(handlePresenceSendReminder(app.PresenceReminder)))
		}
	}

	if app.Personality != nil {
		ph := newAgentPersonalityHandler(app.Personality)
		mux.Handle("GET /v1/admin/agent/personality", admin(ph.get))
		mux.Handle("PATCH /v1/admin/agent/personality", admin(ph.update))
		mux.Handle("POST /v1/admin/agent/personality/reset", admin(ph.reset))
	}

	if app.AgentAPIKey != "" {
		agent := func(h http.HandlerFunc) http.Handler {
			return middleware.AgentAuth(app.AgentAPIKey, h)
		}
		tools := newAgentToolsHandler(app)
		mux.Handle("GET /v1/internal/agent/personality", agent(tools.getPersonality))
		mux.Handle("PATCH /v1/internal/agent/personality", agent(tools.updatePersonality))
		mux.Handle("POST /v1/internal/agent/personality/reset", agent(tools.resetPersonality))

		mux.Handle("GET /v1/internal/agent/tools/contacts", agent(tools.searchContacts))
		mux.Handle("POST /v1/internal/agent/tools/contacts", agent(tools.createContact))
		mux.Handle("GET /v1/internal/agent/tools/contacts/due-follow-up", agent(tools.listContactsDueFollowUp))
		mux.Handle("GET /v1/internal/agent/tools/contacts/{id}", agent(tools.getContact))
		mux.Handle("PATCH /v1/internal/agent/tools/contacts/{id}", agent(tools.updateContact))
		mux.Handle("POST /v1/internal/agent/tools/contacts/{id}/interactions", agent(tools.logInteraction))
		mux.Handle("GET /v1/internal/agent/tools/contacts/{id}/finance", agent(tools.getContactFinance))

		mux.Handle("GET /v1/internal/agent/tools/projects", agent(tools.listProjects))
		mux.Handle("POST /v1/internal/agent/tools/projects", agent(tools.createProject))
		mux.Handle("GET /v1/internal/agent/tools/projects/{id}", agent(tools.getProject))

		mux.Handle("GET /v1/internal/agent/tools/finance/dashboard", agent(tools.financeDashboard))
		mux.Handle("GET /v1/internal/agent/tools/finance/summary", agent(tools.financeSummary))
		mux.Handle("GET /v1/internal/agent/tools/finance/calendar", agent(tools.financeCalendar))
		mux.Handle("GET /v1/internal/agent/tools/finance/income-sources", agent(tools.listIncomeSources))
		mux.Handle("POST /v1/internal/agent/tools/finance/income-sources", agent(tools.createIncomeSource))
		mux.Handle("GET /v1/internal/agent/tools/finance/transactions", agent(tools.listTransactions))
		mux.Handle("POST /v1/internal/agent/tools/finance/transactions", agent(tools.createTransaction))

		mux.Handle("GET /v1/internal/agent/tools/presence/posts", agent(tools.listSocialPosts))
		mux.Handle("POST /v1/internal/agent/tools/presence/posts", agent(tools.createSocialPost))
		mux.Handle("GET /v1/internal/agent/tools/presence/templates", agent(tools.listPostTemplates))
		mux.Handle("POST /v1/internal/agent/tools/presence/apply-template", agent(tools.applyPostTemplate))
	}
}
