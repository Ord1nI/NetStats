package server

import "flag"

type Config struct {
	Address string
	DBdsn   string
}

func (s *Server) getConf() error {
	var (
		fAddress = flag.String("a", "127.0.0.1:8080", "enter IP format ip:port")
		fDBdsn   = flag.String("d", "./stats.db", "enter database dsn")
	)

	flag.Parse()

	s.Config = &Config{
		Address: *fAddress,
		DBdsn:   *fDBdsn,
	}

	return nil
}
