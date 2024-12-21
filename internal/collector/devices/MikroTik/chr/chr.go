package chr

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ord1nI/netStats/internal/collector/devices"
	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/sirikothe/gotextfsm"
)

type chr struct {
	*devices.Dev
}

const versionCmd = "system/resource/print without-paging"
const interfaceAboutCmd = "interface/print detail proplist=name,default-name,type,mtu,mac,disabled,running without-paging"
const counterCmd = "interface/print stats-detail proplist=rx-byte,tx-byte,rx-packet,tx-packet,rx-drop,tx-drop,rx-error,tx-error without-paging"

func New(logger logger.Logger, host string, port int, username string, password string) (devices.Device, error) {
	logger.Infoln("Creating Mikrotik CHR device")

	chr := &chr{}

	dev, err := devices.NewDev(
		logger,
		host,

		[]devices.Command{
			devices.Command{versionCmd, versionTemplate, chr.parseVersion},
			devices.Command{interfaceAboutCmd, interfaceAboutTemplate, chr.parseInterface},
			devices.Command{counterCmd, InterfaceCounterTemplate, chr.parseCounter},
		},

		options.WithTransportType("standard"),
		options.WithDefaultLogger(),
		options.WithReadDelay(time.Second * 3),
		options.WithReturnChar("\r"),
		options.WithPromptPattern(regexp.MustCompile(`\[\S*\]\s*>\s`)),
		options.WithPasswordPattern(regexp.MustCompile(`.*Password:\s?$`)),
		options.WithUsernamePattern(regexp.MustCompile(`.*Username:\s?$`)),
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(username),
		options.WithAuthPassword(password),
		options.WithPort(port),
	)

	if err != nil {
		return nil, err
	}
	chr.Dev = dev

	logger.Infoln("End creating Mikrotik CHR device")
	return chr, nil
}


func (c *chr) parseVersion(parser *gotextfsm.ParserOutput, stats *storage.Stat) error {
	uptime, err := time.ParseDuration(parser.Dict[0]["Uptime"].(string))
	if err != nil {
		c.Logger.Errorln("Error when parsing time")
		uptime = time.Duration(0)
	}

	freememory, err := strconv.ParseFloat(parser.Dict[0]["FreeMemory"].(string), 64)
	if err != nil {
		c.Logger.Errorln("Error when parsing freememory")
		freememory = 0
	}
	freememory *= 1024*1024*8

	totalmemory, err := strconv.ParseFloat(parser.Dict[0]["TotalMemory"].(string), 64)
	if err != nil {
		c.Logger.Errorln("Error when parsing totalmemory")
		totalmemory = 0
	}
	totalmemory *= 1024*1024*8


	stats.DevInfo.Version = parser.Dict[0]["Version"].(string)
	stats.DevInfo.Processor = parser.Dict[0]["Cpu"].(string)
	stats.DevInfo.DevType = "router"
	stats.DevInfo.Hostname = parser.Dict[0]["BoardName"].(string)
	stats.DevInfo.MemoryTotalBytes = int64(totalmemory)
	stats.DevInfo.MemoryUsedBytes =  int64(totalmemory-freememory)
	stats.DevInfo.Uptime = int64(uptime.Seconds())

	return nil
}

func (c *chr) parseInterface(parser *gotextfsm.ParserOutput, stat *storage.Stat) error {
	if len(parser.Dict) != len(stat.InterfacesInfo) {
		stat.InterfacesInfo = make([]storage.L2Interface, len(parser.Dict))
	}

	for i, v := range parser.Dict {
		mtu, err := strconv.Atoi(v["MTU"].(string))
		if err != nil {
			mtu = -1
		}

		if (v["Comment"] != nil) {
			stat.InterfacesInfo[i].Description = v["Comment"].(string)
		}

		stat.InterfacesInfo[i].MAC = v["MAC"].(string)
		stat.InterfacesInfo[i].Name = v["Name"].(string)
		stat.InterfacesInfo[i].Running = v["Running"].(string)
		stat.InterfacesInfo[i].Disabled = v["Disabled"].(string)
		stat.InterfacesInfo[i].NameOriginal = v["NameOriginal"].(string)
		stat.InterfacesInfo[i].Ifname = v["Type"].(string)
		stat.InterfacesInfo[i].MTU = int32(mtu)
	}
	return nil
}

func (c *chr) parseCounter(parser *gotextfsm.ParserOutput, stats *storage.Stat) error {
	if len(parser.Dict) != len(stats.InterfacesInfo) {
		stats.InterfacesInfo = make([]storage.L2Interface, len(parser.Dict))
	}

	r := strings.NewReplacer(" ", "")

	for i, v := range parser.Dict {
		InBytes, err := strconv.ParseInt(r.Replace(v["InBytes"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing InBytes")
		}
		OutBytes, err := strconv.ParseInt(r.Replace(v["OutBytes"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing outBytes")
		}
		InPkts, err := strconv.ParseInt(r.Replace(v["InPkts"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing InPkts")
		}
		OutPkts, err := strconv.ParseInt(r.Replace(v["OutPkts"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing OutPkts")
		}
		InDrops, err := strconv.ParseInt(r.Replace(v["InDrops"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing InDrpos")
		}
		OutDrops, err := strconv.ParseInt(r.Replace(v["OutDrops"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing OutDrops")
		}
		ReadError, err := strconv.ParseInt(r.Replace(v["ReadError"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing ReadError")
		}
		OutError, err := strconv.ParseInt(r.Replace(v["OutError"].(string)), 10, 64)
		if (err != nil) {
			c.Logger.Errorln("Error when parsing OutError")
		}

		stats.InterfacesInfo[i].Counters.InBytes = InBytes
		stats.InterfacesInfo[i].Counters.OutBytes = OutBytes
		stats.InterfacesInfo[i].Counters.InPkts = InPkts
		stats.InterfacesInfo[i].Counters.OutPkts = OutPkts
		stats.InterfacesInfo[i].Counters.InDrops = InDrops
		stats.InterfacesInfo[i].Counters.OutDrops = OutDrops
		stats.InterfacesInfo[i].Counters.InErr = ReadError
		stats.InterfacesInfo[i].Counters.OutErr = OutError
	}
	return nil
}


func (c *chr) CollectMetric() error {
	c.Logger.Infoln("Start Metric collecting for MikroTik CHR host: ", c.Host)
	err := c.Open()

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	stats := &storage.Stat{}
	parser := gotextfsm.ParserOutput{}

	tc, stop := context.WithTimeout(context.Background(), time.Second*60)
	defer stop()

	_, err = c.Driver.Channel.ReadUntilPrompt(tc)

	if err != nil {
		log.Fatal("196")
	}

	for _, v := range c.Commands {
		out, err := c.BackOff(v.Command)

		if err != nil {
			log.Fatal("203")
		}


		err = parser.ParseTextString(out.Result, *v.Fsm, true)
		if err != nil {
			log.Fatal("209")
		}

		err = v.ParseFunc(&parser,stats)
		if err != nil {
			log.Fatal("214")
		}

		parser.Dict = nil
	}

	c.Stats = stats
	c.Logger.Infoln("End Metric collecting for MikroTik CHR host: ", c.Driver.Transport.Args.Host)
	return nil
}
