package wasm

/*
#include <stdlib.h>

extern void debug(void *context, int32_t a);
extern void wexit(void *context, int32_t a);
extern void wwrite(void *context, int32_t a);
extern void nanotime(void *context, int32_t a);
extern void walltime(void *context, int32_t a);
extern void scheduleCallback(void *context, int32_t a);
extern void clearScheduledCallback(void *context, int32_t a);
extern void getRandomData(void *context, int32_t a);
extern void stringVal(void *context, int32_t a);
extern void valueGet(void *context, int32_t a);
extern void valueSet(void *context, int32_t a);
extern void valueIndex(void *context, int32_t a);
extern void valueSetIndex(void *context, int32_t a);
extern void valueCall(void *context, int32_t a);
extern void valueInvoke(void *context, int32_t a);
extern void valueNew(void *context, int32_t a);
extern void valueLength(void *context, int32_t a);
extern void valuePrepareString(void *context, int32_t a);
extern void valueLoadString(void *context, int32_t a);
extern void scheduleTimeoutEvent(void *context, int32_t a);
extern void clearTimeoutEvent(void *context, int32_t a);
extern void copyBytesToGo (void *context, int32_t a);
extern void copyBytesToJS (void *context, int32_t a);

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
	"crypto/rand"
	"fmt"
	"log"
	"reflect"
	"syscall"
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

//export debug
func debug(_ unsafe.Pointer, sp int32) {
	log.Println(sp)
}

//export wexit
func wexit(ctx unsafe.Pointer, sp int32) {
	b := getBridge(ctx)
	b.exitCode = int(b.getUint32(sp + 8))
	b.cancF()
}

//export wwrite
func wwrite(ctx unsafe.Pointer, sp int32) {
	b := getBridge(ctx)
	fd := int(b.getInt64(sp + 8))
	p := int(b.getInt64(sp + 16))
	l := int(b.getInt32(sp + 24))
	_, err := syscall.Write(fd, b.mem()[p:p+l])
	if err != nil {
		panic(fmt.Errorf("wasm-write: %v", err))
	}
}

//export nanotime
func nanotime(ctx unsafe.Pointer, sp int32) {
	n := time.Now().UnixNano()
	getBridge(ctx).setInt64(sp+8, n)
}

//export runtimeTicks
func runtimeTicks(_ unsafe.Pointer) float64 {
	fmt.Println("calling runtimeTicks")
	return float64(time.Now().Unix())
}

//export walltime
func walltime(ctx unsafe.Pointer, sp int32) {
	b := getBridge(ctx)
	t := time.Now().UnixNano()
	nanos := t % int64(time.Second)
	b.setInt64(sp+8, t/int64(time.Second))
	b.setInt32(sp+16, int32(nanos))

}

//export scheduleCallback
func scheduleCallback(_ unsafe.Pointer, _ int32) {
	panic("schedule callback")
}

//export clearScheduledCallback
func clearScheduledCallback(_ unsafe.Pointer, _ int32) {
	panic("clear scheduled callback")
}

//export getRandomData
func getRandomData(ctx unsafe.Pointer, sp int32) {
	s := getBridge(ctx).loadSlice(sp + 8)
	_, err := rand.Read(s)
	if err != nil {
		panic("failed: getRandomData")
	}
}

//export stringVal
func stringVal(ctx unsafe.Pointer, sp int32) {
	stringValTiny(ctx, sp+24, sp+8, sp+16, 0, 0)
}

//export stringValTiny
func stringValTiny(ctx unsafe.Pointer, ret, ptr, plen, x, y int32) {
	fmt.Println("calling stringValTiny")
	b := getBridge(ctx)
	str := b.loadStringTiny(ptr, plen)
	b.storeValue(ret, str)
}

//export valueGet
func valueGet(ctx unsafe.Pointer, sp int32) {
	valueGetTiny(ctx, sp+32, sp+8, sp+16, sp+24, 0, 0)
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

//export valueSet
func valueSet(ctx unsafe.Pointer, sp int32) {
	valueSetTiny(ctx, sp+8, sp+16, sp+24, sp+32, 0, 0)
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

//export valueIndex
func valueIndex(ctx unsafe.Pointer, sp int32) {
	valueIndexTiny(ctx, sp+24, sp+8, sp+16, 0, 0)
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

//export valueSetIndex
func valueSetIndex(_ unsafe.Pointer, _ int32) {
	panic("valueSetIndex")
}

//export valueSetIndexTiny
func valueSetIndexTiny(_ unsafe.Pointer, _, _, _, _, _ int32) {
	fmt.Println("calling valueSetIndexTiny")
	panic("valueSetIndexTiny")
}

//export valueCall
func valueCall(ctx unsafe.Pointer, sp int32) {
	valueCallTiny(ctx, sp+56, sp+8, sp+16, sp+24, sp+32, sp+40, sp+48, 0, 0)
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

//export valueInvoke
func valueInvoke(ctx unsafe.Pointer, sp int32) {
	valueInvokeTiny(ctx, sp+40, sp+8, sp+16, sp+24, sp+32, 0, 0)
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

//export valueNew
func valueNew(ctx unsafe.Pointer, sp int32) {
	valueNewTiny(ctx, sp+40, sp+8, sp+16, sp+24, sp+32, 0, 0)
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

//export valueLength
func valueLength(ctx unsafe.Pointer, sp int32) {
	b := getBridge(ctx)
	b.setInt64(sp+16, int64(valueLengthTiny(ctx, sp+8, 0, 0)))
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

//export valuePrepareString
func valuePrepareString(ctx unsafe.Pointer, sp int32) {
	valuePrepareStringTiny(ctx, sp+16, sp+8, 0, 0)
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

//export valueLoadString
func valueLoadString(ctx unsafe.Pointer, sp int32) {
	valueLoadStringTiny(ctx, sp+8, sp+16, sp+24, sp+32, 0, 0)
}

//export valueLoadStringTiny
func valueLoadStringTiny(ctx unsafe.Pointer, vaddr, sptr, slen, scap, x, y int32) {
	fmt.Println("calling valueLoadStringTiny")
	b := getBridge(ctx)
	str := b.loadValue(vaddr).(string)
	sl := b.loadSliceTiny(sptr, slen, scap)
	copy(sl, str)
}

//export scheduleTimeoutEvent
func scheduleTimeoutEvent(_ unsafe.Pointer, _ int32) {
	panic("scheduleTimeoutEvent")
}

//export clearTimeoutEvent
func clearTimeoutEvent(_ unsafe.Pointer, _ int32) {
	panic("clearTimeoutEvent")
}

//export copyBytesToJS
func copyBytesToJS(ctx unsafe.Pointer, sp int32) {
	copyBytesToJSTiny(ctx, sp+40, sp+8, sp+16, sp+24, sp+32, 0, 0)
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

//export copyBytesToGo
func copyBytesToGo(ctx unsafe.Pointer, sp int32) {
	copyBytesToGoTiny(ctx, sp+40, sp+8, sp+16, sp+24, sp+32, 0, 0)
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

var goImports = map[string][]importList{
	// this set satisfies go1.14
	"go": {
		{"debug", debug, C.debug},
		{"runtime.wasmExit", wexit, C.wexit},
		{"runtime.wasmWrite", wwrite, C.wwrite},
		{"runtime.nanotime", nanotime, C.nanotime},
		{"runtime.walltime", walltime, C.walltime},
		{"runtime.scheduleCallback", scheduleCallback, C.scheduleCallback},
		{"runtime.clearScheduledCallback", clearScheduledCallback, C.clearScheduledCallback},
		{"runtime.getRandomData", getRandomData, C.getRandomData},
		{"runtime.scheduleTimeoutEvent", scheduleTimeoutEvent, C.scheduleTimeoutEvent},
		{"runtime.clearTimeoutEvent", clearTimeoutEvent, C.clearTimeoutEvent},
		{"syscall/js.stringVal", stringVal, C.stringVal},
		{"syscall/js.valueGet", valueGet, C.valueGet},
		{"syscall/js.valueSet", valueSet, C.valueSet},
		{"syscall/js.valueIndex", valueIndex, C.valueIndex},
		{"syscall/js.valueSetIndex", valueSetIndex, C.valueSetIndex},
		{"syscall/js.valueCall", valueCall, C.valueCall},
		{"syscall/js.valueInvoke", valueInvoke, C.valueInvoke},
		{"syscall/js.valueNew", valueNew, C.valueNew},
		{"syscall/js.valueLength", valueLength, C.valueLength},
		{"syscall/js.valuePrepareString", valuePrepareString, C.valuePrepareString},
		{"syscall/js.valueLoadString", valueLoadString, C.valueLoadString},
		{"syscall/js.copyBytesToGo", copyBytesToGo, C.copyBytesToGo},
		{"syscall/js.copyBytesToJS", copyBytesToJS, C.copyBytesToJS},
	},
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

// addImports adds go Bridge imports in "go" namespace.
func (b *Bridge) addImports(imps *wasmer.Imports) error {
	var err error
	for k, implist := range goImports {
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
