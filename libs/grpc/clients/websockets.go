package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/websockets"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewWebsocket(env string) websockets.WebsocketClient {
	serverAddress := createClientAddr(env, "websockets", constants.WEBSOCKET_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(
				`{"loadBalancingConfig": [{"%s":{}}]}`,
				roundrobin.Name,
			),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := websockets.NewWebsocketClient(conn)

	return c
}
