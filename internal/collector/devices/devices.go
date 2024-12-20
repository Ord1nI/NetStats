package devices

import (
	"time"

	"github.com/Ord1nI/netStats/internal/storage"
	"github.com/Ord1nI/netStats/internal/logger"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
	"github.com/sirikothe/gotextfsm"
)

type Device interface{
	CollectMetric() (error)
	GetStats() (*storage.Stats)
	Ping() (error)
}

type command struct {
	Command string
	Fsm *gotextfsm.TextFSM
	ParseFunc func(parser *gotextfsm.ParserOutput, stats *storage.Stats) error
}

type Command struct {
	Command string
	Fsm string
	ParseFunc func(parser *gotextfsm.ParserOutput, stats *storage.Stats) error
}

type Dev struct {
	Logger logger.Logger
	Driver *generic.Driver
	Stats *storage.Stats
	Commands []command
}

func (d *Dev) GetStats() (*storage.Stats) {
	return d.Stats
}

func (d *Dev) BackOff(command string) (*response.Response, error) {
	var err error
	var out *response.Response

	for i := 0; i < 3; i++ {
		out, err = d.Driver.SendCommand(command)
		if  err == nil {
			return out, nil
		}
		time.Sleep(3)
	}

	return nil, err
}

func (d *Dev) Ping() (error) {
	err := d.Driver.Open()
	defer d.Driver.Close()

	return err
}

func NewDev(logger logger.Logger, host string, cmds []Command, opts ...util.Option) (*Dev, error) {
	d := &Dev{}
	d.Logger = logger

	d.Logger.Infoln("Start creating Generic Device")

	var err error

	d.Driver, err = generic.NewDriver(host, opts...)

	if err != nil {
		d.Logger.Errorln("Error while creating driver")
		return nil, err
	}

	d.Commands = make([]command,len(cmds))

	for i, v := range cmds {

		fsm := &gotextfsm.TextFSM{}
		err = fsm.ParseString(v.Fsm)

		if err != nil {
			d.Logger.Errorln("Error while creating fsm parser")
			return nil, err
		}

		d.Commands[i].Command = v.Command
		d.Commands[i].Fsm = fsm
		d.Commands[i].ParseFunc = v.ParseFunc
	}

	d.Logger.Infoln("End creating Generic device")

	return d, nil
}
