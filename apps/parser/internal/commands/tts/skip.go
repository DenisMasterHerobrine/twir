package tts

import (
	"context"

	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/libs/grpc/generated/websockets"
)

var SkipCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts skip",
		Description: null.StringFrom("Skip current saying message in tts"),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		_, err := parseCtx.Services.GrpcClients.WebSockets.TextToSpeechSkip(
			context.Background(), &websockets.TTSSkipMessage{
				ChannelId: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while sending message to tts service",
				Err:     err,
			}
		}

		return result, nil
	},
}
