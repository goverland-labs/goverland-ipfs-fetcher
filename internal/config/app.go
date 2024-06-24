package config

type App struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	Prometheus  Prometheus
	Health      Health
	Nats        Nats
	DB          DB
	InternalAPI InternalAPI

	IpfsConsumersCount int `env:"IPFS_CONSUMERS_COUNT" envDefault:"1"`
}
