package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type validateReq struct {
    Token string `json:"token"`
}

type validateResp struct {
    Valid   bool   `json:"valid"`
    UserID  string `json:"userId,omitempty"`
    Role    string `json:"role,omitempty"`
    Message string `json:"message,omitempty"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req validateReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }

    // Simple static token for tests
    if req.Token == "test-token" {
        resp := validateResp{Valid: true, UserID: "11111111-1111-1111-1111-111111111111", Role: "user"}
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(resp)
        return
    }

    resp := validateResp{Valid: false, Message: "invalid token"}
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    http.HandleFunc("/api/v1/auth/validate", validateHandler)

    log.Printf("mock auth listening on :%s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("failed to start mock auth: %v", err)
    }
}
