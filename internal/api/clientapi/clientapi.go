package clientapi

import (
	"context"
	"time"

	"github.com/Ord1nI/netStats/internal/converter"
	"github.com/Ord1nI/netStats/internal/storage/stat"
	"google.golang.org/grpc"
	pb "github.com/Ord1nI/netStats/internal/proto"
	"google.golang.org/grpc/credentials/insecure"
)

type clientAPI struct{
	pb.UsersClient
}

func (a *clientAPI) AddStats(stats []stat.Stat) error {
	gStats := &pb.Stats{}
	gStats.Stat = make([]*pb.Stat,len(stats))
	for i, v := range stats {
		gStats.Stat[i] = converter.ToGrpsStats(&v)
	}
	_, err := a.UsersClient.AddStats(context.Background(), gStats)

	return err
}

func (a *clientAPI) GetStats(time time.Time) ([]stat.Stat, error){
	return nil, nil
}


func New(host string) (*clientAPI, error) {

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	a := &clientAPI{pb.NewUsersClient(conn)}

	return a, nil
}
