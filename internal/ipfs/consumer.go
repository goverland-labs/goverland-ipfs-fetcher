package ipfs

import (
	"context"
	"fmt"

	pevents "github.com/goverland-labs/goverland-platform-events/events/ipfs"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-ipfs-fetcher/internal/config"
)

const (
	groupName                = "ipfs"
	maxPendingAckPerConsumer = 10
)

type Consumer struct {
	nc      *nats.Conn
	service *Service

	msgCreatedClient *client.Consumer[pevents.MessagePayload]
}

func NewConsumer(nc *nats.Conn, service *Service) *Consumer {
	return &Consumer{
		nc:      nc,
		service: service,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	cc, err := client.NewConsumer(ctx, c.nc, group, pevents.SubjectMessageCreated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectMessageCreated, err)
	}

	c.msgCreatedClient = cc

	log.Info().Msg("ipfs consumer is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *Consumer) stop() error {
	if c.msgCreatedClient != nil {
		if err := c.msgCreatedClient.Close(); err != nil {
			return fmt.Errorf("close message created consumer: %w", err)
		}
	}

	return nil
}

func (c *Consumer) handler() pevents.MessageHandler {
	return func(payload pevents.MessagePayload) error {
		err := c.service.Process(context.TODO(), payload.IpfsID, payload.Type)
		if err != nil {
			log.Error().Err(err).Msg("process message")

			return err
		}

		return nil
	}
}
