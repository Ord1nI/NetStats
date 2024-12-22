package grpcserv

import (
	"testing"
	"time"

	"github.com/Ord1nI/netStats/internal/api/clientapi"
	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/storage/stat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
)

type storageMock struct {
	Data []stat.Stat
}

func (s *storageMock) Add(data []stat.Stat, smth string) error {
	s.Data = data
	return nil
}

func (s *storageMock) Get(time.Time) ([]stat.Stat, error) {
	return []stat.Stat{
		stat.Stat{
			DevInfo: stat.DeviceInfo{
				Uptime:           123,
				DevType:          "router",
				MemoryTotalBytes: 1,
				MemoryUsedBytes:  2,
			},
		},
	}, nil
}

func grpcServTest(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	serv, err := New(zap.NewNop().Sugar())
	storMock := &storageMock{}
	serv.Stor = storMock

	require.NoError(t, err)

	pb.RegisterUsersServer(serv.Gserv, serv)

	go func() {
		err := serv.Gserv.Serve(lis)
		require.NoError(t, err)
	}()

	api, err := clientapi.New("")
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

	err = api.AddStats(reqData)

	require.NoError(t, err)

	assert.Equal(t, storMock.Data, reqData)

	stat, err := api.GetStats(time.Now())

	require.NoError(t, err)

	assert.Equal(t, reqData, stat)
}
