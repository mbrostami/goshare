package grpc

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/resolver/manual"

	"google.golang.org/grpc/balancer"

	"google.golang.org/grpc/resolver"

	"github.com/mbrostami/goshare/api/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

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
