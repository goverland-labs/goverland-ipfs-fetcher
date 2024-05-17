package message

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-ipfs-fetcher/protocol/ipfsfetcherpb"
)

type Server struct {
	ipfsfetcherpb.UnimplementedMessageServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetByID(_ context.Context, _ *ipfsfetcherpb.GetByIDRequest) (*ipfsfetcherpb.GetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByID not implemented")
}
