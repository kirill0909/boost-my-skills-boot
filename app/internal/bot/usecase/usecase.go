package usecase

import (
	"boost-my-skills-bot/internal/bot"
)

type BotUC struct{}

func NewBotUC() bot.Usecase {
	return &BotUC{}
}
