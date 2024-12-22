package grpcserv

import (
	"context"
	"testing"

	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandlers(t *testing.T) {
	serv, err := New(zap.NewNop().Sugar())

	require.NoError(t, err)

	storMock := &storageMock{}

	serv.Stor = storMock

	pbReqData := pb.Stats{
		Stat: []*pb.Stat{
			&pb.Stat{
				DevInfo: &pb.DevInfo{
					Uptime:           123,
					DevType:          "router",
					MemoryTotalBytes: 1,
					MemoryUsedBytes:  2,
				},
			},
		},
	}

	serv.AddStats(context.Background(), &pbReqData)

	assert.Equal(t, storMock.Data[0].DevInfo.Uptime, pbReqData.Stat[0].DevInfo.Uptime)
	assert.Equal(t, storMock.Data[0].DevInfo.DevType, pbReqData.Stat[0].DevInfo.DevType)
	assert.Equal(t, storMock.Data[0].DevInfo.MemoryTotalBytes, pbReqData.Stat[0].DevInfo.MemoryTotalBytes)
	assert.Equal(t, storMock.Data[0].DevInfo.MemoryUsedBytes, pbReqData.Stat[0].DevInfo.MemoryUsedBytes)
}
