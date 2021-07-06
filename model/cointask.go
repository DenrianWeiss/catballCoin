package model

type CoinTask struct {
	TaskID          int64  `json:"-"`
	TargetWallet    string `json:"target_wallet" binding:"required"`
	SourceKey       string `json:"source_key" binding:"required"`
	ChatID          int64  `json:"chat_id" binding:"required"`
	GasMode         string `json:"gas_mode" binding:"required"`
	EthRpc          string `json:"eth_rpc" binding:"required"`
	LowerThreshold  int64  `json:"threshold" binding:"required"`
	HigherThreshold int64  `json:"high_threshold" binding:"required"`
	Percents        int64  `json:"percents" binding:"required"`
}
