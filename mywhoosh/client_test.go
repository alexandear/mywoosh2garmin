package mywhoosh

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFormattedDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"01:31:54.000", "1h31m54s"},
		{"10:00:05.000", "10h00m05s"},
		{"00:45:10.000", "45m10s"},
		{"00:45:10", "45m10s"},
		{"", ""},
		{"45:10", "45:10"},
	}
	for _, tt := range tests {
		a := Activity{RideDuration: tt.in}
		if got := a.FormattedDuration(); got != tt.want {
			t.Errorf("FormattedDuration(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestFormattedDistance(t *testing.T) {
	t.Parallel()

	if got := (&Activity{Distance: 42.56}).FormattedDistance(); got != "42.6 km" {
		t.Errorf("got %q, want %q", got, "42.6 km")
	}
	if got := (&Activity{}).FormattedDistance(); got != "" {
		t.Errorf("zero distance: got %q, want empty", got)
	}
}

func TestDisplayName(t *testing.T) {
	t.Parallel()

	if got := (&Activity{Name: "Alpe du Zwift"}).DisplayName(); got != "Alpe du Zwift" {
		t.Errorf("got %q", got)
	}
	if got := (&Activity{}).DisplayName(); got != "Activity" {
		t.Errorf("empty name: got %q, want %q", got, "Activity")
	}
}

func TestFormattedDate(t *testing.T) {
	t.Parallel()

	const layout = "Mon 02 Jan 2006 15:04"
	start := time.Date(2026, time.March, 4, 17, 30, 0, 0, time.UTC)
	date := time.Date(2025, time.December, 1, 8, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		activity Activity
		want     string
	}{
		{
			name:     "start datetime",
			activity: Activity{StartDatetime: "2026-03-04T17:30:00.000Z", Date: float64(date.Unix())},
			want:     time.Unix(start.Unix(), 0).Format(layout),
		},
		{
			name:     "invalid start datetime falls back to date",
			activity: Activity{StartDatetime: "yesterday", Date: float64(date.Unix())},
			want:     time.Unix(date.Unix(), 0).Format(layout),
		},
		{
			name:     "date only",
			activity: Activity{Date: float64(date.Unix())},
			want:     time.Unix(date.Unix(), 0).Format(layout),
		},
	}
	for _, tt := range tests {
		if got := tt.activity.FormattedDate(); got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestSaveTokenResume(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(t.TempDir(), "nested")
	c := NewClient(dir)
	c.token = "access"
	c.whooshID = "w42"
	if err := c.saveToken(); err != nil {
		t.Fatal(err)
	}

	resumed := NewClient(dir)
	if err := resumed.Resume(); err != nil {
		t.Fatal(err)
	}
	if resumed.token != "access" || resumed.whooshID != "w42" || !resumed.IsAuthenticated() {
		t.Errorf("resumed: got token %q whooshID %q", resumed.token, resumed.whooshID)
	}
}

func TestResumeErrors(t *testing.T) {
	t.Parallel()

	if err := NewClient("").Resume(); err == nil {
		t.Error("no token dir: got nil error")
	}
	if err := NewClient(t.TempDir()).Resume(); err == nil {
		t.Error("missing file: got nil error")
	}

	for name, content := range map[string]string{
		"broken JSON": "{broken",
		"empty token": `{"access_token":"","whoosh_id":"w42"}`,
	} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "mywhoosh_token.json"), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		c := NewClient(dir)
		if err := c.Resume(); err == nil {
			t.Errorf("%s: got nil error", name)
		}
		if c.IsAuthenticated() {
			t.Errorf("%s: client authenticated", name)
		}
	}
}

func TestSaveTokenNoDir(t *testing.T) {
	t.Parallel()

	c := NewClient("")
	c.token = "access"
	if err := c.saveToken(); err != nil {
		t.Errorf("got %v, want nil", err)
	}
}

func TestNotAuthenticated(t *testing.T) {
	t.Parallel()

	c := NewClient(t.TempDir())
	if c.IsAuthenticated() {
		t.Error("new client authenticated")
	}
	if _, err := c.GetActivities(t.Context()); err == nil {
		t.Error("GetActivities: got nil error")
	}
	if _, err := c.GetRecentActivities(t.Context(), 10); err == nil {
		t.Error("GetRecentActivities: got nil error")
	}
	if _, err := c.DownloadFitFile(t.Context(), "file"); err == nil {
		t.Error("DownloadFitFile: got nil error")
	}
}
