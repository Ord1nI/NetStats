package collector

import (
	"sync"
	"time"

	"github.com/Ord1nI/netStats/internal/collector/devices"
	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage"
	"golang.org/x/sync/errgroup"
)

type Collector struct {
	logger logger.Logger
	Schedule time.Duration
	Devices []devices.Device
	stop chan struct{}
	overAllStats chan []storage.Stats
}

func NewCollector(logger logger.Logger, schedule time.Duration, devices ...devices.Device) *Collector {
	collector := &Collector{logger:logger,Schedule:schedule, Devices:devices}
	return collector
}

func (c *Collector) Ping() error {
	var g errgroup.Group

	for _, i := range c.Devices {
		g.Go(func() error {
			return i.Ping()
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (c *Collector) Start() error {
	// if err := c.Ping(); err != nil {
	// 	return err
	// }

	c.stop = make(chan struct{})

	c.overAllStats = make(chan []storage.Stats)

	ticker := time.NewTicker(c.Schedule)

	go func() {
		defer close(c.overAllStats)
		for {
			select {
			case <-c.stop:
				c.logger.Infoln("Stop Collector")
				return
			case <-ticker.C:

				c.startWorkers(c.devPool())

				overAllStats := make([]storage.Stats,len(c.Devices))

				for v, i := range c.Devices {
					overAllStats[v] = *i.GetStats()
				}

				c.overAllStats<-overAllStats
			}

		}
	}()

	return nil
}

func (c *Collector) GetStats() []storage.Stats{
	return <-c.overAllStats
}

func (c *Collector) Stop() {
	close(c.stop)
}

func (c *Collector) devPool() <-chan devices.Device{
	devPool := make(chan devices.Device)
	go func() {
		defer close(devPool)
			for _, i := range c.Devices {
				devPool <- i
			}
	}()
	return devPool
}

func (c *Collector) startWorkers(devicePoll <-chan devices.Device) {

	var wg sync.WaitGroup

	for i := range 2 {
		wg.Add(1)
		c.logger.Infoln("start", i, "worker")

		go func() {
			defer wg.Done()
			for j := range devicePoll {

				err := j.CollectMetric()

				if err != nil {
					c.logger.Errorln(err)
				}

				c.logger.Infoln("End Worker")
			}
		}()
	}
	wg.Wait()
}
