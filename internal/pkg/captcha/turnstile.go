package captcha

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
)

type TurnstileConfig struct {
	Secret        string
	SiteVerifyURL string
	HTTPClient    *http.Client
	ExpectedHost  string // optional: verify hostname
}

type Response struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
	Action      string   `json:"action,omitempty"`
	CData       string   `json:"cdata,omitempty"`
}

type Turnstile struct {
	cfg TurnstileConfig
}

// NewTurnstileFromEnv creates a Turnstile with values from environment.
// Use this in production (cmd/api).
func NewTurnstileFromEnv() *Turnstile {
	return NewTurnstile(TurnstileConfig{
		Secret:        os.Getenv("TURNSTILE_SECRET"),
		SiteVerifyURL: "https://challenges.cloudflare.com/turnstile/v0/siteverify",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ExpectedHost: os.Getenv("EXPECTED_HOSTNAME"), // optional
	})
}

// NewTurnstile allows full DI – perfect for tests.
func NewTurnstile(cfg TurnstileConfig) *Turnstile {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	if cfg.SiteVerifyURL == "" {
		cfg.SiteVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	}
	return &Turnstile{cfg: cfg}
}

func (t *Turnstile) Verify(ctx context.Context, token, remoteIP string) (bool, *errx.Error) {
	if token == "" {
		return false, errx.ErrCaptcha
	}

	data := url.Values{
		"secret":   {t.cfg.Secret},
		"response": {token},
	}
	if remoteIP != "" {
		if net.ParseIP(remoteIP) == nil {
			return false, errx.ErrCaptcha
		}
		data.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.cfg.SiteVerifyURL, nil)
	if err != nil {
		sentry.CaptureException(err)
		return false, errx.InternalError()
	}
	req.PostForm = data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.cfg.HTTPClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return false, errx.InternalError()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sentry.CaptureException(err)
		return false, errx.InternalError()
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1 MB max
	if err != nil {
		sentry.CaptureException(err)
		return false, errx.InternalError()
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		sentry.CaptureException(err)
		return false, errx.InternalError()
	}

	if !r.Success {
		return false, errx.ErrCaptcha
	}

	// Optional: verify hostname
	if t.cfg.ExpectedHost != "" && r.Hostname != t.cfg.ExpectedHost {
		return false, errx.ErrCaptcha
	}

	// Optional: verify timestamp is recent (prevent replay)
	if r.ChallengeTs != "" {
		ts, err := time.Parse(time.RFC3339, r.ChallengeTs)
		if err != nil {
			return false, errx.ErrCaptcha
		}
		if time.Since(ts) > 5*time.Minute {
			return false, errx.ErrCaptcha
		}
	}

	return true, nil
}
