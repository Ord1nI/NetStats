package converter

import (
	pb "github.com/Ord1nI/netStats/internal/proto"
	"github.com/Ord1nI/netStats/internal/storage"
)

func ToGrpsStats(stats *storage.Stat) *pb.Stat {
	gStat := &pb.Stat{}
	gStat.L2Interface = make([]*pb.L2Interface, 0, len(stats.InterfacesInfo))

	gStat.DevInfo.Version = stats.DevInfo.Version
	gStat.DevInfo.Hostname = stats.DevInfo.Hostname
	gStat.DevInfo.Uptime = stats.DevInfo.Uptime
	gStat.DevInfo.Processor = stats.DevInfo.Processor
	gStat.DevInfo.DevType = stats.DevInfo.DevType
	gStat.DevInfo.MemoryTotalBytes = stats.DevInfo.MemoryTotalBytes
	gStat.DevInfo.MemoryUsedBytes = stats.DevInfo.MemoryUsedBytes


	for i, v := range stats.InterfacesInfo {

		gStat.L2Interface[i] = &pb.L2Interface{}


		gStat.L2Interface[i].Name = v.Name
		gStat.L2Interface[i].NameOriginal = v.NameOriginal
		gStat.L2Interface[i].IfName = v.Ifname
		gStat.L2Interface[i].Disabled = v.Disabled
		gStat.L2Interface[i].Running = v.Running
		gStat.L2Interface[i].Mac = v.MAC
		gStat.L2Interface[i].Mtu = v.MTU
		gStat.L2Interface[i].Description = v.Description

		gStat.L2Interface[i].Counter = &pb.Counter{}

		gStat.L2Interface[i].Counter.InDrops = v.Counters.InDrops
		gStat.L2Interface[i].Counter.OutDrops = v.Counters.OutDrops
		gStat.L2Interface[i].Counter.InPkts = v.Counters.InPkts
		gStat.L2Interface[i].Counter.InBytes = v.Counters.InBytes
		gStat.L2Interface[i].Counter.InErr = v.Counters.InErr
		gStat.L2Interface[i].Counter.OutPkts = v.Counters.OutPkts
		gStat.L2Interface[i].Counter.OutBytes = v.Counters.OutBytes
		gStat.L2Interface[i].Counter.OutErr = v.Counters.OutErr

	}

	return gStat
}

func FromGrpsStats(stats *pb.Stat) *storage.Stat {
	gStat := &storage.Stat{}
	gStat.InterfacesInfo = make([]storage.L2Interface, 0, len(stats.L2Interface))

	gStat.DevInfo.Version = stats.DevInfo.Version
	gStat.DevInfo.Hostname = stats.DevInfo.Hostname
	gStat.DevInfo.Uptime = stats.DevInfo.Uptime
	gStat.DevInfo.Processor = stats.DevInfo.Processor
	gStat.DevInfo.DevType = stats.DevInfo.DevType
	gStat.DevInfo.MemoryTotalBytes = stats.DevInfo.MemoryTotalBytes
	gStat.DevInfo.MemoryUsedBytes = stats.DevInfo.MemoryUsedBytes

	for i, v := range stats.L2Interface {
		gStat.InterfacesInfo[i].Name = v.Name
		gStat.InterfacesInfo[i].NameOriginal = v.NameOriginal
		gStat.InterfacesInfo[i].Ifname = v.IfName
		gStat.InterfacesInfo[i].Disabled = v.Disabled
		gStat.InterfacesInfo[i].Running = v.Running
		gStat.InterfacesInfo[i].MAC = v.Mac
		gStat.InterfacesInfo[i].MTU = v.Mtu
		gStat.InterfacesInfo[i].Description = v.Description

		gStat.InterfacesInfo[i].Counters.InDrops = v.Counter.InDrops
		gStat.InterfacesInfo[i].Counters.OutDrops = v.Counter.OutDrops
		gStat.InterfacesInfo[i].Counters.InPkts = v.Counter.InPkts
		gStat.InterfacesInfo[i].Counters.InBytes = v.Counter.InBytes
		gStat.InterfacesInfo[i].Counters.InErr = v.Counter.InErr
		gStat.InterfacesInfo[i].Counters.OutPkts = v.Counter.OutPkts
		gStat.InterfacesInfo[i].Counters.OutBytes = v.Counter.OutBytes
		gStat.InterfacesInfo[i].Counters.OutErr = v.Counter.OutErr

	}

	return gStat
}
