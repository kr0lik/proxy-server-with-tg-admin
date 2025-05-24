package adblock_test

import (
	"log/slog"
	"os"
	"proxy-server-with-tg-admin/internal/infrastructure/adblock"
	"testing"
)

func TestAdblockLoad(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	a := adblock.New(logger)

	err := a.Load()
	if err != nil {
		t.Fatalf("Adblock.Load() error: %v", err)
	}

	expectedDomains := []string{
		"doubleclick.net",
		"adservice.google.com",
		"mclean.lato.cloud.360safe.com",
		"zy16eoat1w.com",
		"s.206ads.com",
		"api.huq.io",
		"marketplace-android-b235.hyprmx.com",
		"123.ywxww.net",
		"zenglobalenerji.com",
	}

	for _, domain := range expectedDomains {
		if !a.IsMatch(domain) {
			t.Errorf("expected domain %q not found in domains map", domain)
		}
	}
}
