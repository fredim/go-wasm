package wasm

import (
	"encoding/binary"

	"github.com/wasmerio/go-ext-wasm/wasmer"
)

func TinyBridgeFromBytes(name string, bytes []byte, imports *wasmer.Imports) (*Bridge, error) {
	b := new(Bridge)
	if imports == nil {
		imports = wasmer.NewImports()
	}

	b.name = name
	err := b.addTinyImports(imports)
	if err != nil {
		return nil, err
	}

	inst, err := wasmer.NewInstanceWithImports(bytes, imports)
	if err != nil {
		return nil, err
	}

	ctx, err := getCtxData(b)
	if err != nil {
		return nil, err
	}

	b.instance = inst
	inst.SetContextData(ctx)
	b.addValues()
	b.refs = make(map[interface{}]int)
	b.valueIDX = 8
	return b, nil
}

func TinyBridgeFromFile(name, file string, imports *wasmer.Imports) (*Bridge, error) {
	bytes, err := wasmer.ReadBytes(file)
	if err != nil {
		return nil, err
	}

	return TinyBridgeFromBytes(name, bytes, imports)
}

func (b *Bridge) loadSliceTiny(arr, alen, acap int32) []byte {
	mem := b.mem()
	array := binary.LittleEndian.Uint64(mem[arr:])
	length := binary.LittleEndian.Uint64(mem[alen:])
	return mem[array : array+length]
}

func (b *Bridge) loadStringTiny(ptr, plen int32) string {
	d := b.loadSliceTiny(ptr, plen, 0)
	return string(d)
}

func (b *Bridge) loadSliceOfValuesTiny(arr, alen, acap int32) []interface{} {
	array := int(b.getInt64(arr))
	length := int(b.getInt64(alen))
	vals := make([]interface{}, length, length)
	for i := 0; i < length; i++ {
		vals[i] = b.loadValue(int32(array + i*8))
	}

	return vals
}
