package certifications

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http/httptest"
    "testing"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    glebarezsqlite "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "log/slog"
)

func setupApp(t *testing.T) *fiber.App {
    db, err := gorm.Open(glebarezsqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to open sqlite: %v", err)
    }

    if err := db.AutoMigrate(&Certification{}); err != nil {
        t.Fatalf("auto migrate failed: %v", err)
    }

    app := fiber.New()

    // Middleware to set authenticated user in context
    app.Use(func(c *fiber.Ctx) error {
        c.Locals("userID", uuid.MustParse("11111111-1111-1111-1111-111111111111"))
        return c.Next()
    })

    logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

    api := app.Group("/api/v1")
    SetupRoutes(api.Group("/certifications"), db, logger)

    return app
}

func TestCertificationsCRUD(t *testing.T) {
    app := setupApp(t)

    // Create
    payload := map[string]string{
        "name": "Golang Pro",
        "issuer": "Acme",
        "date": "2024-12",
        "url": "https://example.com/cert",
        "description": "Certified Go developer",
    }
    b, _ := json.Marshal(payload)
    req := httptest.NewRequest("POST", "/api/v1/certifications", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("create request failed: %v", err)
    }
    if resp.StatusCode != 201 {
        t.Fatalf("expected 201 created, got %d", resp.StatusCode)
    }
    var created map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
        t.Fatalf("failed decode create response: %v", err)
    }
    data := created["data"].(map[string]interface{})
    id := data["id"].(string)

    // List
    req = httptest.NewRequest("GET", "/api/v1/certifications", nil)
    resp, err = app.Test(req)
    if err != nil {
        t.Fatalf("list request failed: %v", err)
    }
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200 list, got %d", resp.StatusCode)
    }
    var listResp map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
        t.Fatalf("failed decode list response: %v", err)
    }
    items := listResp["data"].([]interface{})
    if len(items) != 1 {
        t.Fatalf("expected 1 item, got %d", len(items))
    }

    // Get
    req = httptest.NewRequest("GET", "/api/v1/certifications/"+id, nil)
    resp, err = app.Test(req)
    if err != nil {
        t.Fatalf("get request failed: %v", err)
    }
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200 get, got %d", resp.StatusCode)
    }

    // Update
    updatePayload := map[string]string{"issuer":"New Issuer"}
    b, _ = json.Marshal(updatePayload)
    req = httptest.NewRequest("PUT", "/api/v1/certifications/"+id, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err = app.Test(req)
    if err != nil {
        t.Fatalf("update request failed: %v", err)
    }
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200 update, got %d", resp.StatusCode)
    }
    var updated map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
        t.Fatalf("failed decode update response: %v", err)
    }
    updatedData := updated["data"].(map[string]interface{})
    if updatedData["issuer"].(string) != "New Issuer" {
        t.Fatalf("issuer not updated")
    }

    // Delete
    req = httptest.NewRequest("DELETE", "/api/v1/certifications/"+id, nil)
    resp, err = app.Test(req)
    if err != nil {
        t.Fatalf("delete request failed: %v", err)
    }
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200 delete, got %d", resp.StatusCode)
    }

    // List again
    req = httptest.NewRequest("GET", "/api/v1/certifications", nil)
    resp, err = app.Test(req)
    if err != nil {
        t.Fatalf("list request failed: %v", err)
    }
    var finalList map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&finalList); err != nil {
        t.Fatalf("failed decode final list response: %v", err)
    }
    items = finalList["data"].([]interface{})
    if len(items) != 0 {
        t.Fatalf("expected 0 items after delete, got %d", len(items))
    }
}
