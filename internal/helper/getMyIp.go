package helper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var errGetMyIPParse = errors.New("getmyip: invalid IP address in response")

func GetMyIp(ctx context.Context) (string, error) {
	const op = "GetMyIp"

	client := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
		Timeout:   10 * time.Second, //nolint: mnd
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://api.ipify.org", nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = resp.Body.Close() }()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if net.ParseIP(string(ip)) == nil {
		return "", errGetMyIPParse
	}

	return string(ip), nil
}
