package usecase

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"boost-my-skills-bot/app/pkg/utils"
	"context"
)

func (u *botUC) HandleActionsWithMainKeyboardCallback(ctx context.Context, params models.HandleActionsWithMainKeyboardCallbackParams) error {
	userInfo, err := u.pgRepo.GetUserInfo(ctx, params.ChatID)
	if err != nil {
		return err
	}

	if !userInfo.IsAdmin || !userInfo.IsActive {
		return nil
	}

	editMessageParams := models.EditMessageParams{
		ChatID:    params.ChatID,
		MessageID: params.MessageID,
		Keyboard:  u.createActionsWithMainKeyboard(),
		Text:      "select the action you want to perform on the main keyboard"}
	u.editMessage(editMessageParams)

	return nil
}

func (u *botUC) HandleAddForUserMainKeyboardActionCallback(ctx context.Context, params models.HandleAddForUserMainKeyboardActionCallbackParams) error {
	userInfo, err := u.pgRepo.GetUserInfo(ctx, params.ChatID)
	if err != nil {
		return err
	}

	if !userInfo.IsAdmin || !userInfo.IsActive {
		return nil
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingNewMainButtonNameForUserStatus}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	editMessageParams := models.EditMessageParams{
		ChatID:    params.ChatID,
		MessageID: params.MessageID,
		Keyboard:  u.createComeBackToServiceButtonsKeyboard(),
		Text:      "enter name of new button for user",
	}
	u.editMessage(editMessageParams)

	return nil
}

func (u *botUC) HandleComeBackeToServiceButtonsCallback(ctx context.Context, params models.HandleComeBackeToServiceButtonsCallbackParams) {
	editMessageParams := models.EditMessageParams{
		ChatID:    params.ChatID,
		MessageID: params.MessageID,
		Keyboard:  u.createServiceKeyboard(),
		Text:      "choose service action"}
	u.editMessage(editMessageParams)
}
