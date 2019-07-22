package p2p

const (
	ADD_PEER            = 0  // 新增節點
	QUERY_LATEST        = 1  // 取得最後一個區塊
	QUERY_ALL           = 2  // 取得全部區塊
	RESPONSE_BLOCKCHAIN = 3  // 回應區塊鏈
	ERROR               = 10 // 錯誤訊息
)
