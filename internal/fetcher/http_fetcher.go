package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/goverland-labs/goverland-ipfs-fetcher/internal/config"
)

const (
	baseDomain = "https://4everland.io/ipfs/"
)

type Http struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	reqTimeout time.Duration
}

// NewFetcher creates new fetcher
func NewFetcher(httpConf config.HttpClient) *Http {
	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(httpConf.RatePerMinute)), 1)
	return &Http{
		httpClient: &http.Client{
			Timeout: httpConf.TimeoutMS,
		},
		limiter:    limiter,
		reqTimeout: httpConf.TimeoutMS,
	}
}

func (f *Http) Fetch(ctx context.Context, ipfsID string) (json.RawMessage, error) {
	if err := f.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, f.reqTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", buildIpfsUrl(ipfsID), nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ipfs %s: %w", ipfsID, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response: %w", err)
	}

	isSuccessStatus := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !isSuccessStatus {
		return nil, fmt.Errorf("ipfs: %s, http status code: %d, body: %s", ipfsID, resp.StatusCode, body)
	}

	isJsonValid := json.Valid(body)
	if !isJsonValid {
		return nil, fmt.Errorf("invalid json: %s", body)
	}

	return body, nil
}

func buildIpfsUrl(ipfsID string) string {
	return fmt.Sprintf("%s%s", baseDomain, ipfsID)
}
