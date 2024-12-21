package main

import (
	"fmt"

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

	err = serv.Run()

	if err != nil {
		l.Fatal(err)
	}
}
