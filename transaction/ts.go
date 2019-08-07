package ts

type UspentTxOut struct {
	Id      string
	Index   int
	Address string
	Amount  int
}

type TxIn struct {
	TxOutId    string
	TxOutIndex int
	Sinature   string
}

type TxOut struct {
	Address string
	Amount  int
}

type Transaction struct {
	Id     string	// sha256
	TxIns  []TxIn
	TxOuts []TxOut
}

// 計算Ts包內容的sha256
func (ts *Transaction) GetId() {
	txInContent := 
	for _, in := range ts.TxIns {
		in.
	}
}

func (ts *Transaction) ValiDate(UnTxOuts []UspentTxOut) {

}


func (tIn *TxIn) GetTxInAmount(UnTxOuts []UspentTxOut) {

}

func findUnTxOut(tsId string, index int, UnTxOuts []UspentTxOut) {

}

func getCoinbaseTs(addr string, bIndex int) {

}

// 交易簽名
func signTxIn(ts Transaction, txInIndex int, privateKey string, UnTxOuts []UspentTxOut) {

}

func (ts *Transaction) GetTxInsContent() *string{
	
}