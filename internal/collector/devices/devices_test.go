package devices

import (
	"errors"
	"testing"
	"time"

	"github.com/Ord1nI/netStats/internal/storage/stat"
	"github.com/scrapli/scrapligo/response"
	"github.com/scrapli/scrapligo/util"
	"github.com/sirikothe/gotextfsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var fsm string = `Value Type (\S*)

Start
 ^.*;;;\s${Comment} -> Continue
`

type mockDriver struct{}

func (*mockDriver) SendCommand(command string, opts ...util.Option) (*response.Response, error) {
	return nil, errors.New("mock")
}
func (*mockDriver) Open() error  { return nil }
func (*mockDriver) Close() error { return nil }

func TestDev(t *testing.T) {
	dev, err := NewDev(
		zap.NewNop().Sugar(),
		"127.0.0.1",
		[]Command{
			Command{
				Command:   "command",
				Fsm:       fsm,
				ParseFunc: func(parser *gotextfsm.ParserOutput, stats *stat.Stat) error { return nil },
			},
			Command{
				Command:   "command",
				Fsm:       fsm,
				ParseFunc: func(parser *gotextfsm.ParserOutput, stats *stat.Stat) error { return nil },
			},
		},
	)

	dev.Driver = &mockDriver{}

	require.NoError(t, err)
	assert.Equal(t, dev.Host, "127.0.0.1")
	assert.Equal(t, dev.Commands[0].Command, "command")
	assert.Equal(t, dev.Commands[1].Command, "command")

	tn := time.Now()

	res, err := dev.BackOff("123")

	te := time.Since(tn)

	var respnil *response.Response = nil

	assert.Greater(t, te, 8*time.Second)
	assert.Equal(t, res, respnil)
	assert.Equal(t, err, errors.New("mock"))
}
