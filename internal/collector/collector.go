package collector

import (
	"sync"
	"time"

	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage/stat"
)

type Device interface {
	CollectMetric() error
	GetStats() *stat.Stat
}

type Collector struct {
	logger       logger.Logger
	Schedule     time.Duration
	RateLimit    int
	Devices      []Device
	stop         chan struct{}
	overAllStats chan []stat.Stat
}

func NewCollector(logger logger.Logger, schedule time.Duration, RateLimit int, devices ...Device) *Collector {
	return &Collector{logger: logger, Schedule: schedule, Devices: devices, RateLimit: RateLimit}
}

func (c *Collector) Add(d Device) {
	c.Devices = append(c.Devices, d)
}

func (c *Collector) Start() error {

	c.stop = make(chan struct{})

	c.overAllStats = make(chan []stat.Stat)

	overAllStats := make([]stat.Stat, len(c.Devices))

	go func() {

		c.startWorkers(c.devPool(c.stop))

		for v, i := range c.Devices {
			overAllStats[v] = *i.GetStats()
		}

		c.overAllStats <- overAllStats

		ticker := time.NewTicker(c.Schedule)

		defer close(c.overAllStats)
		for {

			select {

			case <-c.stop:
				c.logger.Infoln("Stop Collector")
				return
			case <-ticker.C:

				c.startWorkers(c.devPool(c.stop))

				for v, i := range c.Devices {
					overAllStats[v] = *i.GetStats()
				}

				c.overAllStats <- overAllStats
			}

		}
	}()

	return nil
}

func (c *Collector) GetStatsCh() <-chan []stat.Stat {
	return c.overAllStats
}

func (c *Collector) Stop() {
	close(c.stop)
}

func (c *Collector) devPool(stop <-chan struct{}) <-chan Device {
	devPool := make(chan Device)

	go func() {

		defer close(devPool)

		for _, i := range c.Devices {

			select {
			case <-stop:

				c.logger.Infoln("devPool stop")
				return

			default:

				devPool <- i
			}
		}
	}()
	return devPool
}

func (c *Collector) startWorkers(devicePoll <-chan Device) {

	var wg sync.WaitGroup

	for i := range c.RateLimit {
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
