package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseDomain = "https://4everland.io/ipfs/"
)

type Http struct {
}

// NewFetcher creates new fetcher
func NewFetcher() *Http {
	return &Http{}
}

func (f *Http) Fetch(_ context.Context, ipfsID string) (json.RawMessage, error) {
	req, err := http.NewRequest("GET", buildIpfsUrl(ipfsID), nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ipfs: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response: %w", err)
	}

	isSuccessStatus := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !isSuccessStatus {
		return nil, fmt.Errorf("http status code: %d, body: %s", resp.StatusCode, body)
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
