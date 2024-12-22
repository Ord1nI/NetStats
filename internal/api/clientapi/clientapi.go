package clientapi

import (
	"context"
	"time"

	"github.com/Ord1nI/netStats/internal/converter"
	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/storage/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type clientAPI struct {
	pb.UsersClient
}

func (a *clientAPI) AddStats(stats []stat.Stat) error {
	gStats := &pb.Stats{}
	gStats.Stat = make([]*pb.Stat, len(stats))
	for i, v := range stats {
		gStats.Stat[i] = converter.ToGrpcStats(&v)
	}
	_, err := a.UsersClient.AddStats(context.Background(), gStats)

	return err
}

func (a *clientAPI) GetStats(time time.Time) ([]stat.Stat, error) {

	req := &pb.GetStatsReq{
		Time: time.Unix(),
	}

	pbM, err := a.UsersClient.GetStats(context.Background(), req)
	if err != nil {
		return nil, err
	}

	sStats := make([]stat.Stat, len(pbM.Stats.Stat))

	for i, v := range pbM.Stats.Stat {
		sStats[i] = *converter.FromGrpcStats(v)
	}

	return sStats, nil

}

func New(host string) (*clientAPI, error) {

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	a := &clientAPI{pb.NewUsersClient(conn)}

	return a, nil
}
