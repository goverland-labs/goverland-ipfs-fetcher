//go:generate mockgen -destination=internal/dao/mocks_test.go -package=dao github.com/goverland-labs/goverland-ipfs-fetcher/internal/dao DataProvider,Publisher,DaoIDProvider

package main
