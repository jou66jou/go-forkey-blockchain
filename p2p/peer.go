package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/goinggo/mapstructure"

	"github.com/gorilla/websocket"
	"github.com/jou66jou/go-forky-blockchain/block"
)

// 節點
type Peer struct {
	socket *websocket.Conn
	send   chan []byte
	Taget  string
}

// 建立新節點
func NewPeer(conn *websocket.Conn, target string) Peer {
	return Peer{conn, make(chan []byte), target}
}

// 監聽訊息
func (p *Peer) Read() {
	defer func() {
		p.socket.Close()
	}()

	for {
		_, message, err := p.socket.ReadMessage()
		if err != nil {
			p.socket.Close()
			break
		}
		m := msg{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			fmt.Println("Peer Read() err : " + err.Error())
			continue
		}
		switch m.Event {
		case ADD_PEER:
			ConnectionToAddr(m.Content.(string), true)
		case RESPONSE_BLOCKCHAIN:
			var newBCs []block.Block
			if err := mapstructure.Decode(m.Content, &newBCs); err != nil { // map to slice
				fmt.Println("mapstructure err : ", err)
				continue
			}
			block.ReplaceChain(newBCs)
		}
	}
}

// send channel有訊息時寫入websocket
func (p *Peer) Write() {
	defer func() {
		p.socket.Close()
	}()

	for {
		select {
		case message, ok := <-p.send:
			if !ok {
				p.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			p.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
