package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) RefreshChatOverlaySettings(
	ctx context.Context,
	req *websockets.RefreshChatSettingsRequest,
) (
	*emptypb.Empty,
	error,
) {
	if err := c.chatServer.SendSettings(req.ChannelId); err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
