package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	"github.com/woragis/management/backend/server/internal/content/repository"
	"github.com/woragis/management/backend/server/internal/middleware"
	"github.com/woragis/management/backend/server/internal/testutil"
)

func TestAdminWhatsappPreviewE2E(t *testing.T) {
	db := testutil.OpenSQLite(t)
	content := contentsvc.New(repository.New(db), nil, nil, nil, "", "")
	ctx := t.Context()
	if err := content.EnsureWhatsappDefaults(ctx); err != nil {
		t.Fatal(err)
	}

	title := "Valid Parentheses"
	video, err := content.CreateVideo(ctx, contentsvc.CreateVideoInput{
		Title:        "LC #20",
		Status:       "draft",
		ProblemTitle: &title,
	})
	if err != nil {
		t.Fatal(err)
	}

	ch := newContentHandler(content, nil)
	preview := middleware.AdminAuth("admin-test", http.HandlerFunc(ch.whatsappPreview))

	req := httptest.NewRequest(http.MethodGet, "/v1/admin/content/leetcode/videos/"+video.ID.String()+"/whatsapp-preview", nil)
	req.SetPathValue("id", video.ID.String())
	req.Header.Set("X-Admin-Key", "admin-test")
	rec := httptest.NewRecorder()
	preview.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body.Message, "Valid Parentheses") {
		t.Fatalf("unexpected message: %q", body.Message)
	}
}

func TestInternalDispatchAuthE2E(t *testing.T) {
	db := testutil.OpenSQLite(t)
	content := contentsvc.New(repository.New(db), nil, nil, nil, "", "")

	worker := middleware.WorkerAuth("worker-test", handleInternalDispatch(content))
	req := httptest.NewRequest(http.MethodGet, "/v1/internal/content/leetcode/dispatch?type=problem", nil)
	rec := httptest.NewRecorder()
	worker.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	req.Header.Set("X-Worker-Key", "worker-test")
	rec = httptest.NewRecorder()
	worker.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body %s", rec.Code, rec.Body.String())
	}
}
