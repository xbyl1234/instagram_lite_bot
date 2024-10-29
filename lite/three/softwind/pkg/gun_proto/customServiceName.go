package gun_proto

import (
	"context"
	"google.golang.org/grpc"
)

func ServerDesc(name string) grpc.ServiceDesc {
	return grpc.ServiceDesc{
		ServiceName: name,
		HandlerType: (*GunServiceServer)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "Tun",
				Handler:       _GunService_Tun_Handler,
				ServerStreams: true,
				ClientStreams: true,
			},
			{
				StreamName:    "TunDatagram",
				Handler:       _GunService_TunDatagram_Handler,
				ServerStreams: true,
				ClientStreams: true,
			},
		},
		Metadata: "gun_proto.proto",
	}
}

func (c *gunServiceClient) TunCustomName(ctx context.Context, name string, opts ...grpc.CallOption) (GunService_TunClient, error) {
	stream, err := c.cc.NewStream(ctx, &ServerDesc(name).Streams[0], "/"+name+"/Tun", opts...)
	if err != nil {
		return nil, err
	}
	x := &gunServiceTunClient{stream}
	return x, nil
}

type GunServiceClientX interface {
	TunCustomName(ctx context.Context, name string, opts ...grpc.CallOption) (GunService_TunClient, error)
	Tun(ctx context.Context, opts ...grpc.CallOption) (GunService_TunClient, error)
	TunDatagram(ctx context.Context, opts ...grpc.CallOption) (GunService_TunDatagramClient, error)
}

func RegisterGunServiceServerX(s *grpc.Server, srv GunServiceServer, name string) {
	desc := ServerDesc(name)
	s.RegisterService(&desc, srv)
}
