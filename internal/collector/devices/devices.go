package devices

import (
	"time"

	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage/stat"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
	"github.com/sirikothe/gotextfsm"
)

type command struct {
	Command   string
	Fsm       *gotextfsm.TextFSM
	ParseFunc func(parser *gotextfsm.ParserOutput, stats *stat.Stat) error
}

type Command struct {
	Command   string
	Fsm       string
	ParseFunc func(parser *gotextfsm.ParserOutput, stats *stat.Stat) error
}

type Driver interface {
	SendCommand(string, ...util.Option) (*response.Response, error)
	Open() error
	Close() error
}

type Dev struct {
	Logger logger.Logger

	Driver Driver

	Host    string
	Options []util.Option

	Stats    *stat.Stat
	Commands []command
}

func NewDev(logger logger.Logger, host string, cmds []Command, opts ...util.Option) (*Dev, error) {
	var err error

	d := &Dev{}
	d.Logger = logger

	d.Logger.Infoln("Start creating Generic Device")

	d.Host = host
	d.Options = opts

	d.Commands = make([]command, len(cmds))

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

func (d *Dev) Open() error {
	var err error

	d.Driver, err = generic.NewDriver(d.Host, d.Options...)

	if err != nil {
		d.Logger.Errorln("Error while creating driver")
		return err
	}

	return d.Driver.Open()
}

func (d *Dev) Close() error {
	return d.Driver.Close()
}

func (d *Dev) GetStats() *stat.Stat {
	return d.Stats
}

func (d *Dev) BackOff(command string) (*response.Response, error) {
	var err error
	var out *response.Response

	for i := 0; i < 3; i++ {
		out, err = d.Driver.SendCommand(command)
		if err == nil {
			return out, nil
		}
		time.Sleep(3 * time.Second)
	}

	return nil, err
}
