package chr

import (
	"strconv"
	"testing"

	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockChr struct{}

func (*mockChr) SendCommand(command string, opts ...util.Option) (*response.Response, error) {
	res := response.NewResponse(versionCmd, "host", 22, []string{})
	switch command {
	case versionCmd:
		res.Result = `                   uptime: 1m58s
                  version: 7.16 (stable)
               build-time: 2024-09-20 13:00:27
         factory-software: 7.1
              free-memory: 178.1MiB
             total-memory: 384.0MiB
                      cpu: QEMU
                cpu-count: 1
            cpu-frequency: 4192MHz
                 cpu-load: 0%
           free-hdd-space: 71.2MiB
          total-hdd-space: 89.2MiB
  write-sect-since-reboot: 296
         write-sect-total: 296
        architecture-name: x86_64
               board-name: CHR QEMU Standard PC (i440FX + PIIX, 1996)
                 platform: MikroTik`
	case interfaceAboutCmd:
		res.Result = `Flags: D - dynamic; X - disabled; I - inactive, R - running; S - slave; P - passthrough
 0   R   ;;; 123
         name="ether1" default-name="ether1" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:00 disabled=no
         running=yes

 1       name="ether2" default-name="ether2" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:01 disabled=no running=no

 2       name="ether3" default-name="ether3" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:02 disabled=no running=no

 3       name="ether4" default-name="ether4" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:03 disabled=no running=no

 4       name="ether5" default-name="ether5" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:04 disabled=no running=no

 5       name="ether6" default-name="ether6" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:05 disabled=no running=no

 6       name="ether7" default-name="ether7" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:06 disabled=no running=no

 7       ;;; 1234
         name="ether8" default-name="ether8" type="ether" mtu=1500 mac-address=0C:93:1B:22:00:07 disabled=no
         running=no

 8   R   name="lo" type="loopback" mtu=65536 mac-address=00:00:00:00:00:00 disabled=no running=yes`
	case counterCmd:
		res.Result = `Flags: D - dynamic; X - disabled; I - inactive, R - running; S - slave; P - passthrough
 0   R   ;;; 123
         rx-byte=54 447 tx-byte=74 828 rx-packet=501 tx-packet=413 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 1       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 2       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 3       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 4       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 5       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 6       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 7       rx-byte=0 tx-byte=0 rx-packet=0 tx-packet=0 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0

 8   R   rx-byte=1 296 tx-byte=1 296 rx-packet=8 tx-packet=8 rx-drop=0 tx-drop=0 rx-error=0 tx-error=0
`
	case hostNameCmd:
		res.Result = `name: MikroTik`
	}
	return res, nil
}

func (*mockChr) Open() error  { return nil }
func (*mockChr) Close() error { return nil }

func (c *chr) Open() error {
	return c.Driver.Open()
}

func TestChr(t *testing.T) {
	l := zap.NewNop().Sugar()

	dev, err := New(l, "10.10.5.2", 22, "user", "123")

	require.NoError(t, err)

	dev.Driver = &mockChr{}

	dev.CollectMetric()

	stats := dev.GetStats()

	assert.Equal(t, stats.DevInfo.Uptime, int64(118))
	assert.Equal(t, stats.DevInfo.Processor, "QEMU")
	assert.Equal(t, stats.DevInfo.Version, "7.16 (stable)")
	assert.Equal(t, stats.DevInfo.MemoryUsedBytes, int64(215_901_798))
	assert.Equal(t, stats.DevInfo.MemoryTotalBytes, int64(402_653_184))
	assert.Equal(t, stats.DevInfo.DevType, "router")
	assert.Equal(t, stats.DevInfo.Hostname, "MikroTik")

	assert.Equal(t, stats.InterfacesInfo[0].Running, "yes")
	assert.Equal(t, stats.InterfacesInfo[0].Name, "ether1")
	assert.Equal(t, stats.InterfacesInfo[0].NameOriginal, "ether1")
	assert.Equal(t, stats.InterfacesInfo[0].Ifname, "ether")
	assert.Equal(t, stats.InterfacesInfo[0].MTU, int32(1500))
	assert.Equal(t, stats.InterfacesInfo[0].MAC, "0C:93:1B:22:00:00")
	assert.Equal(t, stats.InterfacesInfo[0].Disabled, "no")
	assert.Equal(t, stats.InterfacesInfo[0].Description, "123")

	assert.Equal(t, stats.InterfacesInfo[0].Counters.OutPkts, int64(413))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.InPkts, int64(501))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.InBytes, int64(54447))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.OutBytes, int64(74828))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.OutErr, int64(0))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.InErr, int64(0))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.InDrops, int64(0))
	assert.Equal(t, stats.InterfacesInfo[0].Counters.OutDrops, int64(0))

	for i := range 7 {
		assert.Equal(t, stats.InterfacesInfo[i+1].Running, "no")
		assert.Equal(t, stats.InterfacesInfo[i+1].Name, ("ether" + strconv.Itoa(i+2)))
		assert.Equal(t, stats.InterfacesInfo[i+1].NameOriginal, "ether"+strconv.Itoa(i+2))
		assert.Equal(t, stats.InterfacesInfo[i+1].Ifname, "ether")
		assert.Equal(t, stats.InterfacesInfo[i+1].MTU, int32(1500))
		assert.Equal(t, stats.InterfacesInfo[i+1].MAC, "0C:93:1B:22:00:0"+strconv.Itoa(i+1))
		assert.Equal(t, stats.InterfacesInfo[i+1].Disabled, "no")

		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.OutPkts, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.InPkts, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.InBytes, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.OutBytes, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.OutErr, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.InErr, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.InDrops, int64(0))
		assert.Equal(t, stats.InterfacesInfo[i+1].Counters.OutDrops, int64(0))
	}

	assert.Equal(t, stats.InterfacesInfo[8].Running, "yes")
	assert.Equal(t, stats.InterfacesInfo[8].Name, "lo")
	assert.Equal(t, stats.InterfacesInfo[8].Ifname, "loopback")
	assert.Equal(t, stats.InterfacesInfo[8].MTU, int32(65536))
	assert.Equal(t, stats.InterfacesInfo[8].MAC, "00:00:00:00:00:00")
	assert.Equal(t, stats.InterfacesInfo[8].Disabled, "no")

	assert.Equal(t, stats.InterfacesInfo[8].Counters.OutPkts, int64(8))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.InPkts, int64(8))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.InBytes, int64(1296))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.OutBytes, int64(1296))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.OutErr, int64(0))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.InErr, int64(0))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.InDrops, int64(0))
	assert.Equal(t, stats.InterfacesInfo[8].Counters.OutDrops, int64(0))

}
