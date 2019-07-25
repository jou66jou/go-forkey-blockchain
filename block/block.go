package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/jou66jou/go-forky-blockchain/common"
)

type Block struct {
	Index      int    `json:"index"`
	Timestamp  string `json:"timesp"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prehash"`
	Wallet     int    `json:"wallet"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
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

func (block *Block) findBlock() string {
	var checkHead string
	block.Nonce = 0
	for {
		h := block.CalculateHash() // 64個十六進位數字
		endIndex := block.Difficulty/4 + 1
		for i := 0; i < endIndex; i += 16 {
			// 一次僅能處理64個二進位==16個十六進位數字
			checkHead += h[i : i+(endIndex%16)+1]
		}
		block.Nonce += 1
	}
}

// 產生一個block的SHA256
func (block *Block) CalculateHash() string {
	record := string(block.Index) + block.Timestamp + string(block.Wallet) + block.PrevHash + string(block.Nonce)
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

// 驗證鏈
func BlockChainValid(newBlocks []Block) (event int, content interface{}) {
	if len(newBlocks) == 0 {
		fmt.Println("new blockchain len is 0")
		return -1, ""
	}

	lastNewBlock := newBlocks[len(newBlocks)-1]
	lastheldBlock := GetLatestBlock()
	if lastNewBlock.Index > lastheldBlock.Index {
		if lastNewBlock.IsBlockValid() {
			BCs = append(BCs, lastNewBlock)
			// 廣播新區塊
			fmt.Println("broadcast new block to other peer")
			return common.RESPONSE_BLOCKCHAIN, []Block{lastNewBlock}
		} else if len(newBlocks) == 1 {
			// 請求其他節點的鏈
			fmt.Println("query chain form other peer")
			return common.QUERY_ALL, ""
		} else {
			fmt.Println("replace now chain")
			BCs = newBlocks
			return -1, ""
		}
	}
	fmt.Println("new blockchain len is not longger than loacl blockchain")
	return -1, ""

}

func hashMatchesDifficulty(h string, diff int) bool {

}

// 取得最後一塊block
func GetLatestBlock() Block {
	if len(BCs) == 0 { // 若鏈上無區塊則產生初始block
		t := time.Now()
		genesisBlock := new(Block)
		genesisBlock.Timestamp = t.String()
		genesisBlock.Hash = genesisBlock.CalculateHash()
		BCs = append(BCs, *genesisBlock)
	}
	return BCs[len(BCs)-1]
}

func HexToBin(hex string) (string, error) {
	ui, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}

	format := fmt.Sprintf("%%0%db", len(hex)*4)
	return fmt.Sprintf(format, ui), nil
}
