//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package serialization

import (
	"context"
	"encoding/binary"
	"io"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/consensus/gcpv2/api/phases"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/cryptkit"
)

var (
	defaultByteOrder = binary.BigEndian
)

type PacketHeaderAccessor interface {
	GetProtocolType() ProtocolType
	GetPacketType() phases.PacketType
	GetSourceID() insolar.ShortNodeID
	HasFlag(flag Flag) bool
	GetFlagRangeInt(from, to uint8) uint8
	IsRelayRestricted() bool
	IsBodyEncrypted() bool
}

type PacketHeaderModifier interface {
	SetFlag(flag Flag)
	ClearFlag(flag Flag)
}

type PacketBody interface {
	ContextSerializerTo
	ContextDeserializerFrom

	String(ctx PacketContext) string
}

type FieldContext uint

const (
	NoContext = FieldContext(iota)
	ContextMembershipAnnouncement
	ContextNeighbourAnnouncement
)

type PacketContext interface {
	PacketHeaderAccessor
	context.Context

	InContext(ctx FieldContext) bool
	SetInContext(ctx FieldContext)
	GetNeighbourNodeID() insolar.ShortNodeID
	SetNeighbourNodeID(nodeID insolar.ShortNodeID)
	GetAnnouncedJoinerNodeID() insolar.ShortNodeID
	SetAnnouncedJoinerNodeID(nodeID insolar.ShortNodeID)
}

type SerializeContext interface {
	PacketHeaderModifier
	PacketContext
}

type DeserializeContext interface {
	PacketContext
}

type SerializerTo interface {
	SerializeTo(ctx context.Context, writer io.Writer, digester cryptkit.DataDigester, signer cryptkit.DigestSigner) (int64, error)
}

type ContextSerializerTo interface {
	SerializeTo(ctx SerializeContext, writer io.Writer) error
}

type DeserializerFrom interface {
	DeserializeFrom(ctx context.Context, reader io.Reader) (int64, error)
}

type ContextDeserializerFrom interface {
	DeserializeFrom(ctx DeserializeContext, reader io.Reader) error
}
