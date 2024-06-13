package message

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/goverland-labs/goverland-ipfs-fetcher/internal/ipfs"
	"github.com/goverland-labs/goverland-ipfs-fetcher/protocol/ipfsfetcherpb"
)

type ipfsService interface {
	GetByID(_ context.Context, id string) (*ipfs.Data, error)
}

type Server struct {
	ipfsfetcherpb.UnimplementedMessageServer

	ipfsService ipfsService
}

func NewServer(ipfsService ipfsService) *Server {
	return &Server{
		ipfsService: ipfsService,
	}
}

func (s *Server) GetByID(ctx context.Context, req *ipfsfetcherpb.GetByIDRequest) (*ipfsfetcherpb.GetByIDResponse, error) {
	data, err := s.ipfsService.GetByID(ctx, req.IpfsId)
	if err != nil {
		log.Error().Err(err).Msg("failed to get data")

		return nil, status.Errorf(codes.Internal, "failed to get data: %v", err)
	}

	if data == nil {
		log.Error().Msgf("data with id %s not found", req.IpfsId)

		return nil, status.Errorf(codes.NotFound, "data with id %s not found", req.IpfsId)
	}

	return &ipfsfetcherpb.GetByIDResponse{
		RawMessage: &anypb.Any{
			Value: data.Data,
		},
	}, nil
}
