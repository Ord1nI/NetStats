package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
)
func ciscoOnOpen(d *generic.Driver) error {
	if _, err := d.SendCommand("terminal length 0"); err != nil {
		return err
	}
	if _, err := d.SendCommand("enable"); err != nil {
		return err
	}
	return nil
}

func main() {
    d, err := generic.NewDriver(
		"192.168.1.2",
		// options.WithTimeoutSocket(5 * time.Second),
		options.WithTransportType("standard"),
		options.WithDefaultLogger(),
		options.WithReadDelay(time.Second * 3),
		options.WithReturnChar("\r\n"),
		options.WithPromptPattern(regexp.MustCompile(`\[\S*\]\s*>\s`)),
		options.WithPasswordPattern(regexp.MustCompile(`.*Password:\s?$`)),
		options.WithUsernamePattern(regexp.MustCompile(`.*Username:\s?$`)),
		options.WithAuthNoStrictKey(),
		options.WithTimeoutOps(20*time.Second),
		options.WithAuthUsername("admin"),
		// options.WithOnOpen(ciscoOnOpen),
		options.WithAuthPassword("123"),
		options.WithPort(22),
	)

	if err != nil {
		log.Fatal("FUCK")
	}

	d.Open()

	res, err := d.GetPrompt()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)


	d.Close()

	err = d.Open()

	for {
		b, err := d.Channel.Read()
		if err != nil {
			log.Fatal(err)
		}
		if b == nil {
			fmt.Println("nil")
		}
		fmt.Println(string(b))
	}
}
