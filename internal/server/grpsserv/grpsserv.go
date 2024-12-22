package grpcserv

import (
	"context"
	"net"
	"time"

	"github.com/Ord1nI/netStats/internal/converter"
	"github.com/Ord1nI/netStats/internal/logger"
	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/server"
	"github.com/Ord1nI/netStats/internal/storage/stat"
	"google.golang.org/grpc"
)

type GrpcServ struct {
	*server.Server
	Gserv *grpc.Server
	pb.UnimplementedUsersServer
}

func New(l logger.Logger) (*GrpcServ, error) {

	mServ, err := server.New(l)
	if err != nil {
		return nil, err
	}

	return &GrpcServ{Server: mServ, Gserv: grpc.NewServer()}, nil

}

func (s *GrpcServ) Run() error {
	s.Logger.Infoln("Startin sever host: ", s.Config.Address)

	listen, err := net.Listen("tcp", s.Config.Address)
	if err != nil {
		return err
	}

	pb.RegisterUsersServer(s.Gserv, s)

	if err := s.Gserv.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *GrpcServ) AddStats(ctx context.Context, gStats *pb.Stats) (*pb.AddStatsRes, error) {
	stats := make([]stat.Stat, len(gStats.Stat))

	for i, v := range gStats.Stat {
		stats[i] = *converter.FromGrpcStats(v)
	}

	err := s.Stor.Add(stats, "123")

	if err != nil {
		return &pb.AddStatsRes{Error: "Error while adding metrics"}, err
	}

	return &pb.AddStatsRes{Error: "All good"}, err
}

func (s *GrpcServ) GetStats(ctx context.Context, req *pb.GetStatsReq) (*pb.GetStatsRes, error) {
	stats, err := s.Stor.Get(time.Unix(req.Time, 0))

	if err != nil {
		return nil, err
	}

	if len(stats) == 0 {
		return nil, err
	}

	res := make([]*pb.Stat, len(stats))

	for i, v := range stats {
		res[i] = converter.ToGrpcStats(&v)
	}
	return &pb.GetStatsRes{Stats: &pb.Stats{Stat: res}}, err
}
