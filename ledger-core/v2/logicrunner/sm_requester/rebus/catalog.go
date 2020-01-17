package rebus

import (
	"fmt"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
)

type MessageKey [28]byte

func (k *MessageKey) ByteString() string {
	return fmt.Sprintf("message-%s", k[:])
}

type MessageCatalog struct{}

func (c MessageCatalog) Get(ctx smachine.ExecutionContext, key MessageKey) SharedMessageAccessor {
	if v, ok := c.TryGet(ctx, key); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", key.ByteString()))
}

func (MessageCatalog) TryGet(ctx smachine.ExecutionContext, key MessageKey) (SharedMessageAccessor, bool) {
	if v := ctx.GetPublishedLink(key.ByteString()); v.IsAssignableTo((*SharedMessageAccessor)(nil)) {
		return SharedMessageAccessor{v}, true
	}
	return SharedMessageAccessor{}, false
}

func (MessageCatalog) Create(ctx smachine.ExecutionContext, key MessageKey, value *SharedResponseBucket) bool {
	sharedValue := ctx.Share(value, 0)
	sharedWrapper := SharedMessageAccessor{sharedValue}

	return ctx.Publish(key.ByteString(), sharedWrapper) && ctx.PublishGlobalAlias(key.ByteString())
}

func (MessageCatalog) Finish(ctx smachine.ExecutionContext, key MessageKey) bool {
	return ctx.UnpublishGlobalAlias(key.ByteString()) && ctx.Unpublish(key.ByteString())
}

// //////////////////////////////////////

type SharedMessageAccessor struct {
	smachine.SharedDataLink
}

func (a SharedMessageAccessor) Prepare(fn func(*SharedResponseBucket)) smachine.SharedDataAccessor {
	return a.PrepareAccess(func(data interface{}) bool {
		fn(data.(*SharedResponseBucket))
		return false
	})
}
