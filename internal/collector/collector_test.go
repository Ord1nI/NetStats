package collector

import (
	"testing"
	"time"

	"github.com/Ord1nI/netStats/internal/storage/stat"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type DevMock struct{}

func (*DevMock) CollectMetric() error {
	time.Sleep(time.Second * 5)
	return nil
}

func (*DevMock) GetStats() *stat.Stat {
	return &stat.Stat{}
}

func TestCollector(t *testing.T) {
	l := zap.NewNop().Sugar()

	collector := NewCollector(l, time.Second*5, 3, &DevMock{}, &DevMock{})

	collector.Start()

	statsCh := collector.GetStatsCh()

	tn := time.Now()
	<-statsCh
	te := time.Since(tn)

	assert.Less(t, te, 6*time.Second)

	<-statsCh
	te = time.Since(tn)

	assert.Less(t, te, 16*time.Second)
}
