package clientapi

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/storage/stat"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type userClientMock struct {
	Data *pb.Stats
}

func (u *userClientMock) AddStats(ctx context.Context, in *pb.Stats, opts ...grpc.CallOption) (*pb.AddStatsRes, error) {
	u.Data = in
	return nil, errors.New("123")
}

func (u *userClientMock) GetStats(ctx context.Context, in *pb.GetStatsReq, opts ...grpc.CallOption) (*pb.GetStatsRes, error) {
	return &pb.GetStatsRes{
		Stats: &pb.Stats{
			Stat: []*pb.Stat{
				&pb.Stat{
					DevInfo: &pb.DevInfo{
						Uptime:           123,
						DevType:          "router",
						MemoryTotalBytes: 1,
						MemoryUsedBytes:  2,
					},
					L2Interface: []*pb.L2Interface{},
				},
			},
		},
	}, nil
}

func TestAddStats(t *testing.T) {
	ucm := &userClientMock{}
	clientApi := clientAPI{ucm}

	reqData := []stat.Stat{
		stat.Stat{
			DevInfo: stat.DeviceInfo{
				Uptime:           123,
				DevType:          "router",
				MemoryTotalBytes: 1,
				MemoryUsedBytes:  2,
			},
		},
	}

	pbReqData := pb.Stats{
		Stat: []*pb.Stat{
			&pb.Stat{
				DevInfo: &pb.DevInfo{
					Uptime:           123,
					DevType:          "router",
					MemoryTotalBytes: 1,
					MemoryUsedBytes:  2,
				},
				L2Interface: []*pb.L2Interface{},
			},
		},
	}
	err := clientApi.AddStats(reqData)

	assert.Equal(t, err, errors.New("123"))
	assert.Equal(t, *ucm.Data, pbReqData)
}

func TestGetStat(t *testing.T) {
	ucm := &userClientMock{}
	clientApi := clientAPI{ucm}
	stats, err := clientApi.GetStats(time.Now())

	reqData := []stat.Stat{
		stat.Stat{
			DevInfo: stat.DeviceInfo{
				Uptime:           123,
				DevType:          "router",
				MemoryTotalBytes: 1,
				MemoryUsedBytes:  2,
			},
			InterfacesInfo: []stat.L2Interface{},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, stats, reqData)
}
