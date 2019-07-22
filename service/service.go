package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jou66jou/go-forky-blockchain/common"
	"github.com/jou66jou/go-forky-blockchain/service/handler"
)

func RunHTTP(httpAddr string) error {
	mux := makeMuxRouter()
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
	rName := common.RouteName
	muxRouter := mux.NewRouter()
	// 列出所有區塊
	muxRouter.HandleFunc("/blocks", handler.GetBlockchain).Methods("GET")
	// 列出所有節點
	muxRouter.HandleFunc(rName["getAllPeers"], handler.GetPeers).Methods("GET")
	// 寫入新區塊
	muxRouter.HandleFunc("/newblock", handler.WriteBlock).Methods("POST")
	// 建立新websocket連線
	muxRouter.HandleFunc(rName["newWS"], handler.NewWS)

	return muxRouter
}
