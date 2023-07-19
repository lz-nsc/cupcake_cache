package cupcake_cache

import (
	"context"
	"fmt"
	"net"

	"github.com/lz-nsc/cupcake_cache/log"
	pb "github.com/lz-nsc/cupcake_cache/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type cacheGrpcServer struct {
	pb.GroupCacheServer
	name        string
	addr        string
	cacheServer *CacheServer
}

var _ Server = (*cacheGrpcServer)(nil)

func init() {
	serverMap["grpc"] = NewCacheGrpcServer
}

func NewCacheGrpcServer(name string, addr string, cs *CacheServer) Server {
	log.WithServer(name, addr)
	return &cacheGrpcServer{
		name:        name,
		addr:        addr,
		cacheServer: cs,
	}
}

func (cs *cacheGrpcServer) Run() error {
	lis, err := net.Listen("tcp", cs.addr)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return err
	}

	server := grpc.NewServer()
	pb.RegisterGroupCacheServer(server, cs)

	return server.Serve(lis)
}
func (cs *cacheGrpcServer) Get(cxt context.Context, in *pb.Request) (*pb.Response, error) {
	log.Infof("new grpc request: %s %s", in.Group, in.Key)

	groupName := in.Group
	key := in.Key
	group := GetGroup(groupName)

	if group == nil {
		return nil, fmt.Errorf("gourp %s not found", groupName)
	}

	bv, err := group.Get(key)
	if err != nil {

		return nil, err
	}

	return &pb.Response{
		Value: bv.bytes,
	}, nil
}

func (remote *cacheGrpcServer) RemoteGet(group string, key string) ([]byte, error) {
	req := &pb.Request{
		Group: group,
		Key:   key,
	}

	peerAddr, ok := remote.cacheServer.GetNode(key)
	if !ok {
		return nil, fmt.Errorf("failed to get remote node with key %s", key)
	}

	if peerAddr == "" {
		return nil, nil
	}
	log.Debugf("Successfully got remote peer, addr: %s", peerAddr)

	conn, err := grpc.Dial(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to get data from remote cache, err: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := pb.NewGroupCacheClient(conn)

	resp, err := client.Get(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}
