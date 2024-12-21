package api

import (
	"context"
	"time"

	"github.com/Ord1nI/netStats/internal/converter"
	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type api struct{
	pb.UsersClient
}

type API interface {
	AddStats(stats []storage.Stat) error
	GetStats(time time.Time) ([]storage.Stat, error)
}

func (a *api) AddStats(stats []storage.Stat) error {
	gStats := &pb.Stats{}
	gStats.Stat = make([]*pb.Stat,0,len(stats))
	for i, v := range stats {
		gStats.Stat[i] = converter.ToGrpsStats(&v)
	}
	_, err := a.UsersClient.AddStats(context.Background(), gStats)

	return err
}

func (a *api) GetStats(time time.Time) ([]storage.Stat, error){
	return nil, nil
}


func NewApi(host string) (API, error) {

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	a := &api{pb.NewUsersClient(conn)}

	return a, nil
}
