package server

import (
	"time"

	"github.com/Ord1nI/netStats/internal/logger"
	"github.com/Ord1nI/netStats/internal/storage/db"
	"github.com/Ord1nI/netStats/internal/storage/stat"
)

type storage interface {
	Add([]stat.Stat, string) error
	Get(time.Time) ([]stat.Stat, error)
}

type Server struct {
	Logger logger.Logger
	Config *Config
	Stor   storage
}

func New(logger logger.Logger) (*Server, error) {
	serv := &Server{}

	serv.Logger = logger

	err := serv.getConf()
	if err != nil {
		return nil, err
	}

	db, err := db.NewDb(serv.Config.DBdsn)
	if err != nil {
		return nil, err
	}
	serv.Stor = db

	return serv, nil
}
