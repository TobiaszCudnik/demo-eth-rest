package server

import (
	"context"
	"fmt"
	"github.com/TobiaszCudnik/infura-interview/internal/service"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	service *service.Service
	port    string
}

func New(port string, service *service.Service) *Server {
	return &Server{
		service: service,
		port:    port,
	}
}

func (s *Server) Start(ctx context.Context, ready chan<- bool) {
	// init
	router := httprouter.New()
	router.GET("/", s.heartbeat)
	router.GET("/v1/block/:number", s.getBlock)
	// TODO
	//router.GET("/v1/block/:number/tx/:index", getBlockTx)
	ready <- true
	log.Print(http.ListenAndServe(":"+s.port, router))
	<-ctx.Done()
	// teardown
}

// @Summary ETH block data
// @Description Returns an ETH block by it's hex number.
// @Accept */*
// @Param number query integer true "Block number (hex)"
// @Success 200 {object} Response "Block data"
// @Failure 400 "Request error"
// @Failure 500 "Server error"
// @Router /v1/block/:number [get]
func (s *Server) getBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// validate
	_, err := strconv.ParseInt(ps.ByName("number"), 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}
	payload, err := s.service.GetBlock(r.Context(), ps.ByName("number"))
	if err != nil || payload == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}
	fmt.Fprintf(w, *payload)
}

func (s *Server) heartbeat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
