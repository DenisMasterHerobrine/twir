package events

import (
	"context"
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Activity) SwitchEmoteOnly(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
	dbEntity, dbEntityErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixApiClient(ctx, dbEntity.BotID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	resp, err := twitchClient.UpdateChatSettings(
		&helix.UpdateChatSettingsParams{
			BroadcasterID: data.ChannelID,
			ModeratorID:   dbEntity.BotID,
			EmoteMode: lo.ToPtr(
				lo.
					If(operation.Type == model.OperationEnableEmoteOnly, true).
					Else(false),
			),
		},
	)
	if err != nil {
		return err
	}
	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}

func (c *Activity) SwitchSubMode(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
	dbEntity, dbEntityErr := c.getChannelDbEntity(ctx, data.ChannelID)
	if dbEntityErr != nil {
		return dbEntityErr
	}

	twitchClient, twitchClientErr := c.getHelixApiClient(ctx, dbEntity.BotID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	resp, err := twitchClient.UpdateChatSettings(
		&helix.UpdateChatSettingsParams{
			BroadcasterID: data.ChannelID,
			ModeratorID:   dbEntity.BotID,
			SubscriberMode: lo.ToPtr(
				lo.
					If(operation.Type == model.OperationEnableSubMode, true).
					Else(false),
			),
		},
	)
	if err != nil {
		return err
	}
	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}
