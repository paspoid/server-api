package paspoid_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	paspoid "github.com/paspoid/server-api-go"
)

func TestClient_GetKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/ext/get-key" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req["api_key"] != "test-key" || req["api_secret"] != "test-secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"key":               "nonce-123",
			"validation_window": "30s",
		})
	}))
	defer ts.Close()

	client := paspoid.NewClient(ts.URL, "test-key", "test-secret")
	resp, err := client.GetKey(context.Background(), "srv-1", "phones")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Key != "nonce-123" {
		t.Errorf("expected key 'nonce-123', got '%s'", resp.Key)
	}
	if resp.ValidationWindow != "30s" {
		t.Errorf("expected validation_window '30s', got '%s'", resp.ValidationWindow)
	}
}

func TestClient_Validate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/ext/validate" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req["nonce"] != "nonce-123" {
			http.Error(w, "invalid nonce", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":     "success",
			"data_type":  "phone",
			"data_value": "+77000000000",
		})
	}))
	defer ts.Close()

	client := paspoid.NewClient(ts.URL, "test-key", "test-secret")
	resp, err := client.Validate(context.Background(), "nonce-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("expected status 'success', got '%s'", resp.Status)
	}
	if resp.DataType == nil || *resp.DataType != "phone" {
		t.Errorf("expected data_type 'phone', got %v", resp.DataType)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"service not found"}`))
	}))
	defer ts.Close()

	client := paspoid.NewClient(ts.URL, "test-key", "test-secret")
	_, err := client.GetKey(context.Background(), "missing", "phones")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !contains(err.Error(), "status 404") {
		t.Errorf("expected error to contain status 404, got '%v'", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
