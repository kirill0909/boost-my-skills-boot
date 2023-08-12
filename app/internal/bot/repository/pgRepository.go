package repository

import (
	"boost-my-skills-bot/internal/bot"
)

type BotPGRepo struct{}

func NewBotPGRepo() bot.PgRepository {
	return &BotPGRepo{}
}
