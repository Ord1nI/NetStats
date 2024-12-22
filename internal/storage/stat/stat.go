package stat

type Stat struct {
	DevInfo        DeviceInfo
	InterfacesInfo []L2Interface
}

type DeviceInfo struct {
	Version          string `json:"version"`
	Hostname         string `json:"hostname"`
	Uptime           int64  `json:"uptime"`
	Processor        string `json:"processor"`
	DevType          string `json:"devtype"`
	MemoryTotalBytes int64  `json:"MemoryTotalBytes"`
	MemoryUsedBytes  int64  `json:"MemoryUsedBytes"`
}

type Counters struct {
	InDrops  int64 `json:"InDrops"`
	OutDrops int64 `json:"OutDrops"`
	InPkts   int64 `json:"InPkts"`
	InBytes  int64 `json:"InBytes"`
	InErr    int64 `json:"InErr"`
	OutPkts  int64 `json:"OutPkts"`
	OutBytes int64 `json:"OutBytes"`
	OutErr   int64 `json:"OutErr"`
}

type L2Interface struct {
	Name         string   `json:"Name"`
	NameOriginal string   `json:"NameOriginal"`
	Ifname       string   `json:"Ifname"`
	Disabled     string   `json:"Disabled"`
	Running      string   `json:"Running"`
	MAC          string   `json:"Mac"`
	MTU          int32    `json:"Mtu"`
	Description  string   `json:"Description"`
	Counters     Counters `json:"Counters"`
}
