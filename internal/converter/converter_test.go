package converter

import (
	"testing"

	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/stretchr/testify/assert"

	"github.com/Ord1nI/netStats/internal/storage/stat"
)

func TestConverter(t *testing.T) {
	pbStat := &pb.Stat{
		DevInfo: &pb.DevInfo{
			Hostname:         "Host",
			Uptime:           24425234,
			DevType:          "router",
			Version:          "v1.07",
			Processor:        "amd",
			MemoryUsedBytes:  1234,
			MemoryTotalBytes: 4321,
		},
		L2Interface: []*pb.L2Interface{
			&pb.L2Interface{
				IfName:       "interface",
				Description:  "very good interface",
				NameOriginal: "if1",
				Name:         "if2",
				Mac:          "23:42:23:32",
				Mtu:          1111,
				Counter: &pb.Counter{
					OutBytes: 123,
					InDrops:  0,
					InErr:    1,
					InPkts:   1000000,
					OutErr:   89,
					InBytes:  8888,
					OutPkts:  666,
					OutDrops: 34,
				},
				Running:  "no",
				Disabled: "yes",
			},
		},
	}

	stata := stat.Stat{
		DevInfo: stat.DeviceInfo{
			Hostname:         "Host",
			Uptime:           24425234,
			DevType:          "router",
			Version:          "v1.07",
			Processor:        "amd",
			MemoryUsedBytes:  1234,
			MemoryTotalBytes: 4321,
		},
		InterfacesInfo: []stat.L2Interface{
			stat.L2Interface{
				Name:         "if2",
				Ifname:       "interface",
				Description:  "very good interface",
				NameOriginal: "if1",
				MAC:          "23:42:23:32",
				MTU:          1111,
				Counters: stat.Counters{
					OutBytes: 123,
					InDrops:  0,
					InErr:    1,
					InPkts:   1000000,
					OutErr:   89,
					InBytes:  8888,
					OutPkts:  666,
					OutDrops: 34,
				},
				Running:  "no",
				Disabled: "yes",
			},
		},
	}

	convertedGrpc := ToGrpcStats(&stata)
	assert.Equal(t, convertedGrpc, pbStat)

	convertedStat := FromGrpcStats(pbStat)
	assert.Equal(t, stata, *convertedStat)

}
