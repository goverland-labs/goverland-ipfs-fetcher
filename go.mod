module github.com/goverland-labs/goverland-ipfs-fetcher

replace github.com/goverland-labs/goverland-ipfs-fetcher/protocol => ./protocol

go 1.22

require (
	github.com/caarlos0/env/v6 v6.10.1
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/goverland-labs/goverland-ipfs-fetcher/protocol v0.0.0
	github.com/goverland-labs/goverland-platform-events v0.2.3-0.20240518122348-b6815a9d8d7b
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/nats-io/nats.go v1.30.2
	github.com/prometheus/client_golang v1.18.0
	github.com/rs/zerolog v1.31.0
	github.com/s-larionov/process-manager v0.0.1
	github.com/shopspring/decimal v1.3.1
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
	gorm.io/datatypes v1.2.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.2
)

require (
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/nats-io/nkeys v0.4.5 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	gorm.io/driver/mysql v1.4.7 // indirect
)
