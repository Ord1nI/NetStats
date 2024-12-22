package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ord1nI/netStats/internal/logger"
	grpcserv "github.com/Ord1nI/netStats/internal/server/grpsserv"
)

func main() {
	l, err := logger.New()
	if err != nil {
		fmt.Println("GG")
		return
	}

	serv, err := grpcserv.New(l)

	if err != nil {
		l.Fatal(err)
	}

	err = serv.Run()

	if err != nil {
		l.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	fmt.Println("End program")
	serv.Stop()

}
