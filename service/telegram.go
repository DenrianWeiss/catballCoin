package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	teleBot *tgbotapi.BotAPI
)

func telegramInit() {
	if !GlobalConfig.Telegram.Enable {
		return
	}
	bot, err := tgbotapi.NewBotAPI(GlobalConfig.Telegram.BotToken)
	if err != nil {
		panic("Init panic: Invalid bot.")
	}
	teleBot = bot
}

func SendTelegramMessage(message string, chatID int64) {
	if !GlobalConfig.Telegram.Enable {
		return
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, _ = teleBot.Send(msg)
	return
}
