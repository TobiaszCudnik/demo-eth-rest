package main

import (
	"context"
	"github.com/TobiaszCudnik/infura-interview/internal/eth"
	rdb "github.com/TobiaszCudnik/infura-interview/internal/redis"
	"github.com/TobiaszCudnik/infura-interview/internal/server"
	"github.com/TobiaszCudnik/infura-interview/internal/service"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"log"
	"os"
)

type ready struct {
	redis   chan bool
	server  chan bool
	service chan bool
	eth     chan bool
}

func main() {
	log.Println("Stating...")

	// init
	osKill := make(chan os.Signal, 1)
	signal.Notify(osKill, syscall.SIGINT, syscall.SIGTERM)
	ctx, shutdown := context.WithCancel(context.Background())
	r := ready{
		redis:   make(chan bool, 1),
		server:  make(chan bool, 1),
		service: make(chan bool, 1),
		eth:     make(chan bool, 1),
	}
	var gracefulExit sync.WaitGroup
	v := reflect.ValueOf(r)
	gracefulExit.Add(v.NumField())

	// deps
	redis := rdb.New(os.Getenv("REDIS_ADDR"))
	ethC := eth.New(os.Getenv("ETH_ADDR"))
	go func() {
		redis.Start(ctx, r.redis)
		gracefulExit.Done()
	}()
	go func() {
		ethC.Start(ctx, r.eth)
		gracefulExit.Done()
	}()
	<-r.redis
	<-r.eth

	// service
	serv := service.New(redis, ethC)
	go func() {
		serv.Start(ctx, r.service)
		gracefulExit.Done()
	}()
	<-r.service

	// server
	s := server.New(os.Getenv("HTTP_PORT"), serv)
	go func() {
		s.Start(ctx, r.server)
		gracefulExit.Done()
	}()
	<-r.server
	log.Println("Ready")

	// exit
	<-osKill
	log.Println("Shutting down...")
	shutdown()
	gracefulExit.Wait()

	// end
	log.Println("Bye...")
}
