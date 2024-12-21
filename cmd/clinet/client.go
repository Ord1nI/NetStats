package main

import (
	"fmt"
	"time"

	"github.com/Ord1nI/netStats/internal/client"
	"github.com/Ord1nI/netStats/internal/logger"
)

func main() {
	l,err := logger.New()

	if err != nil {
		fmt.Println(err)
	}

	c, err := client.NewClient(l)

	if err != nil {
		l.Fatal(err)
	}

	c.Start()
	defer c.Stop()

	time.Sleep(time.Hour * 1)
}
