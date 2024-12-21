package main

import (
	"fmt"
	"time"

	"github.com/Ord1nI/netStats/internal/collector"
	"github.com/Ord1nI/netStats/internal/collector/devices/MikroTik/chr"
	"github.com/Ord1nI/netStats/internal/logger"
)

func main() {

	l, err := logger.New()

	if err != nil {
		l.Fatal(err)
	}


	dev1, err := chr.New(l, "192.168.1.2", 22, "admin", "123");
	dev2, err := chr.New(l, "192.168.1.3", 22, "admin", "123");

	if err != nil {
		l.Fatal(err)
	}

	c := collector.NewCollector(l, time.Second*20, dev1, dev2)

	c.Start()
	defer c.Stop()
	stats := c.GetStats()
	stats = c.GetStats()
	stats = c.GetStats()

	for _, i := range stats {
		fmt.Println(i)
	}
}
