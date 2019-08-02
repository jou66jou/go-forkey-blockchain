package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/jou66jou/go-forky-blockchain/common"
)

// 產生塊頻率
const BLOCK_GENERATION_INTERVAL = 10

// 調整難度週期
const DIFFICULTY_ADJUSTMENT_INTERVAL = 2

type Block struct {
	Index      int    `json:"index"`
	Timestamp  int64  `json:"timesp"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prehash"`
	Wallet     int    `json:"wallet"`
	Difficulty int    `json:"difficulty"`
	Nonce      uint64 `json:"nonce"`
}

var (
	BCs []Block
)

// 建立新block
func (block *Block) GenerateBlock(Wallet int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = block.Index + 1
	newBlock.Timestamp = t.Unix()
	newBlock.Wallet = Wallet
	newBlock.PrevHash = block.Hash
	newBlock.Difficulty = newBlock.GetDifficulty()
	err := newBlock.findBlock()
	if err != nil {
		return newBlock, err
	}
	return newBlock, nil
}

// 產生工作量證明 proof of work
func (block *Block) findBlock() error {
	var err error
	block.Nonce = 0
	if block.Difficulty < 1 {
		h := block.CalculateHash() // 64個十六進位數字
		block.Hash = h
		return nil
	}
	// 固定資料先轉為字串，計算hash時就不必一直重複轉換
	record := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + strconv.Itoa(block.Wallet) + block.PrevHash + strconv.Itoa(block.Difficulty)
	for {
		var checkHead, binStr string

		s := record + strconv.FormatUint(block.Nonce, 10) //基礎字串＋變動數字
		h := GetHash(s)                                   // 64個十六進位數字

		// 十六轉二進制一次僅能處理16個十六進位=>64個二進位
		// 從hash字串第一位開始，提取 {Difficulty/4 無條件進位} 的長度出來
		// 例如Difficulty = 6，則從
		i := 0
		endIndex := block.Difficulty / 4
		if block.Difficulty%4 != 0 {
			endIndex += 1
		}
		for {
			if endIndex >= 16 {
				binStr, err = HexToBin(h[i : i+16])
				if err != nil {
					return err
				}
			} else {
				binStr, err = HexToBin(h[i : i+endIndex])
				if err != nil {
					return err
				}
				checkHead += binStr
				break
			}

			checkHead += binStr //累加二進位字串
			endIndex -= 16
			i += 16
		}

		if hasMatchesDif(block.Difficulty, binStr) { // 確認前導零個數
			block.Hash = h
			return nil
		}
		block.Nonce++
	}
}

// 產生一個block的SHA256
func (block *Block) CalculateHash() string {
	var record string
	record += strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + strconv.Itoa(block.Wallet) + block.PrevHash + strconv.Itoa(block.Difficulty) + strconv.FormatUint(block.Nonce, 10)
	return GetHash(record)
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
	if !block.IsBlockValid() {
		return false
	}
	return true
}

// 時間戳記驗證
func (b *Block) IsTimeVaild() bool {
	lastBlock := GetLatestBlock()
	return b.Timestamp > (lastBlock.Timestamp-60) && b.Timestamp < (lastBlock.Timestamp+60)
}

// 取得Difficulty
func (b *Block) GetDifficulty() int {
	if b.Index%DIFFICULTY_ADJUSTMENT_INTERVAL == 0 && b.Index != 0 {
		return AdjustedDif()
	}
	return BCs[len(BCs)-1].Difficulty
}

// 驗證鏈
func BlockChainValid(c *[]Block) (event int, content interface{}) {
	if len(*c) == 0 {
		fmt.Println("new blockchain len is 0")
		return -1, ""
	}

	lastNewBlock := (*c)[len(*c)-1]
	lastheldBlock := GetLatestBlock()
	if lastNewBlock.Index > lastheldBlock.Index {
		if lastNewBlock.IsBlockValid() {
			BCs = append(BCs, lastNewBlock)
			// 廣播新區塊
			fmt.Println("broadcast new block to other peer")
			return common.RESPONSE_BLOCKCHAIN, []Block{lastNewBlock}
		} else if len(*c) == 1 {
			// 請求其他節點的鏈
			fmt.Println("query chain form other peer")
			return common.QUERY_ALL, ""
		} else {
			// 計算Difficulty
			if GetAccumulateDif(c) > GetAccumulateDif(&BCs) {
				fmt.Println("replace now chain")
				BCs = *c
				return -1, ""
			} else {
				fmt.Println("replace chain fail")
				return -1, ""
			}

		}
	}
	fmt.Println("new blockchain len is not longger than loacl blockchain")
	return -1, ""

}

// 統計鏈的difficulty
func GetAccumulateDif(c *[]Block) int {
	var countDiff = 0
	for _, b := range *c {
		countDiff += b.Difficulty
	}
	return countDiff
}

// 取得最後一塊block
func GetLatestBlock() Block {
	if len(BCs) == 0 { // 若鏈上無區塊則產生初始block
		t := time.Now()
		genesisBlock := new(Block)
		genesisBlock.Timestamp = t.Unix()
		genesisBlock.Wallet = 50
		genesisBlock.Difficulty = 1
		genesisBlock.findBlock()
		BCs = append(BCs, *genesisBlock)
	}
	return BCs[len(BCs)-1]
}

// Hash256
func GetHash(s string) string {
	d := sha256.New()
	d.Write([]byte(s))
	hashed := d.Sum(nil)
	return hex.EncodeToString(hashed)
}

// 調整Difficulty
func AdjustedDif() int {
	lastBlock := BCs[len(BCs)-1]
	preAdjBlock := BCs[len(BCs)-DIFFICULTY_ADJUSTMENT_INTERVAL]
	timeExpected := int64(BLOCK_GENERATION_INTERVAL * DIFFICULTY_ADJUSTMENT_INTERVAL)
	timeTaken := lastBlock.Timestamp - preAdjBlock.Timestamp

	if timeTaken < timeExpected/2 {
		return preAdjBlock.Difficulty + 1
	} else if timeTaken > timeExpected*2 {
		return preAdjBlock.Difficulty - 1
	}
	return preAdjBlock.Difficulty

}

// 十六進位轉二進位（max:16 hex numbers）
func HexToBin(hex string) (string, error) {
	ui, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}

	format := fmt.Sprintf("%%0%db", len(hex)*4)
	return fmt.Sprintf(format, ui), nil
}

// 驗證Difficulty
func hasMatchesDif(dif int, binStr string) bool {
	for i := 0; i < dif; i++ {
		if binStr[i] != '0' {
			return false
		}
	}
	return true
}
