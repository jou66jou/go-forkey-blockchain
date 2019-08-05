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
	Id     string
	TxIns  []TxIn
	TxOuts []TxOut
}

func (ts *Transaction) GetId() {

	for _, in := range ts.TxIns {

	}
}

func (ts *Transaction) ValiDate(UnTxOuts []UspentTxOut) {

}

func ValiDateBlockTss(tss []Transaction, UnTxOuts []UspentTxOut, bIndex int) {

}

func HasDuplicates(txIns []TxIn) {

}

func ValiDateCoinBaseTx(ts Transaction, bIndex int) {

}

func (tIn *TxIn) ValiDateTxIn(ts Transaction, UnTxOuts []UspentTxOut) {

}

func (tIn *TxIn) GetTxInAmount(UnTxOuts []UspentTxOut) {

}

func findUnTxOut(tsId string, index int, UnTxOuts []UspentTxOut) {

}

func getCoinbaseTs(addr string, bIndex int) {

}

func signTxIn(ts Transaction, txInIndex int, privateKey string, UnTxOuts []UspentTxOut) {

}
