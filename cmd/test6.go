package main

import (
	"encoding/json"
	"fmt"

	"github.com/Ord1nI/netStats/internal/client"
)

func main() {
	c := []client.DevConnectionParams {
		client.DevConnectionParams{
			DeviceName:"MicroTik",
			Host: "192.168.1.2",
			Port: 22,
			Username: "admin",
			Password: "123",
		},client.DevConnectionParams{
			DeviceName:"MicroTik",
			Host: "192.168.1.2",
			Port: 22,
			Username: "admin",
			Password: "123",
		},
	}

	b, _ := json.Marshal(c)
	fmt.Println(string(b))
}
