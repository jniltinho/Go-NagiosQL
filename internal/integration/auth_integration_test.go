//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestIntegration_Login_Success(t *testing.T) {
	srv, _, _ := newTestServer(t)

	code, body := authPost(t, srv, "", "/api/v1/auth/login",
		`{"username":"admin","password":"admin123"}`)
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", code, body)
	}
	if body["access_token"] == nil {
		t.Error("expected access_token in response")
	}
}

func TestIntegration_Login_WrongPassword(t *testing.T) {
	srv, _, _ := newTestServer(t)

	code, _ := authPost(t, srv, "", "/api/v1/auth/login",
		`{"username":"admin","password":"wrongpassword"}`)
	if code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", code)
	}
}

func TestIntegration_Login_UnknownUser(t *testing.T) {
	srv, _, _ := newTestServer(t)

	code, _ := authPost(t, srv, "", "/api/v1/auth/login",
		`{"username":"nobody","password":"x"}`)
	if code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", code)
	}
}

func TestIntegration_Logout(t *testing.T) {
	srv, _, _ := newTestServer(t)
	token := login(t, srv, "admin", "admin123")

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("logout: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestIntegration_RefreshToken(t *testing.T) {
	srv, _, _ := newTestServer(t)

	// Login to get the refresh_token httpOnly cookie.
	jar := newSimpleCookieJar()
	client := &http.Client{Jar: jar}

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/login",
		strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	var loginBody map[string]any
	json.NewDecoder(resp.Body).Decode(&loginBody)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: %d %v", resp.StatusCode, loginBody)
	}

	// Refresh — the jar sends the refresh_token cookie automatically.
	req2, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/auth/refresh", nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 on refresh, got %d", resp2.StatusCode)
	}
}

func TestIntegration_Protected_NoToken(t *testing.T) {
	srv, _, _ := newTestServer(t)

	resp, err := http.Get(srv.URL + "/api/v1/hosts")
	if err != nil {
		t.Fatalf("GET /hosts: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 without token, got %d", resp.StatusCode)
	}
}

// simpleCookieJar stores all cookies without domain matching logic.
type simpleCookieJar struct {
	cookies []*http.Cookie
}

func newSimpleCookieJar() *simpleCookieJar { return &simpleCookieJar{} }

func (j *simpleCookieJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	j.cookies = append(j.cookies, cookies...)
}

func (j *simpleCookieJar) Cookies(_ *url.URL) []*http.Cookie {
	return j.cookies
}
