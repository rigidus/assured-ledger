//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package unsafekit

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/longbits"
)

// WARNING! The given struct MUST be immutable. Expects struct ptr.
// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (v) stays alive while the resulting ByteString is in use.
func WrapSlice(v interface{}) longbits.ByteString {
	vt := reflect.ValueOf(v)
	if vt.Kind() != reflect.Slice {
		panic("illegal value")
	}
	return wrapSlice(vt)
}

func wrapSlice(vt reflect.Value) longbits.ByteString {
	n := uintptr(vt.Len()) * vt.Type().Elem().Size()
	if n == 0 {
		return ""
	}
	return wrapUnsafePtr(vt.Pointer(), n)
}

func WrapSliceOf(v interface{}, mt MMapSliceType) longbits.ByteString {
	vt := reflect.ValueOf(v)
	if vt.Type() != mt.ReflectType() {
		panic("illegal value type")
	}
	return wrapSlice(vt)
}

func UnwrapAsSliceOf(s longbits.ByteString, mt MMapSliceType) interface{} {
	t := mt.ReflectType()
	if t.Kind() != reflect.Slice { // double-check
		panic("illegal value")
	}

	itemSize := int(t.Elem().Size())
	switch {
	case len(s) == 0:
		return reflect.Zero(t).Interface()
	case len(s)%itemSize != 0:
		panic(fmt.Sprintf("illegal value - length is unaligned: dataLen=%d itemSize=%d", len(s), itemSize))
	}

	slice := reflect.New(t)
	itemCount := len(s) / itemSize
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(slice.Pointer()))
	sliceHeader.Data = _unwrapUnsafeUintptr(s)
	sliceHeader.Cap = itemCount
	sliceHeader.Len = itemCount

	slice = slice.Elem()
	if slice.Len() != itemCount {
		panic("unexpected")
	}
	runtime.KeepAlive(s)

	return slice.Interface()
}
