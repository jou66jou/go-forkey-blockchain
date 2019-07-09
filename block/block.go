package block

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
	BPM       int //記錄心跳數

}

var (
	BCs []Block
)

// 建立新block
func (block *Block) GenerateBlock(BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = block.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = block.Hash
	newBlock.Hash = newBlock.CalculateHash()
	return newBlock, nil
}

// 產生一個block的SHA256
func (block *Block) CalculateHash() string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 驗證block
func (block *Block) IsBlockValid() bool {
	if BCs[len(BCs)-1].Index+1 != block.Index {
		return false
	}
	if BCs[len(BCs)-1].Hash != block.PrevHash {
		return false
	}
	if block.CalculateHash() != block.Hash {
		return false
	}
	return true
}

// 替換舊鏈
func ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(BCs) {
		BCs = newBlocks
	}
}
