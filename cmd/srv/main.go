package main

import (
	"time"

	"github.com/khezen/bulklog/pkg/config"
	"github.com/khezen/bulklog/pkg/log"
	"github.com/khezen/bulklog/pkg/server"
)

const (
	maxTries    = 30
	retryPeriod = 5 * time.Second
)

var (
	quit  chan error
	cfg   *config.Config
	serv  *server.Server
	timer *time.Timer
	i     int
)

func main() {
	quit = make(chan error)
	var err error
	cfg, err = config.Get()
	if err != nil {
		panic(err)
	}
	for {
		serv, err = server.New(cfg, quit)
		if err != nil {
			if i < maxTries {
				log.Err().Println(err)
				timer = time.NewTimer(retryPeriod)
				<-timer.C
				i++
			} else {
				panic(err)
			}
		} else {
			break
		}
	}
	if err != nil {
		panic(err)
	}
	go serv.ListenAndServe()
	err = <-quit
	panic(err)
}
