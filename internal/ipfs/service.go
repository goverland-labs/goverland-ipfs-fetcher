package ipfs

import (
	"context"
	"encoding/json"
	"fmt"

	pevents "github.com/goverland-labs/goverland-platform-events/events/ipfs"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
)

type fetcher interface {
	Fetch(ctx context.Context, ipfsID string) (json.RawMessage, error)
}

type publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type Service struct {
	repo *Repo

	fetcher   fetcher
	publisher publisher
}

func NewService(repo *Repo, fetcher fetcher, publisher publisher) *Service {
	return &Service{
		repo:      repo,
		fetcher:   fetcher,
		publisher: publisher,
	}
}

func (s *Service) Process(_ context.Context, ipfsID, ipfsType string) error {
	var err error

	data, err := s.GetByID(context.TODO(), ipfsID)
	if err != nil {
		return fmt.Errorf("get data by id: %w", err)
	}

	if data != nil {
		return nil
	}

	log.Info().Msgf("processings ipfs: %s", ipfsID)

	fetchedData, err := s.fetcher.Fetch(context.TODO(), ipfsID)
	if err != nil {
		return fmt.Errorf("fetch ipfs: %w", err)
	}

	if err := s.repo.Create(Data{
		IpfsID: ipfsID,
		Type:   ipfsType,
		Data:   datatypes.JSON(fetchedData),
	}); err != nil {
		return fmt.Errorf("create ipfs data: %w", err)
	}

	log.Info().Msgf("ipfs was processed: %s", ipfsID)

	err = s.publisher.PublishJSON(context.TODO(), pevents.SubjectMessageCollected, pevents.MessagePayload{
		IpfsID: ipfsID,
		Type:   ipfsType,
	})
	if err != nil {
		return fmt.Errorf("publish ipfs collected message: %w", err)
	}

	log.Info().Msgf("ipfs published collected message: %s", ipfsID)

	return nil
}

func (s *Service) GetByID(_ context.Context, id string) (*Data, error) {
	return s.repo.GetByID(id)
}
