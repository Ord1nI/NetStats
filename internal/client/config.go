package client

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"
)

type DevConnectionParams struct {
	DeviceName string `json:"deviceName"`
	Host string `json:"host"`
	Port int `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type config struct {
	Address string
	Schedule time.Duration
	DevParams []DevConnectionParams
	RateLimit int
}

func getDevicesFromFile(fileName string) ([]DevConnectionParams, error) {

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644);

	if err != nil {
		return nil, err
	}

	var devs = []DevConnectionParams{}

	bFile, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bFile, &devs)

	if err != nil {
		return nil, err
	}

	return devs, nil
}

func (c *Client) getConf() error {

	var (
		fAddress = flag.String("a", "127.0.0.1:8080", "enter IP format ip:port")
		fRateLimit = flag.Int("l", 3, "enter Rate limit")
		fSchedule = flag.Int64("s", 1, "enter device polling period in minutes")
		devParamsFilePath = flag.String("d", "./devices.json", "enter path to file with devices")
	)

	flag.Parse()


	c.config = &config{
		Address: *fAddress,
		Schedule: time.Duration(*fSchedule) * time.Minute,
		RateLimit: *fRateLimit,
	}


	devs, err := getDevicesFromFile(*devParamsFilePath)

	if err != nil {
		c.logger.Errorln("Fail to get devices from file ", err)
		return err
	}

	c.config.DevParams = devs


	//TODO need to be Enabled
	// if c.config.Schedule < 30 * time.Minute {
	// 	c.logger.Info("Shedule les than 30 minutes is not allowed.\n Go back to 30")
	// 	c.config.Schedule = 30 * time.Minute
	// }

	return nil
}
