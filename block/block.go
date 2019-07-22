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
	Wallet    int
}

var (
	BCs []Block
)

// 建立新block
func (block *Block) GenerateBlock(Wallet int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = block.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Wallet = Wallet
	newBlock.PrevHash = block.Hash
	newBlock.Hash = newBlock.CalculateHash()
	return newBlock, nil
}

// 產生一個block的SHA256
func (block *Block) CalculateHash() string {
	record := string(block.Index) + block.Timestamp + string(block.Wallet) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 驗證block
func (block *Block) IsBlockValid() bool {
	if GetLatestBlock().Index+1 != block.Index {
		return false
	}
	if GetLatestBlock().Hash != block.PrevHash {
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

//取得最後一塊block
func GetLatestBlock() Block {
	if len(BCs) == 0 { //若鏈上無區塊則產生初始block
		t := time.Now()
		genesisBlock := Block{0, t.String(), "", "", 0}
		genesisBlock.Hash = genesisBlock.CalculateHash()
		BCs = append(BCs, genesisBlock)
	}
	return BCs[len(BCs)-1]
}
