package adblock

import (
	"bufio"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const updateAtHour = 3

var domainRegex = regexp.MustCompile(`^(?i)([a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)

type Adblock struct {
	domains map[string]struct{}
	client  http.Client
	logger  *slog.Logger
}

func New(logger *slog.Logger) *Adblock {
	return &Adblock{
		domains: make(map[string]struct{}),
		client: http.Client{
			Transport: &http.Transport{DisableKeepAlives: true},
		},
		logger: logger,
	}
}

func (a *Adblock) Start() error {
	const op = "adblock.Start"

	if err := a.load(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go a.scheduleUpdate()

	return nil
}

func (a *Adblock) load() error {
	const op = "adblock.load"

	for _, url := range sources {
		err := a.downloadAndAddDomains(url)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	a.logger.Info("Adblock", "loaded domains", len(a.domains))

	return nil
}

func (a *Adblock) scheduleUpdate() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), updateAtHour, 0, 0, 0, now.Location())

		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}

		sleepDuration := time.Until(next)

		time.Sleep(sleepDuration)

		if err := a.load(); err != nil {
			a.logger.Error("failed to reload adblock list", "err", err)
		}
	}
}

func (a *Adblock) IsMatch(host string) bool {
	_, exists := a.domains[host]

	return exists
}

func (a *Adblock) downloadAndAddDomains(url string) error {
	const op = "adblock.downloadAndAddDomains"

	resp, err := a.client.Get(url)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.addDomains(resp); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Adblock) addDomains(resp *http.Response) error {
	const op = "adblock.addDomains"
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			line = strings.TrimSpace(line)
			if domainRegex.MatchString(line) {
				a.domains[line] = struct{}{}
			}

			continue
		}

		domain := fields[1]
		domain = strings.TrimSpace(domain)
		if domainRegex.MatchString(domain) {
			a.domains[domain] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
