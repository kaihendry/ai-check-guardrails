package banner

import (
	"bytes"
	"os"
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func TestBanner_RunReturnsNoFindings(t *testing.T) {
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Banner: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("banner module should produce no findings, got %v", findings)
	}
}

func TestDisplay_NoOutput100(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	Display(100, config.Config{Modules: config.ModuleToggles{Banner: true}, ScoreThreshold: 70})
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() != 0 {
		t.Errorf("score 100 should produce no banner output, got: %q", buf.String())
	}
}

func TestDisplay_AdvisoryBanner(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	Display(80, config.Config{Modules: config.ModuleToggles{Banner: true}, ScoreThreshold: 70})
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.Len() == 0 {
		t.Error("score 80 (above threshold but < 100) should produce advisory banner")
	}
}

func TestDisplay_WarningBannerWithURL(t *testing.T) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	Display(50, config.Config{
		Modules:        config.ModuleToggles{Banner: true},
		ScoreThreshold: 70,
		BannerURL:      "https://example.com/training",
	})
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	out := buf.String()
	if out == "" {
		t.Error("score below threshold should produce warning banner")
	}
	if !bytes.Contains([]byte(out), []byte("https://example.com/training")) {
		t.Errorf("banner should contain training URL, got: %q", out)
	}
}
