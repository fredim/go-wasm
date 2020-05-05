package wasm

/*
#include <stdlib.h>

extern int32_t wasiFdWrite(void *context, int32_t a, int32_t b, int32_t c, int32_t d);
extern double runtimeTicks(void *context);
extern void finalizeRef(void *context, int32_t a, int32_t x, int32_t y);
extern void stringValTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t x, int32_t y);
extern void valueGetTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t x, int32_t y);
extern void valueSetTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t x, int32_t y);
extern void valueIndexTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t x, int32_t y);
extern void valueSetIndexTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t x, int32_t y);
extern void valueCallTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t e, int32_t f, int32_t g, int32_t x, int32_t y);
extern void valueInvokeTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t e, int32_t x, int32_t y);
extern void valueNewTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t e, int32_t x, int32_t y);
extern int32_t valueLengthTiny(void *context, int32_t a, int32_t x, int32_t y);
extern void valuePrepareStringTiny(void *context, int32_t a, int32_t b, int32_t x, int32_t y);
extern void valueLoadStringTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t x, int32_t y);
extern void copyBytesToJSTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t e, int32_t x, int32_t y);
extern void copyBytesToGoTiny(void *context, int32_t a, int32_t b, int32_t c, int32_t d, int32_t e, int32_t x, int32_t y);

*/
import "C"
import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/wasmerio/go-ext-wasm/wasmer"
)

//export wasiFdWrite
func wasiFdWrite(_ unsafe.Pointer, _, _, _, _ int32) int32 {
	panic("don't call me")
}

//export finalizeRef
func finalizeRef(_ unsafe.Pointer, _, _, _ int32) {
	panic("don't call me")
}

//export runtimeTicks
func runtimeTicks(_ unsafe.Pointer) float64 {
	fmt.Println("calling runtimeTicks")
	return float64(time.Now().Unix())
}

//export stringValTiny
func stringValTiny(ctx unsafe.Pointer, ret, ptr, plen, x, y int32) {
	fmt.Println("calling stringValTiny")
	b := getBridge(ctx)
	str := b.loadStringTiny(ptr, plen)
	b.storeValue(ret, str)
}

func reflectGet(val interface{}, prop string) interface{} {
	if obj, ok := val.(*object); ok {
		if res, ok := obj.props[prop]; ok {
			return res
		}
		panic(fmt.Sprintln("missing property", prop, val))
	}
	return val
}

//export valueGetTiny
func valueGetTiny(ctx unsafe.Pointer, ret, vaddr, ptr, plen, x, y int32) {
	fmt.Println("calling valueGetTiny")
	b := getBridge(ctx)
	prop := b.loadStringTiny(ptr, plen)
	val := b.loadValue(vaddr)
	b.storeValue(ret, reflectGet(val, prop))
	// sp := b.getSP()
}

//export valueSetTiny
func valueSetTiny(ctx unsafe.Pointer, vaddr, ptr, plen, xref, x, y int32) {
	fmt.Println("calling valueSetTiny")
	b := getBridge(ctx)
	val := b.loadValue(vaddr)
	obj := val.(*object)
	prop := b.loadStringTiny(ptr, plen)
	propVal := b.loadValue(xref)
	obj.props[prop] = propVal
}

//export valueIndexTiny
func valueIndexTiny(ctx unsafe.Pointer, ret, vaddr, idx, x, y int32) {
	fmt.Println("calling valueIndexTiny")
	b := getBridge(ctx)
	l := b.loadValue(vaddr)
	i := b.getInt64(idx)
	rv := reflect.ValueOf(l)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	iv := rv.Index(int(i))
	b.storeValue(ret, iv.Interface())
}

//export valueSetIndexTiny
func valueSetIndexTiny(_ unsafe.Pointer, _, _, _, _, _ int32) {
	fmt.Println("calling valueSetIndexTiny")
	panic("valueSetIndexTiny")
}

//export valueCallTiny
func valueCallTiny(ctx unsafe.Pointer, ret, vaddr, mptr, mlen, argptr, arglen, argcap, x, y int32) {
	fmt.Println("calling valueCallTiny")
	b := getBridge(ctx)
	v := b.loadValue(vaddr)
	str := b.loadStringTiny(mptr, mlen)
	args := b.loadSliceOfValuesTiny(argptr, arglen, argcap)
	f, ok := v.(*object).props[str].(Func)
	if !ok {
		panic(fmt.Sprintf("valueCall: prop not found in %v, %s", v.(*object).name, str))
	}
	// sp = b.getSP()
	res, err := f(args)
	if err != nil {
		b.storeValue(ret, err.Error())
		b.setUint8(ret+8, 0)
		return
	}

	b.storeValue(ret, res)
	b.setUint8(ret+8, 1)
}

//export valueInvokeTiny
func valueInvokeTiny(ctx unsafe.Pointer, ret, vaddr, argptr, arglen, argcap, x, y int32) {
	fmt.Println("calling valueInvokeTiny")
	b := getBridge(ctx)
	val := *(b.loadValue(vaddr).(*Func))
	args := b.loadSliceOfValuesTiny(argptr, arglen, argcap)
	res, err := val(args)
	// sp = b.getSP()
	if err != nil {
		b.storeValue(ret, err)
		b.setUint8(ret+8, 0)
		return
	}

	b.storeValue(ret, res)
	b.setUint8(ret+8, 1)
}

//export valueNewTiny
func valueNewTiny(ctx unsafe.Pointer, ret, vaddr, argptr, arglen, argcap, x, y int32) {
	fmt.Println("calling valueNewTiny")
	b := getBridge(ctx)
	val := b.loadValue(vaddr)
	args := b.loadSliceOfValuesTiny(argptr, arglen, argcap)
	res := val.(*object).new(args)
	// sp = b.getSP()
	b.storeValue(ret, res)
	b.setUint8(ret+8, 1)
}

//export valueLengthTiny
func valueLengthTiny(ctx unsafe.Pointer, vaddr, x, y int32) int32 {
	fmt.Println("calling valueLengthTiny")
	b := getBridge(ctx)
	val := b.loadValue(vaddr)
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	var l int
	switch {
	case rv.Kind() == reflect.Slice:
		l = rv.Len()
	case rv.Type() == reflect.TypeOf(array{}):
		l = len(val.(*array).buf)
	default:
		panic(fmt.Sprintf("valueLength on %T", val))
	}
	return int32(l)
}

//export valuePrepareStringTiny
func valuePrepareStringTiny(ctx unsafe.Pointer, ret, vaddr, x, y int32) {
	fmt.Println("calling valuePrepareStringTiny")
	b := getBridge(ctx)
	val := b.loadValue(vaddr)
	var str string
	if val != nil {
		str = fmt.Sprint(val)
	}

	b.storeValue(ret, str)
	b.setInt64(ret+8, int64(len(str)))
}

//export valueLoadStringTiny
func valueLoadStringTiny(ctx unsafe.Pointer, vaddr, sptr, slen, scap, x, y int32) {
	fmt.Println("calling valueLoadStringTiny")
	b := getBridge(ctx)
	str := b.loadValue(vaddr).(string)
	sl := b.loadSliceTiny(sptr, slen, scap)
	copy(sl, str)
}

//export copyBytesToJSTiny
func copyBytesToJSTiny(ctx unsafe.Pointer, ret, dest, sptr, slen, scap, x, y int32) {
	fmt.Println("calling copyBytesToJSTiny")
	b := getBridge(ctx)
	dst, ok := b.loadValue(dest).(*array)
	if !ok {
		b.setUint8(ret+8, 0)
		return
	}
	src := b.loadSliceTiny(sptr, slen, scap)
	n := copy(dst.buf, src[:len(dst.buf)])
	b.setInt64(ret, int64(n))
	b.setUint8(ret+8, 1)
}

//export copyBytesToGoTiny
func copyBytesToGoTiny(ctx unsafe.Pointer, ret, dest, dlen, dcap, source, x, y int32) {
	fmt.Println("calling copyBytesToGoTiny")
	b := getBridge(ctx)
	dst := b.loadSliceTiny(dest, dlen, dcap)
	src, ok := b.loadValue(source).(*array)
	if !ok {
		b.setUint8(ret+8, 0)
		return
	}
	n := copy(dst, src.buf[:len(dst)])
	b.setInt64(ret, int64(n))
	b.setUint8(ret+8, 1)
}

type importList struct {
	name string
	imp  interface{}
	cgo  unsafe.Pointer
}

var tinyImports = map[string][]importList{
	// this set is for tinygo 0.13
	"wasi_unstable": {
		{"fd_write", wasiFdWrite, C.wasiFdWrite},
	},
	"env": {
		{"runtime.ticks", runtimeTicks, C.runtimeTicks},
		// {"runtime.sleepTicks", nil, nil},
		{"syscall/js.finalizeRef", finalizeRef, C.finalizeRef},
		{"syscall/js.stringVal", stringValTiny, C.stringValTiny},
		{"syscall/js.valueGet", valueGetTiny, C.valueGetTiny},
		{"syscall/js.valueSet", valueSetTiny, C.valueSetTiny},
		{"syscall/js.valueIndex", valueIndexTiny, C.valueIndexTiny},
		{"syscall/js.valueSetIndex", valueSetIndexTiny, C.valueSetIndexTiny}, //unsuported
		{"syscall/js.valueCall", valueCallTiny, C.valueCallTiny},
		{"syscall/js.valueInvoke", valueInvokeTiny, C.valueInvokeTiny},
		{"syscall/js.valueNew", valueNewTiny, C.valueNewTiny},
		{"syscall/js.valueLength", valueLengthTiny, C.valueLengthTiny},
		{"syscall/js.valuePrepareString", valuePrepareStringTiny, C.valuePrepareStringTiny},
		{"syscall/js.valueLoadString", valueLoadStringTiny, C.valueLoadStringTiny},
		{"syscall/js.copyBytesToGo", copyBytesToGoTiny, C.copyBytesToGoTiny},
		{"syscall/js.copyBytesToJS", copyBytesToJSTiny, C.copyBytesToJSTiny},
	},
}

// addTinyImports adds go Bridge imports in "env" namespace.
func (b *Bridge) addTinyImports(imps *wasmer.Imports) error {
	var err error
	for k, implist := range tinyImports {
		imps = imps.Namespace(k)
		for _, imp := range implist {
			imps, err = imps.Append(imp.name, imp.imp, imp.cgo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
