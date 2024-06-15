package usecase

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"context"
)

func (u *botUC) HandleActionsWithMainKeyboardCallback(ctx context.Context, params models.HandleActionsWithMainKeyboardCallback) error {
	userInfo, err := u.pgRepo.GetUserInfo(ctx, params.ChatID)
	if err != nil {
		return err
	}

	if !userInfo.IsAdmin || !userInfo.IsActive {
		return nil
	}

	sendMessageParams := models.EditMessageParams{
		ChatID:    params.ChatID,
		MessageID: params.MessageID,
		Keyboard:  u.createActionsWithMainKeyboard(),
		Text:      "select the action you want to perform on the main keyboard"}
	u.editMessage(sendMessageParams)

	return nil
}
