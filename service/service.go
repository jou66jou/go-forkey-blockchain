package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jou66jou/go-forky-blockchain/service/handler"
)

func RunHTTP() error {
	mux := makeMuxRouter()
	httpAddr := "8080"
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handler.GetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handler.WriteBlock).Methods("POST")
	return muxRouter
}
