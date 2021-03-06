// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package builtin

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner/executor/common"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner/executor/common/foundation"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner/executor/common/rpctypes"
)

type ProxyHelper struct {
	common.Serializer
	common.SystemError
	methods common.RunnerRPCStub
}

func NewProxyHelper(runner common.RunnerRPCStub) *ProxyHelper {
	return &ProxyHelper{
		Serializer:  common.NewCBORSerializer(),
		SystemError: common.NewSystemError(),
		methods:     runner,
	}
}

func (h *ProxyHelper) getUpBaseReq() rpctypes.UpBaseReq {
	callContext := foundation.GetLogicalContext()

	return rpctypes.UpBaseReq{
		Mode:            callContext.Mode,
		Callee:          *callContext.Callee,
		CalleePrototype: *callContext.CallerPrototype,
		Request:         *callContext.Request,
	}
}

func (h *ProxyHelper) RouteCall(ref insolar.Reference, immutable bool, saga bool, method string, args []byte,
	proxyPrototype insolar.Reference) ([]byte, error) {

	if h.GetSystemError() != nil {
		return nil, h.GetSystemError()
	}

	res := rpctypes.UpRouteResp{}
	req := rpctypes.UpRouteReq{
		UpBaseReq: h.getUpBaseReq(),

		Object:    ref,
		Immutable: immutable,
		Saga:      saga,
		Method:    method,
		Arguments: args,
		Prototype: proxyPrototype,
	}

	err := h.methods.RouteCall(req, &res)
	if err != nil {
		h.SetSystemError(err)
		return nil, err
	}

	return res.Result, nil
}

func (h *ProxyHelper) SaveAsChild(
	parentRef, classRef insolar.Reference,
	constructorName string, argsSerialized []byte,
) (
	[]byte, error,
) {
	if !parentRef.IsObjectReference() {
		return nil, errors.Errorf("Failed to save AsChild: objRef should be ObjectReference; ref=%s", parentRef.String())
	}

	if h.GetSystemError() != nil {
		// There was a system error during execution of the contract.
		// Immediately return this error to the calling contract - any
		// results will not be registered on LME anyway.
		return nil, h.GetSystemError()
	}

	res := rpctypes.UpSaveAsChildResp{}
	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq: h.getUpBaseReq(),

		Parent:          parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	err := h.methods.SaveAsChild(req, &res)
	if err != nil {
		h.SetSystemError(err)
		return nil, err
	}

	return res.Result, nil
}

func (h *ProxyHelper) DeactivateObject(object insolar.Reference) error {
	if h.GetSystemError() != nil {
		return h.GetSystemError()
	}

	res := rpctypes.UpDeactivateObjectResp{}
	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: h.getUpBaseReq(),
	}

	if err := h.methods.DeactivateObject(req, &res); err != nil {
		h.SetSystemError(err)
		return err
	}
	return nil
}

func (h *ProxyHelper) MakeErrorSerializable(err error) error {
	if err == nil || err == (*foundation.Error)(nil) || reflect.ValueOf(err).IsNil() {
		return nil
	}
	return &foundation.Error{S: err.Error()}
}
