// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package smachine

import (
	"context"
	"sync"
)

var _ AdapterExecutor = &OverflowPanicCallChannel{}

func NewCallChannelExecutor(ctx context.Context, bufMax int, requestCancel bool, parallelReaders int) (AdapterExecutor, chan AdapterCall) {
	if parallelReaders <= 0 {
		panic("illegal value")
	}
	switch output := make(chan AdapterCall, parallelReaders<<1); {
	case bufMax == 0:
		return WrapCallChannelNoBuffer(requestCancel, output), output
	case bufMax < 0: //unlimited buffer
		return WrapCallChannelNoLimit(ctx, requestCancel, output), output
	default:
		return WrapCallChannel(ctx, bufMax, requestCancel, output), output
	}
}

func WrapCallChannel(ctx context.Context, bufMax int, requestCancel bool, output chan AdapterCall) *OverflowBufferCallChannel {
	if bufMax <= 0 {
		panic("illegal value")
	}
	return &OverflowBufferCallChannel{ctx: ctx, output: output, needCancel: requestCancel, bufMax: bufMax}
}

func WrapCallChannelNoLimit(ctx context.Context, requestCancel bool, output chan AdapterCall) *OverflowBufferCallChannel {
	return &OverflowBufferCallChannel{ctx: ctx, output: output, needCancel: requestCancel}
}

func WrapCallChannelNoBuffer(requestCancel bool, output chan AdapterCall) OverflowPanicCallChannel {
	return OverflowPanicCallChannel{output, requestCancel}
}

type channelRecord = AdapterCall

type OverflowPanicCallChannel struct {
	output     chan channelRecord
	needCancel bool
}

func (v OverflowPanicCallChannel) TrySyncCall(AdapterCallFunc) (bool, AsyncResultFunc) {
	return false, nil
}

func (v OverflowPanicCallChannel) StartCall(fn AdapterCallFunc, callback *AdapterCallback, needCancel bool) context.CancelFunc {
	switch {
	case fn == nil:
		panic("illegal value")
	case callback == nil:
		panic("illegal value")
	}

	r := channelRecord{CallFn: fn, Callback: callback}
	cancelFn := callback.Prepare(needCancel || v.needCancel)
	v.queueCall(r)

	return cancelFn
}

func (v OverflowPanicCallChannel) SendNotify(fn AdapterNotifyFunc) {
	if fn == nil {
		panic("illegal value")
	}
	r := channelRecord{CallFn: func(svc interface{}) AsyncResultFunc {
		fn(svc)
		return nil
	}}
	v.queueCall(r)
}

func (v OverflowPanicCallChannel) Channel() chan channelRecord {
	return v.output
}

func (v OverflowPanicCallChannel) queueCall(r channelRecord) {
	select {
	case v.output <- r:
	default:
		panic("overflow")
	}
}

var _ AdapterExecutor = &OverflowBufferCallChannel{}

// This wrapper doesn't allocate a buffer unless the channel is full
type OverflowBufferCallChannel struct {
	ctx        context.Context
	mutex      sync.Mutex
	output     chan channelRecord
	buffer     []channelRecord
	bufMax     int
	needCancel bool
}

func (p *OverflowBufferCallChannel) TrySyncCall(AdapterCallFunc) (bool, AsyncResultFunc) {
	return false, nil
}

func (p *OverflowBufferCallChannel) StartCall(fn AdapterCallFunc, callback *AdapterCallback, needCancel bool) context.CancelFunc {
	switch {
	case fn == nil:
		panic("illegal value")
	case callback == nil:
		panic("illegal value")
	}

	r := channelRecord{CallFn: fn, Callback: callback}
	cancelFn := callback.Prepare(needCancel || p.needCancel)
	p.queueCall(r)

	return cancelFn
}

func (p *OverflowBufferCallChannel) SendNotify(fn AdapterNotifyFunc) {
	if fn == nil {
		panic("illegal value")
	}
	r := channelRecord{CallFn: func(svc interface{}) AsyncResultFunc {
		fn(svc)
		return nil
	}}
	p.queueCall(r)
}

func (p *OverflowBufferCallChannel) queueCall(r channelRecord) {
	switch {
	case p.append(r, false):
	case p.send(r):
	case p.append(r, true):
	default:
		panic("overflow")
	}
}

func (p *OverflowBufferCallChannel) Channel() chan channelRecord {
	return p.output
}

func (p *OverflowBufferCallChannel) Context() context.Context {
	return p.ctx
}

func (p *OverflowBufferCallChannel) append(r channelRecord, force bool) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch {
	case len(p.buffer) > 0:
		break
	case !force:
		return false
	default:
		go p.sendWorker() // wont start because of lock
	}

	if p.bufMax > 0 && len(p.buffer) >= p.bufMax {
		panic("overflow")
	}
	p.buffer = append(p.buffer, r)
	return true
}

func (p *OverflowBufferCallChannel) send(r channelRecord) bool {
	select {
	case p.output <- r:
		return true
	default:
		return false
	}
}

const bigBufferTrimSize = 65536

func (p *OverflowBufferCallChannel) sendWorker() {
	defer func() {
		_ = recover()
	}()

	bufPos := 0
	for bufPos >= 0 {
		p.mutex.Lock()
		r := p.buffer[bufPos]
		p.buffer[bufPos] = channelRecord{}
		bufPos++

		switch n := len(p.buffer); {
		case bufPos == n:
			p.buffer = p.buffer[:0]
			bufPos = -1
		case p.bufMax > 0 && bufPos > p.bufMax>>1:
			fallthrough
		case n > bigBufferTrimSize && bufPos > bigBufferTrimSize>>1:
			copy(p.buffer, p.buffer[bufPos:])
			p.buffer = p.buffer[:n-bufPos]
			bufPos = 0
		}
		p.mutex.Unlock()

		select {
		case <-p.ctx.Done():
			return
		case p.output <- r:
		}
	}
}
