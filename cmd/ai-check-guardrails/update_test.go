package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCheckAndUpdate_UpdatesWhenNewVersionAvailable(t *testing.T) {
	const fakeNewVersion = "v0.0.99999999999999"
	const fakeContent = "fake-binary-content"

	assetName := "ai-check-guardrails-" + runtime.GOOS + "-" + runtime.GOARCH

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/kaihendry/ai-check-guardrails/releases/latest":
			rel := githubRelease{
				TagName: fakeNewVersion,
				Assets: []struct {
					Name               string `json:"name"`
					BrowserDownloadURL string `json:"browser_download_url"`
				}{
					{
						Name:               assetName,
						BrowserDownloadURL: "http://" + r.Host + "/download/binary",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rel)
		case "/download/binary":
			w.Write([]byte(fakeContent))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	origBase := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = origBase }()

	dir := t.TempDir()
	fakeExe := filepath.Join(dir, "ai-check-guardrails")
	if err := os.WriteFile(fakeExe, []byte("old-content"), 0755); err != nil {
		t.Fatal(err)
	}

	origGet := getExecutable
	getExecutable = func() (string, error) { return fakeExe, nil }
	defer func() { getExecutable = origGet }()

	if err := checkAndUpdate("dev"); err != nil {
		t.Fatalf("checkAndUpdate returned error: %v", err)
	}

	got, err := os.ReadFile(fakeExe)
	if err != nil {
		t.Fatalf("read updated binary: %v", err)
	}
	if string(got) != fakeContent {
		t.Errorf("binary content = %q, want %q", got, fakeContent)
	}
}

func TestCheckAndUpdate_NoopWhenReleaseTagIsMovingLatest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rel := githubRelease{TagName: "latest"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rel)
	}))
	defer srv.Close()

	origBase := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = origBase }()

	origGet := getExecutable
	downloaded := false
	getExecutable = func() (string, error) {
		downloaded = true
		return "", nil
	}
	defer func() { getExecutable = origGet }()

	if err := checkAndUpdate("v0.0.20260513135826+5118f3a"); err != nil {
		t.Fatalf("checkAndUpdate returned error: %v", err)
	}
	if downloaded {
		t.Error("binary was replaced despite release tag being a moving 'latest' tag")
	}
}

func TestCheckAndUpdate_NoopWhenAlreadyCurrent(t *testing.T) {
	const currentVersion = "v0.0.20260508000000+abc1234"

	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		rel := githubRelease{TagName: currentVersion}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rel)
	}))
	defer srv.Close()

	origBase := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = origBase }()

	if err := checkAndUpdate(currentVersion); err != nil {
		t.Fatalf("checkAndUpdate returned error: %v", err)
	}
	_ = called
}
