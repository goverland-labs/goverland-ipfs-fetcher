package config

import (
	"time"
)

type HttpClient struct {
	TimeoutMS     time.Duration `env:"HTTP_CLIENT_TIMEOUT_MS" envDefault:"2000ms"`
	RatePerMinute int           `env:"HTTP_CLIENT_RATE_PER_MINUTE" envDefault:"250"`
}
