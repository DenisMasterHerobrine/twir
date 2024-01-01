package manage

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
)

var RemoveAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases remove",
		Description: null.StringFrom("Remove aliase from command"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result, nil
		}

		args := strings.Split(*parseCtx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result, nil
		}

		commandName := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		aliase := strings.ToLower(strings.ReplaceAll(strings.Join(args[1:], " "), "!", ""))

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, commandName).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, "Command not found.")
			return result, nil
		}

		if !lo.Contains(cmd.Aliases, aliase) {
			result.Result = append(result.Result, "That aliase not in the command")
			return result, nil
		}

		index := lo.IndexOf(cmd.Aliases, aliase)
		cmd.Aliases = append(cmd.Aliases[:index], cmd.Aliases[index+1:]...)

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&cmd).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot save command",
				Err:     err,
			}
		}

		result.Result = append(result.Result, "✅ Aliase removed.")
		return result, nil
	},
}
