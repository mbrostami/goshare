package grpc

import (
	"context"
	"fmt"
	"github.com/mbrostami/goshare/api/grpc/pb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"strings"
)

type Client struct {
	conn pb.GoShareClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	var client Client
	client.conn = pb.NewGoShareClient(conn)

	return &client, nil
}

func NewClientLoadBalancer(servers []string) (*Client, error) {
	var addresses []resolver.Address

	log.Debug().Msgf("connecting to servers: %s", strings.Join(servers, " - "))
	for _, server := range servers {
		addresses = append(addresses, resolver.Address{Addr: server})
	}

	res := manual.NewBuilderWithScheme("manual")
	res.InitialState(resolver.State{Addresses: addresses})

	balancer.Register(balancer.Get(roundrobin.Name))

	conn, err := grpc.Dial(
		res.Scheme()+":///test",
		grpc.WithResolvers(res),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, roundrobin.Name),
		),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	var client Client
	client.conn = pb.NewGoShareClient(conn)

	return &client, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.conn.Ping(ctx, &pb.PingMsg{Ping: true})
	return err
}
