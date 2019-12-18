package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"ServerInstance", args{&Config{"key"}}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServer(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_handleUnprotectedAPI(t *testing.T) {
	type args struct {
		config *Config
	}
	type jsonResponse struct {
		Payload string `json:"data"`
	}
	tests := []struct {
		name           string
		config         Config
		want           *jsonResponse
		wantStatusCode int
	}{
		{"Simple", Config{"Key"}, &jsonResponse{"some unprotected data to send back"}, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := NewServer(&tt.config)
			r := httptest.NewRequest(http.MethodGet, "/unprotectedAPI", nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)
			if w.Code != tt.wantStatusCode {
				t.Fatalf("handleUnprotectedAPI() = %v, want %v", w.Code, tt.wantStatusCode)
			}
			if tt.want != nil {
				// Check the json payload
				got := &jsonResponse{}
				err := json.NewDecoder(w.Body).Decode(got)
				if err != nil {
					t.Fatal("handleUnprotectedAPI() failed to decode response", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Server.handleunprotectedAPI() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_handleProtectedAPI(t *testing.T) {
	type response struct {
		Payload  string `json:"data"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}
	config := Config{"secreteKey"}
	// Generate tokens needed for test
	tokenAuthorized, err := GenerateToken("guest", false, 90, []byte(config.JWTKey))
	if err != nil {
		t.Fatal("GenerateToken() failed:", err)
	}
	tests := []struct {
		name           string
		token          string
		want           *response
		wantStatusCode int
	}{
		{"Authorized", tokenAuthorized, &response{"some data to send back", "guest", false}, http.StatusOK},
		{"UnAuthorized", "ab" + tokenAuthorized[2:], nil, http.StatusUnauthorized},
		{"MissingAuthorization", "", nil, http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := NewServer(&config)
			r := httptest.NewRequest(http.MethodGet, "/protectedAPI", nil)
			if len(tt.token) > 0 {
				var bearer = "Bearer " + tt.token
				r.Header.Add("Authorization", bearer)
			}
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)
			if w.Code != tt.wantStatusCode {
				t.Fatalf("Server.handleProtectedAPI() = %v, want %v", w.Code, tt.wantStatusCode)
			}
			if tt.want != nil {
				// Check the json payload
				got := &response{}
				err := json.NewDecoder(w.Body).Decode(got)
				if err != nil {
					t.Fatal("Server.handleProtectedAPI() failed to decode response", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Server.handleProtectedAPI() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestServer_handleAdminProtectedAPI(t *testing.T) {
	type response struct {
		Payload  string `json:"data"`
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}
	config := Config{"secreteKey"}
	// Generate tokens needed for test
	tokenAuthorized, err := GenerateToken("guest", false, 90, []byte(config.JWTKey))
	if err != nil {
		t.Fatalf("GenerateToken(%s) failed: %s", "guest", err)
	}
	tokenAdminAuthorized, err := GenerateToken("sue", true, 90, []byte(config.JWTKey))
	if err != nil {
		t.Fatalf("GenerateToken(%s) failed: %s", "sue", err)
	}
	tests := []struct {
		name           string
		token          string
		want           *response
		wantStatusCode int
	}{
		{"AdminAuthorized", tokenAdminAuthorized, &response{"some data to send back", "sue", true}, http.StatusOK},
		{"Authorized", tokenAuthorized, nil, http.StatusForbidden},
		{"UnAuthorized", "ab" + tokenAuthorized[2:], nil, http.StatusUnauthorized},
		{"MissingAuthorization", "", nil, http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, _ := NewServer(&config)
			r := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if len(tt.token) > 0 {
				var bearer = "Bearer " + tt.token
				r.Header.Add("Authorization", bearer)
			}
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)
			if w.Code != tt.wantStatusCode {
				t.Fatalf("Server.handleAdminProtected() = %v, want %v", w.Code, tt.wantStatusCode)
			}
			if tt.want != nil {
				// Check the json payload
				got := &response{}
				err := json.NewDecoder(w.Body).Decode(got)
				if err != nil {
					t.Fatal("Server.handleAdminProtected() failed to decode response", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Server.handleAdminProtected() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func ExampleNewServer() {
	config := Config{JWTKey: "secretKey"}
	h, err := NewServer(&config)
	if err != nil {
		log.Fatal(err)
	}
	httptest.NewServer(h)
	// Output:
}
