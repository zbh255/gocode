package bs

import (
	"encoding/binary"
	"fmt"
	"sync"
)

type BoxSwap struct {
	box  []byte
	pool sync.Pool
}

func NewBoxSwap(key, iv []byte) (*BoxSwap, error) {
	if len(key) > 256 {
		return nil, fmt.Errorf("key must have a length of 255")
	}
	if len(key) < 256 {
		key2 := make([]byte, 256)
		copy(key2, key)
		for i := len(key); i < len(key2); i++ {
			key2[i] = key[i%len(key)]
		}
		key = key2
	}
	bs := &BoxSwap{
		box: key,
	}
	for i := 0; i < len(iv); i++ {
		idx := i % len(bs.box)
		bs.box[idx] = byte(uint32(bs.box[idx]) * uint32(iv[i]) % 256)
	}
	bs.pool = sync.Pool{
		New: func() interface{} { return make([]byte, len(bs.box)) },
	}
	return bs, nil
}

func (bs *BoxSwap) Encrypt(src []byte) []byte {
	state := bs.pool.Get().([]byte)
	defer bs.pool.Put(state)
	copy(state, bs.box)
	res := make([]byte, 0, len(src))
	maxLen := len(state)
	for i := 0; i < len(src); i += 1 {
		b2Idx := i % maxLen
		v := src[i]
		v2 := state[b2Idx]
		resV := v ^ v2
		res = append(res, resV)
		state[resV] = v
		state[v] = resV
	}
	return res
}

func (bs *BoxSwap) EncryptUint64(src uint64) uint64 {
	return binary.BigEndian.Uint64(bs.Encrypt(binary.BigEndian.AppendUint64(nil, src)))
}
