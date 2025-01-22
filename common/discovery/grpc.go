package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, servicename string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.Discover(ctx, servicename)
	if err != nil {
		return nil, err
	}

	/// Pick random service among all
	return grpc.NewClient(
		addrs[rand.Intn(len(addrs))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
