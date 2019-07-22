package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jou66jou/go-forky-blockchain/block"
	"github.com/jou66jou/go-forky-blockchain/p2p"
)

type Message struct {
	Wallet int `json:"wallet"`
}

// 新websocket請求
func NewWS(res http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	rPort, ok := q["port"]
	if !ok {
		fmt.Println("url value port is nil")
		http.NotFound(res, req)
		return
	}
	ip := strings.Split(req.RemoteAddr, ":")
	// 取得請求端ip:port
	taget := ip[0] + ":" + rPort[0]

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("new client error: " + err.Error())
		http.NotFound(res, req)
		return
	}

	// req帶有brdcst key則不進行廣播，brdcst代表req端是接收到廣播而發起websocket，避免廣播風暴
	v, ok := q["brdcst"]
	if !ok {
		if len(v) == 0 {
			// 廣播新結點
			p2p.BroadcastAddr(taget)
		}
	}

	// p2p
	newPeer := p2p.AppendNewPeer(conn, taget)
	go newPeer.Write()
	go newPeer.Read()
	p2p.RespBLOCKCHAIN(&newPeer)
}

func GetPeers(w http.ResponseWriter, r *http.Request) {
	var addrs []string
	for _, p := range p2p.Peers {
		addrs = append(addrs, p.Taget)
	}
	b, e := json.Marshal(addrs)
	if e != nil {
		fmt.Println(e)
	}
	w.Write(b)
}

func GetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(block.BCs, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func WriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		fmt.Println("decoder err : ", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	lastBlock := block.GetLatestBlock()
	newBlock, err := lastBlock.GenerateBlock(m.Wallet)
	if err != nil {
		fmt.Println("GenerateBlock err : ", err)
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	Blockchain := block.BCs
	if newBlock.IsBlockValid() {
		newBlockchain := append(Blockchain, newBlock)
		block.ReplaceChain(newBlockchain)
		// spew.Dump(Blockchain)
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
