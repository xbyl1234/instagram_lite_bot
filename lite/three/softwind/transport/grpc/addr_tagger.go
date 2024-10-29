package grpc

import (
	"context"
	"google.golang.org/grpc/stats"
	"sync"
)

type addrTagger struct {
	ConnTagInfo *stats.ConnTagInfo
	mu          sync.Mutex
}

func (t *addrTagger) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	t.mu.Lock()
	defer t.mu.Unlock()
	// One GRPC Conn has at most one conn, and we should update it after reconnecting.
	t.ConnTagInfo = cti
	return ctx
}
func (t *addrTagger) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context { return ctx }
func (t *addrTagger) HandleRPC(context.Context, stats.RPCStats)                         {}
func (t *addrTagger) HandleConn(context.Context, stats.ConnStats)                       {}
