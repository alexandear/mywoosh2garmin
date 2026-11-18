package garmin

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParseUploadResult(t *testing.T) {
	t.Parallel()

	if err := parseUploadResult(http.StatusConflict, nil); !errors.Is(err, ErrDuplicate) {
		t.Errorf("409: got %v, want ErrDuplicate", err)
	}

	body := []byte(`{"detailedImportResult":{"failures":[{"messages":["duplicate in another sense"]}]}}`)
	if err := parseUploadResult(http.StatusOK, body); err == nil || errors.Is(err, ErrDuplicate) {
		t.Errorf("failures: got %v, want a non-duplicate error", err)
	}

	if err := parseUploadResult(http.StatusCreated, []byte(`{"detailedImportResult":{"failures":[]}}`)); err != nil {
		t.Errorf("success: got %v, want nil", err)
	}

	if err := parseUploadResult(http.StatusOK, []byte("not json")); err != nil {
		t.Errorf("non-JSON success: got %v, want nil", err)
	}

	err := parseUploadResult(http.StatusInternalServerError, []byte(strings.Repeat("x", 1000)))
	if err == nil || errors.Is(err, ErrDuplicate) {
		t.Fatalf("500: got %v, want a non-duplicate error", err)
	}
	if !strings.HasSuffix(err.Error(), ": "+strings.Repeat("x", 300)) {
		t.Errorf("500: got %q, want the body truncated to 300 bytes", err)
	}
}

func TestResumeValidToken(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	oauth1 := &OAuth1Token{OAuthToken: "t1", Domain: "garmin.cn"}
	oauth2 := &OAuth2Token{AccessToken: "a2", ExpiresAt: time.Now().Add(time.Hour).Unix()}
	if err := SaveTokens(dir, oauth1, oauth2); err != nil {
		t.Fatal(err)
	}

	c := NewClient(dir)
	if err := c.Resume(t.Context()); err != nil {
		t.Fatal(err)
	}
	if c.Domain != "garmin.cn" || c.OAuth2.AccessToken != "a2" {
		t.Errorf("got domain %q access token %q, want garmin.cn a2", c.Domain, c.OAuth2.AccessToken)
	}
}

func TestResumeErrors(t *testing.T) {
	t.Parallel()

	if err := NewClient("").Resume(t.Context()); err == nil {
		t.Error("no token dir: got nil error")
	}
	if err := NewClient(t.TempDir()).Resume(t.Context()); err == nil {
		t.Error("no cached tokens: got nil error")
	}
}

func TestUploadFITNotAuthenticated(t *testing.T) {
	t.Parallel()

	if err := NewClient(t.TempDir()).UploadFIT(t.Context(), "activity.fit"); err == nil {
		t.Error("got nil error")
	}
}
