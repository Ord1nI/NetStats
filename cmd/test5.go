package main

import (
	"fmt"
	"time"

	"github.com/Ord1nI/netStats/internal/collector"
	"github.com/Ord1nI/netStats/internal/collector/devices/MikroTik/chr"
	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage/db"
)

func main() {

	db, err := db.NewDb()
	defer db.Close()
	l, _ := logger.New()

	if err != nil {
		l.Fatal(err)
	}


	dev1, err := chr.New(l, "192.168.1.2", 22, "admin", "123");
	dev2, err := chr.New(l, "192.168.1.2", 22, "admin", "123");

	if err != nil {
		l.Fatal(err)
	}


	c := collector.NewCollector(l, time.Second*150, dev1, dev2)

	c.Start()
	defer c.Stop()
	stats := c.GetStats()
	stats = c.GetStats()
	stats = c.GetStats()

	for _, i := range stats {
		fmt.Println(i)
	}
}
