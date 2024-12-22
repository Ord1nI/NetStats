package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ord1nI/netStats/internal/client"
	"github.com/Ord1nI/netStats/internal/logger"
)

func main() {
	l, err := logger.New()

	if err != nil {
		fmt.Println(err)
	}

	c, err := client.NewClient(l)

	if err != nil {
		l.Fatal(err)
	}

	c.Start()

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	fmt.Println("End program")
	c.Stop()

}
