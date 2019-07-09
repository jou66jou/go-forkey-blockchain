package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jou66jou/go-forky-blockchain/block"
)

type Message struct {
	Wallet int
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
	Blockchain := block.BCs
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	newBlock, err := Blockchain[len(Blockchain)-1].GenerateBlock(m.Wallet)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
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
