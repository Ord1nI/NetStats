package client

import (
	"time"

	"github.com/Ord1nI/netStats/internal/api"
	"github.com/Ord1nI/netStats/internal/collector"
	"github.com/Ord1nI/netStats/internal/collector/devices/MikroTik/chr"
	"github.com/Ord1nI/netStats/internal/logger"
)

type Client struct {
	api api.API
	config *config
	logger logger.Logger
	collector *collector.Collector

	stopCh chan struct{}
}

func NewClient(l logger.Logger) (*Client, error) {
	client := &Client{logger:l}


	err := client.getConf()
	if err != nil {
		return nil, err
	}

	collector := collector.NewCollector(l, time.Duration(client.config.Schedule), client.config.RateLimit)

	for _, v := range client.config.DevParams {
		switch v.DeviceName {
		case "MicroTik":

			dev, err := chr.New(l, v.Host, v.Port, v.Username,v.Password)
			if err != nil {
				l.Errorln("Probably wrong data for device ", v.DeviceName," host: ", v.Host, " ",err)
			}
			collector.Devices = append(collector.Devices, dev)

		default:

			l.Errorln("Unsupported device ",v.DeviceName, " host: ", v.Host)
		}
	}

	client.collector = collector

	api, err := api.NewApi(client.config.Address)

	if err != nil {
		return nil, err
	}

	client.api = api

	return client, nil
}

func (c *Client) Start() error {
	err := c.collector.Start()

	if err != nil {
		c.logger.Errorln("Fail to start collector")
		return err
	}

	c.stopCh = make(chan struct{})

	statsCh := c.collector.GetStatsCh()

	go func() {

		defer c.collector.Stop()

		for {

			select {

			case <-c.stopCh:
				c.logger.Infoln("Client stop")
				return

			case stats := <- statsCh:
				c.api.AddStats(stats)
			}
		}
	}()

	return nil
}

func (c *Client) Stop() {
	close(c.stopCh)
}
