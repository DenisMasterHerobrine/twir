package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Activity) ChangeCategory(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
	hydratedCategory, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)

	if err != nil || len(hydratedCategory) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	twitchClient, twitchClientErr := c.getHelixApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	searchCategory, err := twitchClient.SearchCategories(
		&helix.SearchCategoriesParams{
			Query: hydratedCategory,
		},
	)
	if err != nil {
		return err
	}
	if searchCategory.ErrorMessage != "" {
		return fmt.Errorf("cannot get category with that name: %s", searchCategory.ErrorMessage)
	}

	if len(searchCategory.Data.Categories) == 0 {
		return errors.New("cannot get category with that name")
	}

	editReq, err := twitchClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: data.ChannelID,
			GameID:        searchCategory.Data.Categories[0].ID,
		},
	)
	if err != nil {
		return err
	}

	if editReq.ErrorMessage != "" {
		return errors.New(editReq.ErrorMessage)
	}

	return nil
}

func (c *Activity) ChangeTitle(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
	hydratedTitle, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)

	if err != nil || len(hydratedTitle) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	twitchClient, twitchClientErr := c.getHelixApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	req, err := twitchClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: data.ChannelID,
			Title:         hydratedTitle,
		},
	)
	if err != nil {
		return err
	}

	if req.ErrorMessage != "" {
		return errors.New(req.ErrorMessage)
	}

	return nil
}
