package model

type TelegramConfig struct {
	BotToken string `json:"bot_token"`
	Enable   bool   `json:"enable"`
}

type GlobalConfig struct {
	RpcKey   string         `json:"rpc_key"`
	DBPath   string         `json:"db_path"`
	Interval int            `json:"interval"`
	Telegram TelegramConfig `json:"telegram_config"`
}
